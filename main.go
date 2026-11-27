package main

import (
	"flag"
	"geoip-mmdb/internal/dump"
	"geoip-mmdb/internal/utils"
	"github.com/maxmind/mmdbwriter/v2"
	"github.com/maxmind/mmdbwriter/v2/mmdbtype"
	"github.com/oschwald/maxminddb-golang/v2"
	"log"
	"net/netip"
	"os"
)

var (
	username    string
	password    string
	liteCityASN bool
)

func init() {
	flag.StringVar(&username, "u", "", "")
	flag.StringVar(&password, "p", "", "")
	flag.BoolVar(&liteCityASN, "l", false, "")
	flag.Parse()
}

func main() {
	if liteCityASN {
		genCityAsnLite()
		return
	}
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

func DumpCity(action func(prefix netip.Prefix, city *City)) error {
	cityReader, err := maxminddb.Open("GeoLite2-City-ASN.mmdb")
	if err != nil {
		return err
	}

	for result := range cityReader.Networks() {
		record := City{}
		err := result.Decode(&record)
		if err != nil {
			continue
		}
		action(result.Prefix(), &record)
	}
	return nil
}

func genCityAsnLite() {
	cityWriter, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "GeoLite2-City",
			RecordSize:              28,
			IncludeReservedNetworks: true,
			DisableIPv4Aliasing:     true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = DumpCity(func(prefix netip.Prefix, city *City) {
		if city.IsoCode != "CN" && city.IsoCode != "private" {
			return
		}

		data := mmdbtype.Map{
			"country":         mmdbtype.String(city.Country),
			"province":        mmdbtype.String(city.Province),
			"city":            mmdbtype.String(city.City),
			"district":        mmdbtype.String(city.District),
			"iso_code":        mmdbtype.String(city.IsoCode),
			"organization":    mmdbtype.String(city.Organization),
			"accuracy_radius": mmdbtype.Uint16(city.AccuracyRadius),
			"latitude":        mmdbtype.Float64(city.Latitude),
			"longitude":       mmdbtype.Float64(city.Longitude),
		}
		err := cityWriter.Insert(prefix, data)
		if err != nil {
			log.Fatalln(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	fh, err := os.Create("GeoLite2-City-ASN-lite.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = cityWriter.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
