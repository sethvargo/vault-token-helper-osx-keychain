package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/keybase/go-keychain"
	"github.com/pkg/errors"
)

var (
	stdin  = os.Stdin
	stderr = os.Stderr
	stdout = os.Stdout

	// account is the name of the current user, which is where the keychain item
	// will be stored.
	account string

	accessGroup = "com.hootsuite.vault-token-helper-osx-keychain"
)

const defaultVaultServer = "Default HashiCorp Vault Server"

func init() {
	// Get the current user. This requires cgo, but we already need cgo because
	// we are binding to objective-C libraries.
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	account = u.Username
}

func provideVaultAddress() string {
	addr, isSet := os.LookupEnv("VAULT_ADDR")

	if !isSet {
		return defaultVaultServer
	}

	addr = strings.TrimSpace(addr)
	addr = strings.ToLower(addr)
	return addr
}

func main() {

	if err := realMain(); err != nil {
		fmt.Fprintf(stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func realMain() error {
	args := os.Args[1:]
	if len(args) < 1 {
		return fmt.Errorf("expected at lease 1 argument")
	}

	switch args[0] {
	case "store":
		return handleStore()
	case "get":
		return handleGet()
	case "erase":
		return handleErase()
	default:
		return fmt.Errorf("invalid command %q", args[0])
	}
}

// handleGet retrieves the stored value in the keychain. A missing item is not
// an error, instead the empty string is returned. Errors are returned if there
// are communication issues with the keychain itself.
func handleGet() error {
	query := keychainItem()
	query.SetReturnData(true)
	query.SetMatchLimit(keychain.MatchLimitOne)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return errors.Wrap(err, "failed to query keychain")
	}
	if len(results) == 0 {
		return nil
	}

	fmt.Fprintf(stdout, "%s", results[0].Data)
	return nil
}

// handleStore saves the given value in the keychain, returning any errors that
// occur while attempting to persist the value.
func handleStore() error {
	r := bufio.NewReader(stdin)
	value, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to read value from stdin")
	}

	item := keychainItem()
	item.SetData([]byte(value))

	err = keychain.AddItem(item)

	// Keychain items are not upserted, so we have to catch the duplicate error
	// and attempt an update instead.
	if err == keychain.ErrorDuplicateItem {
		query := keychainItem()
		query.SetReturnData(true)
		query.SetMatchLimit(keychain.MatchLimitOne)

		results, err := keychain.QueryItem(query)
		if err != nil {
			return errors.Wrap(err, "failed to query keychain")
		}
		if len(results) == 0 {
			return errors.New("no results")
		}

		if err := keychain.UpdateItem(query, item); err != nil {
			return errors.Wrap(err, "failed to update item in keychain")
		}
		return nil
	}

	// Handle any other errors
	if err != nil {
		return errors.Wrap(err, "failed to add item to keychain")
	}

	return nil
}

// handleErase removes the entry from the keychain, if it exists. If the entry
// does not exist, the function still returns as successful. Errors are returned
// if there are communication issues with the keychain itself.
func handleErase() error {
	item := keychainItem()
	if err := keychain.DeleteItem(item); err != nil {
		return errors.Wrap(err, "failed to delete item")
	}
	return nil
}

// keychainItem constructs a new keychain item with the given data. This can be
// used as a query or to insert/update a keychain item.
func keychainItem() keychain.Item {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(provideVaultAddress())
	item.SetAccount(account)
	item.SetLabel(provideVaultAddress())
	item.SetAccessGroup(accessGroup)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)
	return item
}
