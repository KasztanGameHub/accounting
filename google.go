package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}


func validateGoogleJWT(cred string) (accClaims, error) {
	token, err := jwt.ParseWithClaims(cred, &accClaims{}, func (t *jwt.Token) (interface{}, error) {
		pem, err := getGooglePublicKey(fmt.Sprintf("%s", t.Header["kid"]))
		if err != nil {
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, err
		}
		return key, nil
	})

	if err != nil {
		return accClaims{}, err
	}

	claims, ok := token.Claims.(*accClaims)
	if !ok {
		return accClaims{}, err
	}

	if claims.Audience != GOOGLE_CLIENT_ID {
		return accClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return accClaims{}, err
	}

	return *claims, nil
}
