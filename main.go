package main

import (
	"bytes"
	"github.com/cention-sany/go.enmime"
	"github.com/cention-sany/net/mail"
	"github.com/google/uuid"
	"github.com/mhale/smtpd"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"strings"
)

var accounts []Account

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	mime, _ := enmime.ParseMIMEBody(msg)
	subject := mime.GetHeader("Subject")

	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)

	t := strings.Split(to[0], "@")
	if len(t) != 2 {
		return nil
	}

	var account *Account

	if _, err := uuid.Parse(to[0]); err == nil {
		for _, a := range accounts {
			if t[0] == a.Webhook.Id {
				account = &a
				break
			}
		}
	} else {
		for _, a := range accounts {
			for _, alias := range a.Alias {
				if t[0] == alias && t[1] == a.Host {
					account = &a
					goto found
				}
			}
		}
	}

found:
	if account == nil {
		return nil
	}

	go gentleman.New().
		Post().
		URL(account.Webhook.Host+"/webhook-test/"+account.Webhook.Id).
		AddQuery("subject", url.QueryEscape(subject)).
		AddHeader("Content-Type", "text/plain").
		BodyString(mime.HTML).
		Do()

	go gentleman.New().
		Post().
		URL(account.Webhook.Host+"/webhook/"+account.Webhook.Id).
		AddQuery("subject", url.QueryEscape(subject)).
		AddHeader("Content-Type", "text/plain").
		BodyString(mime.HTML).
		Do()

	log.Printf("Webhook emitted, %s %s", account.Webhook.Host, account.Webhook.Id)

	return nil
}

func main() {

	conf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Panicln("Failed to open config.yaml, err=", err)
	}
	err = yaml.Unmarshal(conf, &accounts)
	if err != nil {
		log.Panicln("Failed to parse config.yaml, err=", err)
	}

	log.Printf("SMTP server is listening on :2525")
	err = smtpd.ListenAndServe("0.0.0.0:2525", mailHandler, "Postfix", "postfix")
	if err != nil {
		log.Panicln(err)
	}
}
