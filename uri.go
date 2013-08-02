package sendgrid

import (
	"bytes"
	"encoding/json"
	"net/url"
)

type uri string

func (u *uri) String() string {
	return string(*u)
}

// appendArray adds a list in url format.
func (u *uri) appendArray(name string, list []string) {
	for _, value := range list {
		*u += uri("&" + name + "[]=" + url.QueryEscape(value))
	}
}

// appendValue adds a string value to the url.
func (u *uri) appendValue(name, value string) {
	if len(value) == 0 {
		return
	} else {
		*u += uri("&" + name + "=" + url.QueryEscape(value))
	}
}

// appendMap adds a key value pairing map to the uri
func (u *uri) appendMap(name string, data map[string]string) error {
	if len(data) == 0 {
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	*u += uri(bytes.NewBuffer(b).String())
	return nil
}
