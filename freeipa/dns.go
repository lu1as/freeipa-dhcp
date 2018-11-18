package freeipa

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type dnsRecordResult struct {
	Result struct {
		Result []DNSRecord `json:"result"`
	} `json:"result"`
}

type dnsZoneResult struct {
	Result struct {
		Result []DNSZone `json:"result"`
	} `json:"result"`
}

type DNSRecord struct {
	IDNSName    []string `json:"idnsname"`
	ARecord     []string `json:"arecord"`
	AAARecord   []string `json:"aaarecord"`
	NSRecord    []string `json:"nsrecord"`
	TXTRecord   []string `json:"txtrecord"`
	CNAMERecord []string `json:"cnamerecord`
}

type DNSZone struct {
	Idnsname string `json:"idnsname"`
}

func (f *FreeIPAClient) GetDNSRecords(zone string) ([]DNSRecord, error) {
	j := []byte("{\"method\": \"dnsrecord_find\", \"params\": [[\"" + zone + "\", \"\"], {}]}")
	d, err := f.request(j)
	log.Debugf("dnsrecord_find in zone %s result: %s", zone, string(d))
	if err != nil {
		return nil, err
	}

	var r dnsRecordResult
	err = json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return r.Result.Result, nil
}

func (f *FreeIPAClient) GetDNSZones() ([]DNSZone, error) {
	j := []byte("{\"method\": \"dnszone_find\", \"params\":[[\"\"], {}]}")
	d, err := f.request(j)
	log.Debugf("dnszone_find result: %s", string(d))
	if err != nil {
		return nil, err
	}

	var r dnsZoneResult
	err = json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return r.Result.Result, nil
}

func (f *FreeIPAClient) GetDNSRecord(zone string, name string) (*DNSRecord, error) {
	j := []byte("{\"method\": \"dnsrecord_find\", \"params\": [[\"" + zone + "\", \"\"], {\"idnsname\": \"" + name + "\"}]}")
	d, err := f.request(j)
	log.Debugf("dnsrecord_find %s in zone %s result: %s", name, zone, string(d))
	if err != nil {
		return nil, err
	}

	var r dnsRecordResult
	err = json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	} else if len(r.Result.Result) < 1 {
		return nil, fmt.Errorf("dns record not found")
	}

	return &r.Result.Result[0], nil
}
