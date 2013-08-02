package sendgrid

import (
	"encoding/json"
	"github.com/coocood/assrt"
	"net/http"
	"sync"
	"testing"
)

var (
	originalHost string
	statusCode   int
	results      []byte
	server       *http.Server
	assert       *assrt.Assert
	once         sync.Once
)

func setUp() {
	originalHost = Host
	Host = "http://localhost:12345"
	// Custom handler used for mocking request results.
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(statusCode)
		w.Write(results)
	})
	server = &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}
	go func() {
		assert.MustNil(server.ListenAndServe())
	}()
}

func tearDown() {
	Host = originalHost
}

func mockRequest(status int, message string, errors []string) {
	statusCode = status
	resp := &PostResponse{
		Message: message,
		Errors:  errors,
	}
	r, err := json.Marshal(resp)
	results = r
	assert.MustNil(err)
}

func mockSucess() {
	mockRequest(200, "success", nil)
}

func mockError(err string) {
	mockRequest(400, "error", []string{err})
}

func TestMailSend(t *testing.T) {
	once.Do(setUp)
	defer once.Do(tearDown)
	assert = assrt.NewAssert(t)

	c := NewClient("wrong", "creds")
	args := MailArgs{
		To:       []string{"some@gmail.com"},
		From:     "brock@gmail.com",
		FromName: "brock",
		Subject:  "hello thar",
		Text:     "I am you",
	}
	mockError("Bad username / password")
	err := c.MailSend(args)
	assert.MustNotNil(err)
	assert.Equal("Status Code 400: Bad username / password", err.Error())

	c = NewClient("right", "thistime")
	mockSucess()
	err = c.MailSend(args)
	assert.MustNil(err)
}
