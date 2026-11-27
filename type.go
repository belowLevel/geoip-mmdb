package main

type City struct {
	Country        string  `maxminddb:"country"`
	Province       string  `maxminddb:"province"`
	City           string  `maxminddb:"city"`
	District       string  `maxminddb:"district"`
	IsoCode        string  `maxminddb:"iso_code"`
	Organization   string  `maxminddb:"organization"`
	AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	Latitude       float64 `maxminddb:"latitude"`
	Longitude      float64 `maxminddb:"longitude"`
}
