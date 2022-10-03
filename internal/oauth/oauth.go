// Copyright (C) 2018  denis4net
// Copyright (C) 2019-2022 Evgeny Kuznetsov (evgeny@kuznetsov.md)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package oauth

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	photoslibrary "evgenykuznetsov.org/go/gphotoslibrary"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// newToken returns a *oauth2.Token for the given credentials
func newToken(ctx context.Context, clientID, clientSecret, redirectURL string) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		RedirectURL:  redirectURL,
	}
	state, err := generateOAuthState()
	if err != nil {
		return nil, err
	}
	authCodeURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Open the following URL in your browser:\n%s\n", authCodeURL)
	fmt.Println("After authentication, copy the final URL from the browser and paste it here:")

	var uri string
	if _, err := fmt.Scanln(&uri); err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	st := q.Get("state")
	if st != state {
		return nil, fmt.Errorf("state mismatch")
	}
	authCode := q.Get("code")

	return config.Exchange(ctx, authCode)
}

// newClientFromToken creates a new http.Client with a bearer access token
func newClientFromToken(ctx context.Context, clientID, clientSecret, redirectURL string, accessToken *oauth2.Token) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		RedirectURL:  redirectURL,
	}

	return config.Client(ctx, accessToken), nil
}

// NewClient returns a *http.Client ready to be worked with
func NewClient(ctx context.Context, clientID, clientSecret, redirectURL, tokenPath string) (*http.Client, error) {
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		token, err := newToken(ctx, clientID, clientSecret, redirectURL)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(token)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(tokenPath, data, 0600)
		if err != nil {
			return nil, err
		}

		return newClientFromToken(ctx, clientID, clientSecret, redirectURL, token)
	} else {
		data, err := os.ReadFile(tokenPath)
		if err != nil {
			return nil, err
		}
		token := &oauth2.Token{}
		err = json.Unmarshal(data, token)
		if err != nil {
			return nil, err
		}
		fmt.Printf("read a token from \"%s\": %s %s (%s)\n", tokenPath, clientID, clientSecret, redirectURL)
		return newClientFromToken(ctx, clientID, clientSecret, redirectURL, token)
	}
}

func generateOAuthState() (string, error) {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", n), nil
}
