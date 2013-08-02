# sendgrid

This project implements a go client for the SendGrid Web API.

Currently, I've only added what I've needed to use. Pull requests welcome.

### Usage

```go
package main

import (
	"github.com/poptip/sendgrid.go"
)

func main() {
	c := sendgrid.NewClient(username, password)
  args := sendgrid.MailArgs{
    To: []string{"someone@gmail.com"},
    From: "me@gmail.com",
    Subject: "howdy!",
    Text: "Hello there",
  }
  c.MailSend(args)
}
```

### Install

    go get github.com/poptip/sendgrid.go

### License

MIT
