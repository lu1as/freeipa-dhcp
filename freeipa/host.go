package freeipa

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type hostResult struct {
	Result struct {
		Result []Host `json:"result"`
	} `json:"result"`
}

type Host struct {
	Fqdn []string `json:"fqdn"`
	MAC  []string `json:"macaddress"`
}

func (f *FreeIPAClient) GetHosts() ([]Host, error) {
	j := []byte("{\"method\":\"host_find\", \"params\":[[\"\"],{}]}")
	d, err := f.request(j)
	log.Debugf("host_find result: %s", string(d))
	if err != nil {
		return nil, err
	}

	var r hostResult
	err = json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return r.Result.Result, nil
}

func (f *FreeIPAClient) GetHost(fqdn string) (*Host, error) {
	j := []byte("{\"method\":\"host_find\", \"params\":[[\"\"],{\"fqdn\": \"" + fqdn + "\" }]}")
	d, err := f.request(j)
	log.Debugf("host_find %s result: %s", fqdn, string(d))
	if err != nil {
		return nil, err
	}

	var r hostResult
	err = json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	} else if len(r.Result.Result) < 1 {
		return nil, fmt.Errorf("host not found")
	}

	return &r.Result.Result[0], nil
}
