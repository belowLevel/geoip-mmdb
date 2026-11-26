package main

import (
	"fmt"
	"geoip-mmdb/internal/dump"
	"geoip-mmdb/internal/utils"
	"github.com/maxmind/mmdbwriter/v2"
	"github.com/maxmind/mmdbwriter/v2/mmdbtype"
	"go4.org/netipx"
	"log"
	"net/netip"
	"os"
	"strings"
)

func genCityDB(lite bool) {
	dbName := "GeoLite2-merge-City.mmdb"
	if lite {
		dbName = "GeoLite2-merge-City-lite.mmdb"
	}
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

	err = dump.DumpCity(func(prefix netip.Prefix, cityData *dump.City) {
		if lite && cityData.Country.IsoCode != "CN" {
			return
		}
		isoCode, country, province, city, district, longitude, latitude, accuracyRadius := dump.ExpandCity(cityData)
		newCity := City{
			Country:        country,
			Province:       province,
			City:           city,
			District:       district,
			IsoCode:        isoCode,
			AccuracyRadius: accuracyRadius,
			Latitude:       latitude,
			Longitude:      longitude,
		}
		insertCityData(prefix, newCity, cityWriter)
	})
	if err != nil {
		log.Fatal(err)
	}
	mergeip2region(cityWriter, "ipv4_source.txt", lite)
	mergeip2region(cityWriter, "ipv6_source.txt", lite)
	fh, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = cityWriter.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func mergeip2region(cityWriter *mmdbwriter.Tree, filename string, lite bool) {

	utils.ReadTxt(filename, func(line string) {
		d := strings.Split(line, "|")
		if len(d) != 7 {
			log.Fatalln("ip2region txt is not valid")
		}

		prefixes, err := RangeToCIDR(d[0], d[1])
		if err != nil {
			log.Fatalln(err)
		}
		if lite && d[6] != "CN" && d[6] != "0" {
			return
		}
		newCity := City{
			Country:        d[2],
			Province:       d[3],
			City:           d[4],
			District:       "",
			IsoCode:        d[6],
			AccuracyRadius: 0,
			Latitude:       0,
			Longitude:      0,
		}
		for _, prefix := range prefixes {
			c, err := dump.SearchCity(prefix.Addr())
			if err == nil && c != nil {
				_, _, _, _, district, longitude, latitude, accuracyRadius := dump.ExpandCity(c)
				if district != "" {
					newCity.District = district
				}
				newCity.AccuracyRadius = accuracyRadius
				newCity.Longitude = longitude
				newCity.Latitude = latitude
			}
			insertCityData(prefix, newCity, cityWriter)
		}
	})
}

func RangeToCIDR(start, end string) ([]netip.Prefix, error) {
	first, err := netip.ParseAddr(start)
	if err != nil {
		return nil, err
	}

	last, err := netip.ParseAddr(end)
	if err != nil {
		return nil, err
	}

	if first.Is4() != last.Is4() {
		return nil, fmt.Errorf("IP address families do not match")
	}

	if first.Compare(last) > 0 {
		return nil, fmt.Errorf("start IP > end IP")
	}

	r := netipx.IPRangeFrom(first, last)
	if !r.IsValid() {
		return nil, fmt.Errorf("invalid IP range")
	}

	return r.Prefixes(), nil
}

func insertCityData(prefix netip.Prefix, city City, writer *mmdbwriter.Tree) {
	dump.TrimLocationStr(&city.Province)
	dump.TrimLocationStr(&city.City)
	if city.IsoCode == "0" {
		city.IsoCode = "private"
		city.Country = ""
		city.Province = ""
		city.City = ""
		city.District = ""
	}
	data := mmdbtype.Map{
		"country":         mmdbtype.String(city.Country),
		"province":        mmdbtype.String(city.Province),
		"city":            mmdbtype.String(city.City),
		"district":        mmdbtype.String(city.District),
		"iso_code":        mmdbtype.String(city.IsoCode),
		"accuracy_radius": mmdbtype.Uint16(city.AccuracyRadius),
		"latitude":        mmdbtype.Float64(city.Latitude),
		"longitude":       mmdbtype.Float64(city.Longitude),
	}
	err := writer.Insert(prefix, data)
	if err != nil {
		log.Fatalf("fail to insert to writer %v\n", err)
	}
}
