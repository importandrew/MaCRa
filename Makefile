.PHONY: build
build:
	$(MAKE) -C backend build
	$(MAKE) -C frontend build

.PHONY: deploy
deploy:
	$(MAKE) -C backend deploy
	$(MAKE) -C frontend deploy

.PHONY: deploy-dev
deploy-dev:
	$(MAKE) -C backend deploy-dev
	$(MAKE) -C frontend deploy-dev