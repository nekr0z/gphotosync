package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	photoslibrary "github.com/nekr0z/gphotoslibrary"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// NewOAuthClient creates a new http.Client with a bearer access token
func NewOAuthToken(ctx context.Context, clientID string, clientSecret string) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	}
	state, err := generateOAuthState()
	if err != nil {
		return nil, err
	}
	authCodeURL := config.AuthCodeURL(state)
	fmt.Printf("Open %s\n", authCodeURL)
	fmt.Print("Enter code: ")

	var authCode string
	if _, err := fmt.Scanln(&authCode); err != nil {
		return nil, err
	}

	return config.Exchange(ctx, authCode)
}

// NewOAuthClientFromToken creates a new http.Client with a bearer access token
func NewOAuthClientFromToken(ctx context.Context, clientID string, clientSecret string, accessToken *oauth2.Token) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	}

	return config.Client(ctx, accessToken), nil
}

func NewOAuthClient(ctx context.Context, clientID, clientSecret, tokenPath string) (*http.Client, error) {
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		token, err := NewOAuthToken(ctx, clientID, clientSecret)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(token)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(tokenPath, data, 0600)
		if err != nil {
			return nil, err
		}

		return NewOAuthClientFromToken(ctx, clientID, clientSecret, token)
	} else {
		data, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			return nil, err
		}
		token := &oauth2.Token{}
		err = json.Unmarshal(data, token)
		if err != nil {
			return nil, err
		}
		fmt.Printf("read a token from \"%s\": %s %s\n", tokenPath, GoogleClientId, GoogleClientSecret)
		return NewOAuthClientFromToken(ctx, clientID, clientSecret, token)
	}
}

func generateOAuthState() (string, error) {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", n), nil
}
