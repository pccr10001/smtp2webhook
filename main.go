package main

import (
	"bytes"
	"github.com/DusanKasan/parsemail"
	"github.com/google/uuid"
	"github.com/mhale/smtpd"
	"golang.org/x/net/html/charset"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/url"
	"strings"
)

var accounts []Account

func decode(text string) string {
	CharsetReader := func(label string, input io.Reader) (io.Reader, error) {
		label = strings.Replace(label, "windows-", "cp", -1)
		encoding, _ := charset.Lookup(label)
		return encoding.NewDecoder().Reader(input), nil
	}
	dec := mime.WordDecoder{CharsetReader: CharsetReader}

	header, err := dec.DecodeHeader(text)

	if err != nil {
		return text
	}
	return header
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := parsemail.Parse(bytes.NewReader(data))
	subject := decode(msg.Header.Get("Subject"))

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
					break
				}
			}
		}
	}

	if account == nil {
		return nil
	}

	uri := account.Webhook.Host
	if account.IsTest {
		uri += "/webhook-test/"
	} else {
		uri += "/webhook/"
	}
	uri += account.Webhook.Id

	go gentleman.New().
		Post().
		URL(uri).
		AddQuery("subject", url.QueryEscape(decode(subject))).
		AddHeader("Content-Type", "text/plain").
		BodyString(msg.HTMLBody).
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
