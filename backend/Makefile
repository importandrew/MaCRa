BIN_NAME := macra-backend

MAIN_SERVICE = macra-backend
DEV_SERVICE := dev-macra-backend

MAIN_DIR := /opt/macra
DEV_DIR := /opt/macra-dev

.PHONY: build
build:
	mkdir -p dist
	GOOS=linux go build -o dist/$(BIN_NAME) -a cmd/main.go

.PHONY: deploy
deploy:
	scp -i $(MACRA_SSHKEY) -r dist/$(BIN_NAME) $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl stop $(MAIN_SERVICE)
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp $(BIN_NAME) $(MAIN_DIR)/$(BIN_NAME)
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl start $(MAIN_SERVICE)

.PHONY: deploy-dev
deploy-dev:
	scp -i $(MACRA_SSHKEY) -r dist/$(BIN_NAME) $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl stop $(DEV_SERVICE)
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp $(BIN_NAME) $(DEV_DIR)/$(BIN_NAME)
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl start $(DEV_SERVICE)