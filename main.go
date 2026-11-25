package main

import (
	"flag"
	"geoip-mmdb/internal/dump"
	"geoip-mmdb/internal/utils"
	"log"
)

var (
	username string
	password string
)

func init() {
	flag.StringVar(&username, "u", "", "")
	flag.StringVar(&password, "p", "", "")
	flag.Parse()
}

func main() {
	if username == "" || password == "" {
		log.Fatalln("username or password cat not be empty")
	}
	GenCity(username, password)

	//err := dump.InitReader("GeoLite2-City.mmdb", "GeoLite2-Country.mmdb", "GeoLite2-ASN.mmdb")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	GenCountry()
	genasn()
}

func GenCountry() {
	downloadf("https://ispip.clang.cn/all_cn.txt", "all_cn.txt")
	downloadf("https://ispip.clang.cn/all_cn_ipv6.txt", "all_cn_ipv6.txt")
	downloadf("https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt", "china_ip_list.txt")
	genCountry(false)
	genCountry(true)
}

func GenCity(username string, password string) {
	err := utils.DownloadFile("https://download.maxmind.com/geoip/databases/GeoLite2-ASN/download?suffix=tar.gz",
		"GeoLite2-ASN.tar.gz", username, password)
	if err != nil {
		log.Fatalln(err)
	}

	err = utils.DownloadFile("https://download.maxmind.com/geoip/databases/GeoLite2-City/download?suffix=tar.gz",
		"GeoLite2-City.tar.gz", username, password)
	if err != nil {
		log.Fatalln(err)
	}

	err = utils.DownloadFile("https://download.maxmind.com/geoip/databases/GeoLite2-Country/download?suffix=tar.gz",
		"GeoLite2-Country.tar.gz", username, password)
	if err != nil {
		log.Fatalln(err)
	}

	err = utils.UntarGz("GeoLite2-ASN.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	err = utils.UntarGz("GeoLite2-City.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}
	err = utils.UntarGz("GeoLite2-Country.tar.gz")
	if err != nil {
		log.Fatalln(err)
	}

	err = dump.InitReader("GeoLite2-City.mmdb", "GeoLite2-Country.mmdb", "GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatalln(err)
	}

	downloadf("https://raw.githubusercontent.com/lionsoul2014/ip2region/refs/heads/master/data/ipv4_source.txt", "ipv4_source.txt")
	downloadf("https://raw.githubusercontent.com/lionsoul2014/ip2region/refs/heads/master/data/ipv6_source.txt", "ipv6_source.txt")

	genCityDB(false)
	genCityDB(true)
}

func downloadf(url, out string) {
	err := utils.DownloadFile(url, out, "", "")
	if err != nil {
		log.Fatalln(err)
	}
}
