package generate_country_mmdb

import (
	"bufio"
	"encoding/csv"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var (
	writer        *mmdbwriter.Tree
	CountryBlocks = make(map[string]*CountryBlock)
	Locations     = make(map[string]*Location)
	cnCidrs       = make([]string, 0)
)

type CountryBlock struct {
	NetWork   string
	GeoNameId string
}
type Location struct {
	GeonameId      string
	CountryIsoCode string
}

func ReadCountryBlockCsvFile(filename string) {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.Read()
	for {
		strs, err := reader.Read()
		if err == io.EOF {
			break
		}
		countryBlock := CountryBlock{
			NetWork:   strs[0],
			GeoNameId: strs[1],
		}
		CountryBlocks[countryBlock.NetWork] = &countryBlock
	}
}

func ReadLocationBlockCsvFile(filename string) {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	reader.Read()
	for {
		strs, err := reader.Read()
		if err == io.EOF {
			break
		}
		location := Location{
			GeonameId:      strs[0],
			CountryIsoCode: strs[4],
		}
		Locations[location.GeonameId] = &location
	}
}

func ReadCidrTxt(txtfile string) {
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
		cnCidrs = append(cnCidrs, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Generatemmdb(pathbase string) {
	readFiles(pathbase)

	var err error
	writer, err = mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "GeoLite2-Country",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	insertCsvCidrs()

	insertCnCidrs()

	fh, err := os.Create("../GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func insertCsvCidrs() {
	for k, v := range CountryBlocks {
		_, ipnet, err := net.ParseCIDR(k)
		if err != nil {
			log.Panic(err)
		}
		location := getLocation(v, k)

		insertData(location.CountryIsoCode, ipnet)
	}
}

func getLocation(v *CountryBlock, k string) *Location {
	location, ok := Locations[v.GeoNameId]
	if !ok {
		log.Printf("country location not found %s", k)
		location = &Location{
			GeonameId:      "",
			CountryIsoCode: "",
		}
	}
	return location
}

func insertData(iosCode string, ipnet *net.IPNet) {
	data := mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String(iosCode),
		},
		"registered_country": mmdbtype.Map{
			"iso_code": mmdbtype.String(iosCode),
		},
	}
	err := writer.Insert(ipnet, data)
	if err != nil {
		log.Fatalf("fail to insert to writer %v\n", err)
	}
}

func insertCnCidrs() {
	for _, v := range cnCidrs {
		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			log.Panicf("%v %v", err, v)
		}
		insertData("CN", ipnet)
	}
}

func readFiles(pathBase string) {
	ReadCountryBlockCsvFile(filepath.Join(pathBase, "GeoLite2-Country-Blocks-IPv4.csv"))
	ReadCountryBlockCsvFile(filepath.Join(pathBase, "GeoLite2-Country-Blocks-IPv6.csv"))
	ReadLocationBlockCsvFile(filepath.Join(pathBase, "GeoLite2-Country-Locations-zh-CN.csv"))
	ReadCidrTxt(filepath.Join(pathBase, "all_cn.txt"))
	ReadCidrTxt(filepath.Join(pathBase, "all_cn_ipv6.txt"))
	ReadCidrTxt(filepath.Join(pathBase, "china_ip_list.txt"))
}
