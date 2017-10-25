CODESIGN ?=

build:
	@echo "==> Building"
	@go build -a -o bin/vault-token-helper

deps:
	@echo "==> Updating dependencies"
	@dep ensure -update

dist: build dmg sign
	@echo "==> Moving everything into pkg/"
	@cp bin/vault-token-helper pkg/vault-token-helper

dmg:
	@echo "==> Building dmg"
	@mkdir -p pkg/
	@hdiutil create -ov -volname "Vault Token Helper" -srcfolder bin/ "pkg/vault-token-helper.dmg"

sign:
	@echo "==> Signing"
	@codesign --force --sign "${CODESIGN}" "pkg/vault-token-helper.dmg"
	@spctl -a -t open --context "context:primary-signature" "pkg/vault-token-helper.dmg"
