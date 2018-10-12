package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/iampigeon/pigeon"
	"github.com/iampigeon/pigeon/backend"
)

type service struct{}

func (s *service) Approve(content []byte) (valid bool, err error) {
	if content == nil {
		return false, errors.New("Invalid message content")
	}

	fmt.Println(string(content))
	m := new(pigeon.HTTP)

	err = json.Unmarshal(content, m)
	if err != nil {
		return false, err
	}

	// validate topic to avoid  breakout
	fmt.Println(m)

	return true, nil
}

func (s *service) Deliver(content []byte) error {
	// Parse content HTTP to struct
	m := new(pigeon.HTTP)
	err := json.Unmarshal(content, m)
	if err != nil {
		return err
	}

	// TODO(ca): secure this
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Prepare client
	client := &http.Client{Transport: tr}

	// prepare request to target
	req, err := http.NewRequest("POST", m.URL.String(), bytes.NewBufferString(m.Body))
	if err != nil {
		return err
	}

	// Set all headers
	for key, value := range m.Headers {
		req.Header.Set(key, value.(string))
	}

	// Execute POST to target
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	// TODO(ca): check to use the response to send good or bad status to pigeon-central

	fmt.Println("camilito la hizo pigeon-http")

	log.Printf("message received pigeon-http: %s", content)
	return nil
}

func main() {
	host := flag.String("host", "", "host of the service")
	port := flag.Int("port", 9020, "host of the service")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("Serving at %s", addr)

	svc := &service{}

	if err := backend.ListenAndServe(pigeon.NetAddr(addr), svc); err != nil {
		log.Fatal(err)
	}
}
