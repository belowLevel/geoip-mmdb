package main

type City struct {
	Country        string  `maxminddb:"country"`
	Province       string  `maxminddb:"province"`
	City           string  `maxminddb:"city"`
	District       string  `maxminddb:"district"`
	IsoCode        string  `maxminddb:"iso_code"`
	AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	Latitude       float64 `maxminddb:"latitude"`
	Longitude      float64 `maxminddb:"longitude"`
}
