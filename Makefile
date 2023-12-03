BINDIR := $(CURDIR)/bin
BINNAME := artifactory-cleanup
TARGET_BIN := $(BINDIR)/$(BINNAME)
INSTALL_PATH := /usr/local/bin

# ------------------------------------------------------------------------------

.PHONY: build
build:
	@command -v go &> /dev/null || (echo "Please install GoLang" && false)
	go mod download
	go build -ldflags="-s -w" -trimpath -o '$(TARGET_BIN)' ./cmd/clean

.PHONY: clean
clean:
	-@rm -rf $(BINDIR)

.PHONY: install
install:
	@test -e "$(BINDIR)/$(BINNAME)" &> /dev/null || (echo "There are no executable file. Please run 'make build'" && false)
	@install "$(TARGET_BIN)" "$(INSTALL_PATH)/$(BINNAME)"
	-@$(MAKE) clean
	@echo "Binary file installed to '$(INSTALL_PATH)/$(BINNAME)'"

.PHONY: uninstall
uninstall: clean
	@echo "Removing: $(INSTALL_PATH)/$(BINNAME)"
	@rm -f "$(INSTALL_PATH)/$(BINNAME)"

.PHONY: test
test:
	@go test ./... -count=1 -cover
