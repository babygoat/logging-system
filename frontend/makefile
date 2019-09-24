BIN_DIR ?= node_modules/.bin

help:
	@echo "make dev -- start dev server"
	@echo "make build -- transpile es6 files to es5"
	@echo "make start -- start application server"
	@echo "make clean -- clean the old builds"

dev:
	$(BIN_DIR)/nodemon src/index.js --exec $(BIN_DIR)/babel-node

build: clean
	$(BIN_DIR)/babel src/index.js --out-dir dist

start:
	node dist/index.js

clean:
	$(BIN_DIR)/rimraf dist

.PHONY: help build start dev
