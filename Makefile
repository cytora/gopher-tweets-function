.PHONY: build deploy

build:
	sam build

deploy:
	sam deploy --no-confirm-changeset --no-fail-on-empty-changeset
