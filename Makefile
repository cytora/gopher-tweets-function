.PHONY: build deploy start test

build:
	sam build

deploy:
	echo $(TWITTER_CONSUMER_KEY) && sam deploy
	    --parameter-overrides TwitterConsumerKey=$(TWITTER_CONSUMER_KEY)
	    --no-confirm-changeset --no-fail-on-empty-changeset

start:
	sam local start-api -n env.json

test:
	go test ./...