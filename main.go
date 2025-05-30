package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	API struct {
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"api"`
}

type RawToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Token struct {
	RawToken
	Expiration int `json:"expiration"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func fetchToken(config *Config, path string) (*Token, error) {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var cachedToken Token
		err = json.Unmarshal(data, &cachedToken)
		if err != nil {
			return nil, err
		}

		if cachedToken.Expiration >= int(time.Now().Unix()) {
			return &cachedToken, nil
		}
	}

	reqBody := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{config.API.ClientID},
		"client_secret": []string{config.API.ClientSecret},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://accounts.spotify.com/api/token",
		strings.NewReader(reqBody.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawToken RawToken
	err = json.Unmarshal(respBody, &rawToken)
	if err != nil {
		return nil, err
	}

	token := Token{
		RawToken:   rawToken,
		Expiration: int(time.Now().Unix()) + rawToken.ExpiresIn - 15,
	}

	tokenStr, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(path, tokenStr, 0644)

	return &token, nil
}

func main() {
	config, err := loadConfig("./config.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", config)

	token, err := fetchToken(config, "./token.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", token)
}
