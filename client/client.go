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

const apiBase2 = "admin/api/v2"
const apiBase3 = "admin/api/v3"

type Client struct {
	janeClient *http.Client
	logger     logger.Logger
	config     Config
	user       *domain.User
	updateAuth func() error
}

func NewClient(
	logger logger.Logger,
	config Config,
	user *domain.User,
	updateAuth func() error,
) (Client, error) {
	output := Client{
		logger:     logger,
		config:     config,
		user:       user,
		updateAuth: updateAuth,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return Client{}, err
	}

	if user.Auth.AuthCookie != "" || user.Auth.Expires.After(time.Now()) {
		authCookie := http.Cookie{
			Name:     authCookieKey,
			Value:    user.Auth.AuthCookie,
			Path:     "/",
			Domain:   "",
			Expires:  user.Auth.Expires,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		urlDomain, err := url.Parse(output.getDomain())
		if err != nil {
			return output, fmt.Errorf(
				"provided domain %v doesn't correctly parse to URL: %v",
				user.Auth.Domain,
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
	return fmt.Sprintf("https://%v.janeapp.com", client.user.Auth.Domain)
}
