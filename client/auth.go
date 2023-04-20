package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Bekreth/jane_cli/domain"
)

const authCookieKey = "_front_desk_session"
const authQueryParameters = "auth_key=%v&password=%v"
const authPath = "auth/identity/callback"
const failedPath = "/auth/failure"
const idPath = "admin/api/v2/staff_members"

func (client Client) Login(password string) error {
	urlBase := fmt.Sprintf("%v/%v", client.getDomain(), authPath)
	client.logger.Infof("logging in to %v", urlBase)
	loginCredentials := fmt.Sprintf(
		authQueryParameters,
		client.user.Auth.Username,
		password,
	)
	authBody := bytes.NewBufferString(loginCredentials)
	request, err := http.NewRequest(
		http.MethodPost,
		urlBase,
		authBody,
	)
	if err != nil {
		client.logger.Infof("failed to build auth request: %v", err)
		return err
	}

	request.ContentLength = int64(len(loginCredentials))
	request.Header.Del("Accept-Encoding")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("got a bad auth response: %v", err)
		return err
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == authCookieKey {
			client.logger.Debugf("Got a new auth cookie until %v", cookie.Expires)
			client.user.Auth.AuthCookie = cookie.Value
			client.user.Auth.Expires = cookie.Expires
			userID, err := client.getUserID()
			if err != nil {
				return err
			}
			client.user.Auth.UserID = userID
			return client.updateAuth()
		}
	}

	return fmt.Errorf("no cookie was provided")
}

func (client Client) getUserID() (int, error) {
	urlBase := fmt.Sprintf("%v/%v", client.getDomain(), idPath)
	client.logger.Infof("getting user ID for %v", client.user.Auth.Username)
	request, err := http.NewRequest(
		http.MethodGet,
		urlBase,
		nil,
	)

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("unable to get userID: %v", err)
		return 0, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read bytes from response: %v", err)
		return 0, err
	}

	output := []domain.StaffMember{}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		client.logger.Infof("failed to parse json response: %v", err)
		return 0, err
	}

	for _, staffMember := range output {
		if staffMember.Email == client.user.Auth.Username {
			client.logger.Debugf(
				"got client ID of %v for user %v",
				staffMember.ID,
				client.user.Auth.Username,
			)
			client.user.Auth.UserID = staffMember.ID
			return staffMember.ID, err
		}
	}

	return 0, nil
}
