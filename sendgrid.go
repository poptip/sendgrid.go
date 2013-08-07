package sendgrid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	Host = "https://sendgrid.com"
)

type PostResponse struct {
	Message string
	Errors  []string
}

type MailArgs struct {
	To             []string
	ToName         []string
	BCC            []string
	Xsmtapi        map[string]string
	Subject        string
	Text, Html     string
	From, FromName string
	ReplyTo        string
	Date           string
	Files          []string
	Content        map[string]string
	Headers        map[string]string
}

type Client struct {
	username, password string
}

func NewClient(username, password string) *Client {
	return &Client{username: username, password: password}
}

// getBase returns the base url that will be used for the API request.
func (c Client) getBase(command string) string {
	return Host + "/api/" + command + ".json?"
}

func (c Client) setCredentials(v *url.Values) {
	v.Set("api_user", c.username)
	v.Set("api_key", c.password)
}

func (c Client) get(u string, v url.Values) *SendGridError {
	resp, err := http.Get(u + v.Encode())
	if err != nil {
		return &SendGridError{Message: err.Error()}
	}
	return c.response(resp)
}

func (c Client) post(u string, v url.Values) *SendGridError {
	resp, err := http.PostForm(u, v)
	if err != nil {
		return &SendGridError{Message: err.Error()}
	}
	return c.response(resp)
}

func (c Client) response(resp *http.Response) *SendGridError {
	defer resp.Body.Close()

	// Read the JSON message from the body.
	response := PostResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return &SendGridError{
			Message: fmt.Sprintf("JSON unmarshalling failed: %q", err.Error()),
		}
	}

	// Check if the response is successful.
	if len(response.Errors) > 0 ||
		resp.StatusCode != 200 ||
		response.Message != "success" {
		return &SendGridError{StatusCode: resp.StatusCode, Response: response}
	}

	return nil
}

// MailSend will send out emal o list of given recipients.
func (c Client) MailSend(args MailArgs) *SendGridError {
	// Create the address.
	base := c.getBase("mail.send")
	v := url.Values{}
	c.setCredentials(&v)
	addArray(&v, "to", args.To)
	addArray(&v, "toname", args.ToName)
	addArray(&v, "bcc", args.BCC)
	v.Add("from", args.From)
	v.Add("fromname", args.FromName)
	v.Add("subject", args.Subject)
	v.Add("text", args.Text)
	v.Add("html", args.Html)
	if err := addMap(&v, "x-smtapi", args.Xsmtapi); err != nil {
		return &SendGridError{Message: err.Error()}
	}
	if err := addMap(&v, "content", args.Content); err != nil {
		return &SendGridError{Message: err.Error()}
	}
	if err := addMap(&v, "headers", args.Headers); err != nil {
		return &SendGridError{Message: err.Error()}
	}

	// Make the request.
	return c.post(base, v)
}
