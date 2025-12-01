.PHONY: build
build:
	mkdir -p dist
	GOOS=linux go build -o dist/macra-backend -a cmd/main.go

.PHONY: deploy
deploy:
	scp -i $(MACRA_SSHKEY) -r dist/macra-backend $(MACRA_HOST):~/
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl stop macra-backend
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo cp ~/macra-backend /opt/macra/macra-backend
	ssh -i $(MACRA_SSHKEY) $(MACRA_HOST) sudo systemctl start macra-backend