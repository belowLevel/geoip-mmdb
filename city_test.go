package main

import (
	"github.com/oschwald/maxminddb-golang/v2"
	"net/netip"
	"testing"
)

func TestSearch(t *testing.T) {
	cityReader, err := maxminddb.Open("../GeoLite2-merge-City-lite.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	record := &City{}
	ipaddr, err := netip.ParseAddr("114.114.114.114")
	if err != nil {
		t.Fatal(err)
	}
	err = cityReader.Lookup(ipaddr).Decode(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", record)
}
