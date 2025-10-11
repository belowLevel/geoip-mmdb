package generate_test

import (
	"fmt"
	"geoip-mmdb/generate_asn_mmdb"
	"geoip-mmdb/generate_city_mmdb"
	"geoip-mmdb/generate_country_mmdb"
	"geoip-mmdb/merge_ip2region"
	"geoip-mmdb/merge_pure_city_mmdb"
	"geoip-mmdb/pureip"
	"geoip-mmdb/reader"
	"github.com/orestonce/Ip2regionTool"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"os"
	"testing"
)

func Test_Generate_Asn(t *testing.T) {
	generate_asn_mmdb.Generatemmdb("../generate_asn_mmdb")
}

func Test_Generate_City(t *testing.T) {
	generate_city_mmdb.Generatemmdb("../generate_city_mmdb")
}

func Test_Generate_Country(t *testing.T) {
	generate_country_mmdb.Generatemmdb("../generate_country_mmdb")
}

func TestConvert(t *testing.T) {
	pureip.ConvertTxt2Csv("../")
}

func Test_Generate_Merge_City(t *testing.T) {
	merge_pure_city_mmdb.Generatemmdb("../merge_pure_city_mmdb")
}

//	func Test_Generate_Merge_City2(t *testing.T) {
//		merge_pure_city_mmdb.Generatemmdb2("../merge_pure_city_mmdb")
//	}
func TestReader(t *testing.T) {
	cityReader, err := reader.Open("../GeoLite2-City.mmdb")
	require.NoError(t, err)
	defer cityReader.Close()
	asnReader, err := reader.Open("../GeoLite2-ASN.mmdb")
	require.NoError(t, err)
	defer asnReader.Close()

	cidr := ""
	cityRecord, err := cityReader.City(net.ParseIP(cidr))
	require.NoError(t, err)

	asnRecord, err := asnReader.ASN(net.ParseIP(cidr))
	require.NoError(t, err)

	log.Printf("%v %v", cityRecord, asnRecord)
}

func TestReadContry(t *testing.T) {
	countryReader, err := reader.Open("../GeoLite2-Country.mmdb")
	require.NoError(t, err)
	defer countryReader.Close()

	log.Printf("%v", countryReader.Metadata())
	cidr := "101.33.160.0"
	cityRecord, err := countryReader.Country(net.ParseIP(cidr))
	require.NoError(t, err)

	log.Printf("%v", cityRecord)
}

func TestConvertDbToTxt(t *testing.T) {
	var dbPath = "ip2region_v4.xdb"
	var txtFileName = "1.txt"
	errMsg := merge_ip2region.ConvertDbToTxt(Ip2regionTool.ConvertDbToTxt_Req{
		DbFileName:  dbPath,
		TxtFileName: txtFileName,
		Merge:       true,
		DbVersion:   2,
	})
	if errMsg != `` {
		fmt.Println(errMsg)
		os.Exit(-1)
	}
}
