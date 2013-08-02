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
	Bcc            []string
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
	return Host + "/api/" + command + ".json?api_user=" +
		url.QueryEscape(c.username) + "&api_key=" + c.password
}

func (c Client) get(u uri) *SendGridError {
	resp, err := http.Get(string(u))
	if err != nil {
		return &SendGridError{Message: err.Error()}
	}
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
	if len(response.Errors) > 0 || resp.StatusCode != 200 || response.Message != "success" {
		return &SendGridError{StatusCode: resp.StatusCode, Response: response}
	}

	return nil
}

// MailSend will send out emal o list of given recipients.
func (c Client) MailSend(args MailArgs) *SendGridError {
	// Create the address.
	u := uri(c.getBase("mail.send"))
	u.appendArray("to", args.To)
	u.appendArray("toname", args.ToName)
	u.appendArray("bcc", args.Bcc)
	u.appendValue("from", args.From)
	u.appendValue("fromname", args.FromName)
	u.appendValue("subject", args.Subject)
	u.appendValue("text", args.Text)
	u.appendValue("html", args.Html)
	if err := u.appendMap("x-smtapi", args.Xsmtapi); err != nil {
		return &SendGridError{Message: err.Error()}
	}
	if err := u.appendMap("content", args.Content); err != nil {
		return &SendGridError{Message: err.Error()}
	}
	if err := u.appendMap("headers", args.Headers); err != nil {
		return &SendGridError{Message: err.Error()}
	}

	// Make the request.
	return c.get(u)
}
