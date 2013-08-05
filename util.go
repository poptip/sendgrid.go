package sendgrid

import (
	"bytes"
	"encoding/json"
	"net/url"
)

// addArray adds a list in url format.
func addArray(v *url.Values, name string, list []string) {
	for _, value := range list {
		v.Add(name, value)
	}
}

// addMap adds a key value pairing map to the uri
func addMap(v *url.Values, name string, data map[string]string) error {
	if len(data) == 0 {
		return nil
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.Add(name, bytes.NewBuffer(b).String())
	return nil
}
