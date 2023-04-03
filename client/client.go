package client

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
)

/*
path=/;
expires=Wed, 03 Apr 2024 19:35:14 GMT;
HttpOnly;
secure;
SameSite=Lax
*/

type Client struct {
	janeClient *http.Client
	logger     logger.Logger
	config     Config
	auth       *domain.Auth
	updateAuth func() error
}

func NewClient(
	logger logger.Logger,
	config Config,
	auth *domain.Auth,
	updateAuth func() error,
) (Client, error) {
	output := Client{
		logger:     logger,
		config:     config,
		auth:       auth,
		updateAuth: updateAuth,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return Client{}, err
	}

	if auth.AuthCookie != "" || auth.Expires.After(time.Now()) {
		authCookie := http.Cookie{
			Name:     authCookieKey,
			Value:    auth.AuthCookie,
			Path:     "/",
			Domain:   "",
			Expires:  auth.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		urlDomain, err := url.Parse(output.getDomain())
		if err != nil {
			return output, fmt.Errorf(
				"provided domain %v doesn't correctly parse to URL: %v",
				auth.Domain,
				err,
			)
		}
		jar.SetCookies(urlDomain, []*http.Cookie{&authCookie})
	}

	output.janeClient = &http.Client{
		Jar: jar,
	}
	return output, nil
}

func (client Client) getDomain() string {
	//return "http://localhost:2345"
	return fmt.Sprintf("https://%v.janeapp.com", client.auth.Domain)
}
