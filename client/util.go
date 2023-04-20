package client

import (
	"fmt"
	"net/http"
)

func checkStatusCode(response *http.Response) error {
	switch response.StatusCode {
	case 400:
		return fmt.Errorf("bad request to Jane.  Please file a ticket")
	case 401:
		return fmt.Errorf("bad authentication.  Login with 'auth -p ${password}'")
	case 404:
		return fmt.Errorf("bad request to Jane.  Please file a ticket")

	case 500:
		fallthrough
	case 502:
		fallthrough
	case 503:
		fallthrough
	case 504:
		return fmt.Errorf("Jane API isn't responding to requests.  Please try again shortly")
	}
	return nil
}
