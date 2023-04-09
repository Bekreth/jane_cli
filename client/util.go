package client

import (
	"fmt"
	"net/http"
)

func checkStatusCode(response *http.Response) error {
	switch response.StatusCode {
	case 401:
		return fmt.Errorf("bad authentication.  Login with 'auth -p ${password'")
	case 404:
		return fmt.Errorf("bad request to Jane.  Please file a ticket")
	}
	return nil
}
