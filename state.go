package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type State struct {
	TTL      time.Time `json:"ttl"`
	Policies []string  `json:"policies"`
}

func StateFromRequest(r *http.Request, cookieName, cookieEncKey string) (*State, []byte, error) {
	cookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return nil, nil, http.ErrNoCookie
	}
	if err != nil {
		return nil, nil, err
	}
	decodedCookie, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, nil, err
	}
	encryptedCookie := []byte(decodedCookie)
	nonce := encryptedCookie[:12]
	encryptedCookie = encryptedCookie[12:]
	if len(nonce) != 12 {
		return nil, nil, errors.New("Nonce must be 12 characters in length")
	}
	if len(encryptedCookie) == 0 {
		return nil, nil, errors.New("Encrypted Cookie missing")
	}
	b, err := decrypt(encryptedCookie, nonce, []byte(cookieEncKey))
	if err != nil {
		return nil, nil, err
	}
	var state *State
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&state)
	if err != nil {
		return nil, nil, err
	}
	return state, b, nil
}
