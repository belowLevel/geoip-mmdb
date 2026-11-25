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
	urlmap = map[string]string{
		"chinatelecom.txt":      "电信",
		"chinatelecom_ipv6.txt": "电信",
		"unicom_cnc.txt":        "联通",
		"unicom_cnc_ipv6.txt":   "联通",
		"cmcc.txt":              "移动",
		"cmcc_ipv6.txt":         "移动",
		"cernet.txt":            "教育网",
		"cernet_ipv6.txt":       "教育网",
		"gwbn.txt":              "电信通/长城宽带/鹏博士",
		"gwbn_ipv6.txt":         "电信通/长城宽带/鹏博士",
		"chinabtn.txt":          "广电",
		"chinabtn_ipv6.txt":     "广电",
	}
)

func genasn() {
	dbName := "GeoLite2-merge-ASN.mmdb"
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "GeoLite2-ASN",
			RecordSize:              24,
			IncludeReservedNetworks: true,
			DisableIPv4Aliasing:     true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	err = dump.DumpASN(func(prefix netip.Prefix, asn *dump.ASN) {
		insertASN(prefix, asn.AutonomousSystemOrganization, writer)
	})
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range urlmap {
		downloadf("https://ispip.clang.cn/"+k, k)
		f, err := os.Open(k)
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
			insertASN(prefix, v, writer)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

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

func insertASN(prefix netip.Prefix, autonomousSystemOrganization string, writer *mmdbwriter.Tree) {
	data := mmdbtype.Map{
		"autonomous_system_organization": mmdbtype.String(autonomousSystemOrganization),
	}
	err := writer.Insert(prefix, data)
	if err != nil {
		log.Fatalf("fail to insert to writer %v\n", err)
	}
}
