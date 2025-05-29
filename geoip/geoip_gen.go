package geoip

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

type AddressMmdbParam struct {
	En string
	Zh string
}

type GeoIPMmdbInfo struct {
	Ip        string
	Country   AddressMmdbParam
	City      AddressMmdbParam
	Region    AddressMmdbParam
	Latitude  float64
	Longitude float64
}

// GenerateCustomMMDB 生成自定义 .mmdb 文件，保存到指定路径
// readGeoIpInfoPath 读取 GeoIP 信息的文件路径
// 相同ip地址的记录会被覆盖
func GenerateCustomMMDB(outputPath string, readGeoIpInfoPath string, addrInfos ...*GeoIPMmdbInfo) error {
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: "GeoLite2-City",
		RecordSize:   24,
	})
	if err != nil {
		return fmt.Errorf("failed to create mmdb writer: %w", err)
	}
	addrInfos = append(addrInfos, ReadGeoIPInfosFromFile(readGeoIpInfoPath)...)
	if len(addrInfos) == 0 {
		return fmt.Errorf("no geoip info found in file: %s", readGeoIpInfoPath)
	}

	fmt.Printf("Start to generate mmdb file,readGeoIpInfoPath:%s, len(addrInfos):%d\n", readGeoIpInfoPath, len(addrInfos))

	for _, addrInfo := range addrInfos {
		// 写入数据
		_, ipNet, _ := net.ParseCIDR(addrInfo.Ip)
		data := mmdbtype.Map{
			"country": mmdbtype.Map{
				"iso_code": mmdbtype.String("CN"),
				"names": mmdbtype.Map{
					"en":    mmdbtype.String(addrInfo.Country.En),
					"zh-CN": mmdbtype.String(addrInfo.City.Zh),
				},
			},
			"city": mmdbtype.Map{
				"names": mmdbtype.Map{
					"en":    mmdbtype.String(addrInfo.City.En),
					"zh-CN": mmdbtype.String(addrInfo.City.Zh),
				},
			},
			"subdivisions": mmdbtype.Slice{
				mmdbtype.Map{
					"names": mmdbtype.Map{
						"en":    mmdbtype.String(addrInfo.Region.En),
						"zh-CN": mmdbtype.String(addrInfo.Region.Zh),
					},
				},
			},
			"location": mmdbtype.Map{
				"latitude":  mmdbtype.Float64(addrInfo.Latitude),
				"longitude": mmdbtype.Float64(addrInfo.Longitude),
			},
		}
		if err := writer.Insert(ipNet, data); err != nil {
			return fmt.Errorf("failed to insert IP data: %w", err)
		}
	}

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 打开输出文件
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	// 写入 mmdb 数据
	if _, err := writer.WriteTo(f); err != nil {
		return fmt.Errorf("failed to write mmdb data: %w", err)
	}

	return nil
}

// AppendGeoIPInfoToFile 追加 GeoIP 信息到文件
func AppendGeoIPInfoToFile(filePath string, info GeoIPMmdbInfo) error {
	line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%.6f,%.6f\n",
		info.Ip,
		info.Country.En, info.Country.Zh,
		info.Region.En, info.Region.Zh,
		info.City.En, info.City.Zh,
		info.Latitude, info.Longitude,
	)

	// 自动创建文件及目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line)
	return err
}

// ReadGeoIPInfosFromFile 从文件中读取 GeoIP 信息
// eg: Ip, Country.En, Country.Zh, Region.En, Region.Zh, City.En, City.Zh, Latitude, Longitude
// eg: 172.16.58.0/24,China,中国,Guangdong,广东,Shenzhen,深圳,23.1,113.3
func ReadGeoIPInfosFromFile(filePath string) []*GeoIPMmdbInfo {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var infos []*GeoIPMmdbInfo
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) != 9 {
			continue // 跳过格式错误的行
		}

		lat, _ := strconv.ParseFloat(fields[7], 64)
		lon, _ := strconv.ParseFloat(fields[8], 64)

		info := &GeoIPMmdbInfo{
			Ip:        fields[0],
			Country:   AddressMmdbParam{En: fields[1], Zh: fields[2]},
			Region:    AddressMmdbParam{En: fields[3], Zh: fields[4]},
			City:      AddressMmdbParam{En: fields[5], Zh: fields[6]},
			Latitude:  lat,
			Longitude: lon,
		}
		infos = append(infos, info)
	}

	if err := scanner.Err(); err != nil {
		return nil
	}
	return infos
}
