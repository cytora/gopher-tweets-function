.PHONY: build deploy start

build:
	sam build

deploy:
	sam deploy --no-confirm-changeset --no-fail-on-empty-changeset

start:
	sam local start-api -n env.json