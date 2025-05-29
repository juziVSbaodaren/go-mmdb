package main

import (
	"flag"
	"fmt"
	"gommdb/geoip"
)

func main() {
	mmdbPath := flag.String("mmdbPath", "./etc/GeoLite2-City.mmdb", "官方数据库路径")
	customMmdbPath := flag.String("customMmdbPath", "./etc/CustomGeoLite.mmdb", "自定义数据库路径")
	Ip := flag.String("ip", "", "ip地址")
	flag.Parse()
	geoip.InitGeoIP(*mmdbPath, *customMmdbPath)
	loc, err := geoip.GetLocation(*Ip) //103.191.243.121,183.2.172.177,202.201.48.42
	if err != nil {
		panic(err)
	}
	res := geoip.ToReadableString(fmt.Sprintf, "%s-%s", loc.Region, loc.City)

	fmt.Println(res)
}
