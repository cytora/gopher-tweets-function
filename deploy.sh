#!/usr/bin/env bash
sam deploy --parameter-overrides "TwitterConsumerKey=$TWITTER_CONSUMER_KEY" "TwitterConsumerSecret=$TWITTER_CONSUMER_SECRET" --no-confirm-changeset --no-fail-on-empty-changeset
