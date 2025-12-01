.PHONY: build
build:
	$(MAKE) -C backend build
	$(MAKE) -C frontend build

.PHONY: deploy
deploy:
	$(MAKE) -C backend deploy
	$(MAKE) -C frontend deploy