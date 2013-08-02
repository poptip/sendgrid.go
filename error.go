package sendgrid

import (
	"strconv"
	"strings"
)

type SendGridError struct {
	StatusCode int
	Response   PostResponse
	Message    string
}

func (e SendGridError) Error() string {
	var message string
	if e.StatusCode > 0 {
		message = "Status Code " + strconv.Itoa(e.StatusCode) + ": "
	}
	if len(e.Response.Errors) > 0 {
		message += strings.Join(e.Response.Errors, ", ")
	} else if len(e.Response.Message) > 0 {
		message += e.Response.Message
	} else {
		message += e.Message
	}
	return message
}
