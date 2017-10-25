# Vault Token Helper for OS X Keychain

This is sample code and a proof-of-concept for creating an external
[HashiCorp Vault](https://www.vaultproject.io) Token Helper.

By default, Vault authenticates users locally and caches their token in
`~/.vault-token`. For shared systems or systems where security is paramount,
this may not be ideal. Fortunately, this storage mechanism is an abstraction
known as a "token helper".

This code demonstrates one possible example of an external token helper. When
requesting or storing a token, Vault delegates to this binary.


## Installation

1. Download and install the binary from GitHub. I supplied both a signed DMG
with my personal Apple Developer ID or you can download the binary directly. If
neither of those options suffice, you can audit and compile the code yourself.

1. Put the binary somewhere on disk, like `~/.vault.d/token-helpers`:

    ```sh
    $ mv vault-token-helper ~/.vault.d/token-helpers/vault-token-helper
    ```

1. Create a Vault configuration file at `~/.vault` with the contents:

    ```hcl
    token_helper = "/Users/<your username>/.vault.d/token-helpers/vault-token-helper"
    ```

    Be sure to replace `<your username>` with your username. The value must be
    a full path (you cannot use a relative path).

    The local CLI will automatically pickup this configuration value.


## Usage

1. Use Vault normally. Commands like `vault auth` will automatically delegate to
keychain access.


## Development

There's a handy `scripts/dev.sh` that will start a Vault server in development
mode pre-configured with the token helper.


## License & Author

This project is licensed under the MIT license by Seth Vargo
(seth@sethvargo.com).
