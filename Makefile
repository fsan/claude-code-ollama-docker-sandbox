MODEL ?= glm-5:cloud
WORKSPACE ?= $(PWD)
OLLAMA_PORT ?= 11434
FLAGS ?=

# Binary location
BIN_DIR := bin
CLOMA_BIN := $(BIN_DIR)/cloma

.DEFAULT_GOAL := run
.PHONY: setup run doctor shell logs stop clean template template-clean
.PHONY: build install cloma cloma-setup

setup:
	CLAUDE_CODE_MODEL="$(MODEL)" OLLAMA_PORT="$(OLLAMA_PORT)" ./scripts/setup.sh "$(WORKSPACE)"

run:
	CLAUDE_CODE_MODEL="$(MODEL)" OLLAMA_PORT="$(OLLAMA_PORT)" CLAUDE_CODE_FLAGS="$(FLAGS)" ./scripts/run-claude-code.sh "$(WORKSPACE)"

doctor:
	CLAUDE_CODE_MODEL="$(MODEL)" OLLAMA_PORT="$(OLLAMA_PORT)" ./scripts/doctor.sh "$(WORKSPACE)"

shell:
	./scripts/shell.sh "$(WORKSPACE)"

logs:
	./scripts/logs.sh "$(WORKSPACE)"

stop:
	./scripts/stop-sandbox.sh "$(WORKSPACE)"

clean:
	./scripts/clean-sandbox.sh "$(WORKSPACE)"

template:
	./scripts/bake-template.sh

template-clean:
	./scripts/clean-template.sh

# Go CLI targets
build: $(CLOMA_BIN)

$(CLOMA_BIN):
	go build -o $(CLOMA_BIN) ./cmd/cloma

install: build
	@echo "Installing cloma to /usr/local/bin..."
	sudo cp $(CLOMA_BIN) /usr/local/bin/cloma
	@echo "Installed: /usr/local/bin/cloma"

cloma: build
	./$(CLOMA_BIN) $(ARGS)

cloma-setup: build
	@echo "Creating ~/.cloma directory structure..."
	./$(CLOMA_BIN) doctor