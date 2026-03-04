MAIN_DIR := /opt/macra-dist
DEV_DIR := /opt/macra-dist-dev

.PHONY: build
build:
	npm run build

.PHONY: deploy
deploy:
	scp -i $(MACRA_SSHKEY) -r dist $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo rm -rf $(MAIN_DIR)/*
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp -r dist/* $(MAIN_DIR)/

.PHONY: deploy-dev
deploy-dev:
	scp -i $(MACRA_SSHKEY) -r dist $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo rm -rf $(DEV_DIR)/*
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp -r dist/* $(DEV_DIR)/