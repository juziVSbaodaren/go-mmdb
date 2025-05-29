package main

import (
	"flag"

	"github.com/juziVSbaodaren/gommdb/geoip"
)

// GenerateCustomMMDB generates a custom MMDB file
func main() {
	readGeoIpPath := flag.String("readGeoIpPath", "./etc/geoIpInfo.txt", "Path to input file with new IP info")
	writeGeoIpPath := flag.String("writeGeoIpPath", "./etc/CustomGeoLite.mmdb", "Path to output file")
	flag.Parse()
	if err := geoip.GenerateCustomMMDB(*writeGeoIpPath, *readGeoIpPath); err != nil {
		panic(err)
	}
}

// []*geoip.GeoIPMmdbInfo{
// 	{
// 		Ip: "103.191.243.0/24",
// 		Country: geoip.AddressMmdbParam{
// 			En: "China",
// 			Zh: "中国",
// 		},
// 		City: geoip.AddressMmdbParam{
// 			En: "Shengzhen",
// 			Zh: "深圳市",
// 		},
// 		Region: geoip.AddressMmdbParam{
// 			En: "Guangdong",
// 			Zh: "广东省",
// 		},
// 		Latitude:  0.0,
// 		Longitude: 0.0,
// 	}, {
// 		Ip: "202.201.48.0/24",
// 		Country: geoip.AddressMmdbParam{
// 			En: "China",
// 			Zh: "中国",
// 		},
// 		City: geoip.AddressMmdbParam{
// 			En: "Lanzhou",
// 			Zh: "兰州市",
// 		},
// 		Region: geoip.AddressMmdbParam{
// 			En: "Gansu",
// 			Zh: "甘肃省",
// 		},
// 		Latitude:  0.0,
// 		Longitude: 0.0,
// 	},
// }...
