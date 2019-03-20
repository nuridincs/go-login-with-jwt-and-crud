package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

//Request struct
type Request struct {
	from      string
	fromAlias string
	to        []string
	subject   string
	body      string
}

//MIME struct
const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

//NewRequest struct
func NewRequest(from string, fromAlias string, to []string, subject string) *Request {
	return &Request{
		from:      from,
		fromAlias: fromAlias,
		to:        to,
		subject:   subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMail() bool {
	body := "From: " + r.fromAlias + "\r\nTo: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", "smtp.gmail.com", 587)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", "email", "password", "smtp.gmail.com"), r.from, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

//Send struct
func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", r.to)
	} else {
		log.Printf("Failed to send the email to %s\n", r.to)
	}
}
