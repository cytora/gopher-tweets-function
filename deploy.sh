#!/usr/bin/env bash
sam deploy --parameter-overrides "TwitterConsumerKey=$TWITTER_CONSUMER_KEY" --no-confirm-changeset --no-fail-on-empty-changeset
