NAME=neon
VERSION=$(shell changelog release version)
BUILD_DIR=build

YELLOW=\033[93m
CYAN=\033[1m\033[96m
CLEAR=\033[0m

.PHONY: build test

all: clean test build

help:
	@echo "$(YELLOW)Print help$(CLEAR)"
	@echo "$(CYAN)help$(CLEAR)    Print this help screen"
	@echo "$(CYAN)test$(CLEAR)    Run tests"
	@echo "$(CYAN)build$(CLEAR)   Build executable"
	@echo "$(CYAN)archive$(CLEAR) Build binary archive"
	@echo "$(CYAN)release$(CLEAR) Make a release"
	@echo "$(CYAN)clean$(CLEAR)   Clean generated files"

test:
	@echo "$(YELLOW)Running test$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go run $(NAME).go test.md > $(BUILD_DIR)/test.md
	cat $(BUILD_DIR)/test.md
	
build:
	@echo "$(YELLOW)Building executable$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

archive: clean
	@echo "$(YELLOW)Building binary archive$(CLEAR)"
	mkdir -p $(BUILD_DIR)/$(NAME)-$(VERSION)/
	gox -output=$(BUILD_DIR)/$(NAME)-$(VERSION)/{{.Dir}}_{{.OS}}_{{.Arch}}
	cp LICENSE.txt $(BUILD_DIR)/$(NAME)-$(VERSION)/
	cp README.md $(BUILD_DIR)/ && cd $(BUILD_DIR) && md2pdf README.md && cp README.pdf $(NAME)-$(VERSION)/
	cd $(BUILD_DIR) && tar cvzf $(NAME)-bin-$(VERSION).tar.gz $(NAME)-$(VERSION)

release: clean test archive
	@echo "$(YELLOW)Making release $(VERSION)$(CLEAR)"
	release

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)
