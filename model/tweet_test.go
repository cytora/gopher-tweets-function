package model

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTweetRequest(t *testing.T) {
	original := TweetRequest{Tweet: "hello from the other side"}
	obytes, err := json.Marshal(original)
	assert.Nil(t, err)
	closer := ioutil.NopCloser(bytes.NewReader(obytes))
	req := http.Request{Body: closer}
	post, err := NewTweetRequest(&req)
	assert.Nil(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, original, *post)
}
