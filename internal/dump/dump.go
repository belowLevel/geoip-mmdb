package dump

import (
	"errors"
	"github.com/oschwald/maxminddb-golang/v2"
	"net/netip"
)

var (
	cityReader        *maxminddb.Reader
	countryReader     *maxminddb.Reader
	asnReader         *maxminddb.Reader
	ErrDbNotAvailable = errors.New("DB_NOT_AVAILABLE")
)

type City struct {
	City struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`

	Country struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`

	Location struct {
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`

	Subdivisions []struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
}

type Country struct {
	Country struct {
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		IsoCode           string            `maxminddb:"iso_code"`
		Names             map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
}

type ASN struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}

func DumpCountry(action func(prefix netip.Prefix, country *Country)) error {
	if countryReader == nil {
		return ErrDbNotAvailable
	}

	for result := range countryReader.Networks() {
		record := Country{}
		err := result.Decode(&record)
		if err != nil {
			continue
		}
		action(result.Prefix(), &record)
	}
	return nil
}

func DumpCity(action func(prefix netip.Prefix, city *City)) error {
	if cityReader == nil {
		return ErrDbNotAvailable
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

func SearchCity(ip netip.Addr) (city *City, err error) {
	if cityReader == nil {
		return nil, ErrDbNotAvailable
	}
	record := &City{}
	err = cityReader.Lookup(ip).Decode(record)
	return record, err
}

func DumpASN(action func(prefix netip.Prefix, asn *ASN)) error {
	if asnReader == nil {
		return ErrDbNotAvailable
	}

	for result := range asnReader.Networks() {
		record := ASN{}
		err := result.Decode(&record)
		if err != nil {
			continue
		}
		action(result.Prefix(), &record)
	}
	return nil
}

func InitReader(cityMMDB, countryMMDB, asnMMDB string) error {
	var err error
	cityReader, err = maxminddb.Open(cityMMDB)
	if err != nil {
		return err
	}
	countryReader, err = maxminddb.Open(countryMMDB)
	if err != nil {
		return err
	}
	asnReader, err = maxminddb.Open(asnMMDB)
	if err != nil {
		return err
	}
	return nil
}

func TrimLocationStr(location *string) {
	if location == nil || *location == "" {
		return
	}
	pr := []rune(*location)
	prl := len(pr)
	if prl > 0 {
		if string(pr[prl-1]) == "省" {
			*location = string(pr[0 : prl-1])
		}
		if string(pr[prl-1]) == "市" {
			*location = string(pr[0 : prl-1])
		}
		if prl > 3 {
			if string(pr[prl-3:]) == "自治区" {
				*location = string(pr[0 : prl-3])
			}
		}
	}
}

func ExpandCity(cityData *City) (isoCode, country, province, city, district string, longitude, latitude float64, accuracyRadius uint16) {
	country = cityData.Country.Names["zh-CN"]
	city = cityData.City.Names["zh-CN"]
	for i, v := range cityData.Subdivisions {
		switch i {
		case 0:
			province = v.Names["zh-CN"]
		case 1:
			district = v.Names["zh-CN"]
		default:
			break
		}
	}
	isoCode = cityData.Country.IsoCode
	if isoCode == "CN" {
		TrimLocationStr(&province)
		TrimLocationStr(&city)
	}
	longitude = cityData.Location.Longitude
	latitude = cityData.Location.Latitude
	accuracyRadius = cityData.Location.AccuracyRadius
	return
}
