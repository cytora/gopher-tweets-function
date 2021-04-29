package auth

import (
	"fmt"

	"github.com/gorilla/securecookie"
)

type Credentials struct {
	ExtAccessKey    string
	ExtAccessSecret string
}

func (c *Credentials) Encode(sessionKey string) string {
	codecs := securecookie.CodecsFromPairs([]byte(sessionKey))
	enc, err := securecookie.EncodeMulti("token", c, codecs...)
	if err != nil {
		return ""
	}
	return enc
}

func Decode(sessionKey, encoded string) (*Credentials, error) {
	codecs := securecookie.CodecsFromPairs([]byte(sessionKey))
	var cred Credentials
	err := securecookie.DecodeMulti("token", encoded, &cred, codecs...)
	if err != nil {
		return nil, fmt.Errorf("error decoding credentials:%s", err.Error())
	}
	return &Credentials{
		ExtAccessKey:    cred.ExtAccessKey,
		ExtAccessSecret: cred.ExtAccessSecret,
	}, nil
}
