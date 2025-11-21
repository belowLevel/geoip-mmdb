package generate_asn_mmdb

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
	"strconv"
	"strings"
)

var (
	writer    *mmdbwriter.Tree
	AsnBlocks = make([]*AsnBlock, 0)
)

type AsnBlock struct {
	Network                      string
	AutonomousSystemOrganization string
}

func ReadAsnBlockCsvFile(filename string) {
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
		asnblock := AsnBlock{
			Network:                      strs[0],
			AutonomousSystemOrganization: strs[2],
		}
		AsnBlocks = append(AsnBlocks, &asnblock)
	}
}

func ReadCidrTxt(txtfile string, asnName string) {
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
		asnblock := AsnBlock{
			Network:                      line,
			AutonomousSystemOrganization: asnName,
		}
		AsnBlocks = append(AsnBlocks, &asnblock)

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
			DatabaseType: "GeoLite2-ASN",
			RecordSize:   24,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range AsnBlocks {
		_, ipnet, err := net.ParseCIDR(v.Network)
		if err != nil {
			log.Panicf("%v %v", err, v)
		}
		insertData(v, ipnet)
	}
	fh, err := os.Create("../GeoLite2-ASN.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

func readFiles(pathBase string) {
	ReadAsnBlockCsvFile(filepath.Join(pathBase, "GeoLite2-ASN-Blocks-IPv4.csv"))
	ReadAsnBlockCsvFile(filepath.Join(pathBase, "GeoLite2-ASN-Blocks-IPv6.csv"))
	ReadCidrTxt(filepath.Join(pathBase, "cernet.txt"), "中国教育网")
	ReadCidrTxt(filepath.Join(pathBase, "cernet_ipv6.txt"), "中国教育网")
	ReadCidrTxt(filepath.Join(pathBase, "chinatelecom.txt"), "中国电信")
	ReadCidrTxt(filepath.Join(pathBase, "chinatelecom_ipv6.txt"), "中国电信")
	ReadCidrTxt(filepath.Join(pathBase, "cmcc.txt"), "中国移动")
	ReadCidrTxt(filepath.Join(pathBase, "cmcc_ipv6.txt"), "中国移动")
	//ReadCidrTxt(filepath.Join(pathBase, "crtc.txt"), "中国铁通")
	//ReadCidrTxt(filepath.Join(pathBase, "crtc_ipv6.txt"), "中国铁通")
	ReadCidrTxt(filepath.Join(pathBase, "gwbn.txt"), "电信通/长城宽带/鹏博士")
	ReadCidrTxt(filepath.Join(pathBase, "gwbn_ipv6.txt"), "电信通/长城宽带/鹏博士")
	ReadCidrTxt(filepath.Join(pathBase, "unicom_cnc.txt"), "中国联通")
	ReadCidrTxt(filepath.Join(pathBase, "unicom_cnc_ipv6.txt"), "中国联通")
	ReadCidrTxt(filepath.Join(pathBase, "chinabtn.txt"), "中国广电")
	ReadCidrTxt(filepath.Join(pathBase, "chinabtn_ipv6.txt"), "中国广电")
}

func insertData(asn *AsnBlock, ipnet *net.IPNet) {
	data := mmdbtype.Map{
		"autonomous_system_organization": mmdbtype.String(asn.AutonomousSystemOrganization),
	}
	err := writer.Insert(ipnet, data)
	if err != nil {
		log.Fatalf("fail to insert to writer %v\n", err)
	}
}

func str2int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func str2float64(str string) float64 {
	f, err := strconv.ParseFloat(str, 65)
	if err != nil {
		return 0.0
	}
	return f
}

func getaIp(cidr string) string {
	idx := strings.LastIndex(cidr, "/")
	return cidr[:idx]
}
