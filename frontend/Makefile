.PHONY: build
build:
	npm run build

.PHONY: deploy
deploy:
	scp -i $(MACRA_SSHKEY) -r dist $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo rm -rf /opt/macra-dist/*
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp -r dist/* /opt/macra-dist/