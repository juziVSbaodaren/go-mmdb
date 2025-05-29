package geoip

import (
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var (
	db       *geoip2.Reader // 官方数据库
	customDb *geoip2.Reader // 自定义数据库
	once     sync.Once
	err      error
)

// 自定义错误
var (
	ErrDBNotInitialized = &GeoIPError{"GeoIP 数据库未初始化，请调用 InitGeoIP"}
	ErrInvalidIP        = &GeoIPError{"无效的 IP 地址"}
)

type GeoIPError struct {
	Msg string
}

func (e *GeoIPError) Error() string {
	return e.Msg
}

// Location 表示 IP 地理信息
type Location struct {
	Country string
	Region  string
	City    string
	Lat     float64
	Lon     float64
}

// InitGeoIP 初始化数据库（只执行一次）
func InitGeoIP(dbPath, customDbPath string) error {
	once.Do(func() {
		db, err = geoip2.Open(dbPath)
		customDb, err = geoip2.Open(customDbPath)
	})
	return err
}

// GetLocation 返回 IP 对应的地理位置信息（中文）
func GetLocation(ipStr string) (*Location, error) {
	if db == nil {
		return nil, ErrDBNotInitialized
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, ErrInvalidIP
	}

	record, err := db.City(ip) // 官方库使用 City 方法获取位置信息
	if err != nil {
		return nil, err
	}
	if record.City.Names["zh-CN"] == "" && customDb != nil {
		// 请求自定义库
		record, err = customDb.City(ip)
		if err != nil {
			return nil, err
		}
	}
	loc := &Location{
		Country: record.Country.Names["zh-CN"],
		Lat:     record.Location.Latitude,
		Lon:     record.Location.Longitude,
		City:    record.City.Names["zh-CN"],
	}

	if len(record.Subdivisions) > 0 {
		loc.Region = record.Subdivisions[0].Names["zh-CN"]
	}

	return loc, nil
}

// ToReadableString 返回 IP 对应的地理位置信息（中文）
func ToReadableString(toString func(string, ...any) string, args ...any) string {
	if len(args) < 2 {
		return "未知"
	}
	return toString(args[0].(string), args[1:]...)
}

// ToReadableString 返回 IP 对应的地理位置信息（英文）
func ToReadableStringEn(toString func(string, ...any) string, args ...any) string {
	if len(args) < 2 {
		return "Unknown"
	}
	return toString(args[0].(string), args[1:]...)
}
