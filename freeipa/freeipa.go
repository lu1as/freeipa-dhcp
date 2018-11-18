package freeipa

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type FreeIPAClient struct {
	url    string
	client *http.Client
	cookie *http.Cookie
}

func NewFreeIPAClient(serverUrl string) *FreeIPAClient {
	return &FreeIPAClient{
		url:    serverUrl,
		client: &http.Client{},
	}
}

func (f *FreeIPAClient) Login(username string, password string) error {
	d := []byte("user=" + username + "&password=" + password)
	req, err := http.NewRequest("POST", f.url+"/ipa/session/login_password", bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/plain")

	res, err := f.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("login failed")
	} else if len(res.Cookies()) < 1 {
		return fmt.Errorf("no login cookie")
	}
	f.cookie = res.Cookies()[0]
	return nil
}

func (f *FreeIPAClient) AllowInsecure() {
	f.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (f *FreeIPAClient) request(data []byte) ([]byte, error) {
	req, _ := http.NewRequest("POST", f.url+"/ipa/session/json", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", f.url+"/ipa")
	req.AddCookie(f.cookie)

	res, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	f.cookie = res.Cookies()[0]
	return ioutil.ReadAll(res.Body)
}
