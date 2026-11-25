package main

import (
	"bufio"
	"geoip-mmdb/internal/dump"
	"github.com/maxmind/mmdbwriter/v2"
	"github.com/maxmind/mmdbwriter/v2/mmdbtype"
	"log"
	"net/netip"
	"os"
	"strings"
)

var (
	privateCidrs = []string{
		"0.0.0.0/8",
		"10.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"192.168.0.0/16",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/3",
		"::/127",
		"fc00::/7",
		"fe80::/10",
		"ff00::/8",
	}
)

func genCountry(lite bool) {
	dbName := "GeoLite2-merge-Country.mmdb"
	if lite {
		dbName = "country-lite.mmdb"
	}
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "GeoLite2-Country",
			RecordSize:              24,
			IncludeReservedNetworks: true,
			DisableIPv4Aliasing:     true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = dump.DumpCountry(func(prefix netip.Prefix, country *dump.Country) {
		insertCountry(country.Country.IsoCode, prefix, writer, lite)
	})
	if err != nil {
		log.Fatal(err)
	}
	InsertCountryCidrTxt("all_cn.txt", writer, lite)
	InsertCountryCidrTxt("all_cn_ipv6.txt", writer, lite)
	InsertCountryCidrTxt("china_ip_list.txt", writer, lite)
	insertPrivateCidrs(writer)
	fh, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertCountryCidrTxt(txtfile string, writer *mmdbwriter.Tree, lite bool) {
	f, err := os.Open(txtfile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(line)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		insertCountry("CN", prefix, writer, lite)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func insertPrivateCidrs(writer *mmdbwriter.Tree) {
	for _, v := range privateCidrs {
		prefix, err := netip.ParsePrefix(v)
		if err != nil {
			log.Printf("%s", err)
			continue
		}
		insertCountry("private", prefix, writer, false)
	}
}

func insertCountry(iosCode string, prefix netip.Prefix, writer *mmdbwriter.Tree, lite bool) {
	iosCode = strings.ToLower(iosCode)
	if lite && iosCode != "cn" && iosCode != "lan" && iosCode != "tw" && iosCode != "hk" && iosCode != "private" {
		return
	}
	data := mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String(iosCode),
		},
	}
	err := writer.Insert(prefix, data)
	if err != nil {
		log.Fatalf("fail to insert to writer %v\n", err)
	}
}
