package client

import (
	"bytes"
	"fmt"
	"net/http"
)

const authCookieKey = "_front_desk_session"
const authQueryParameters = "auth_key=%v&password=%v"
const authPath = "auth/identity/callback"

func (client Client) Login(password string) error {
	urlBase := fmt.Sprintf("https://%v.janeapp.com/%v", client.auth.Domain, authPath)
	client.logger.Infof("logging in to %v", urlBase)
	loginCredentials := fmt.Sprintf(authQueryParameters, client.auth.Username, password)
	authBody := bytes.NewBufferString(loginCredentials)
	request, err := http.NewRequest(
		http.MethodPost,
		urlBase,
		authBody,
	)
	request.ContentLength = int64(len(loginCredentials))
	if err != nil {
		client.logger.Infof("failed to build auth request: %v", err)
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		client.logger.Infof("got a bad auth response: %v", err)
		return err
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == authCookieKey {
			client.logger.Debugf("Got a new auth cookie until %v", cookie.Expires)
			client.auth.AuthCookie = cookie.Value
			client.auth.Expires = cookie.Expires
			return client.updateAuth()
		}
	}
	return fmt.Errorf("no cookie was provided")
}
