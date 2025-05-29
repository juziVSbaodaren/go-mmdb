# go-mmdb 创建于 2025-02-29
A Go library for reading MaxMind DB files.
## 自定义生成一个GeoLite2-City类型的数据库CustomGeoLite.mmdb，并使用go-mmdb解析

* 目的是将官方的GeoLite2-City.mmdb数据库中识别不了的ip段进行补充

* 步骤：
  1. 下载GeoLite2-City.mmdb数据库，并解压到指定目录；【https://github.com/P3TERX/GeoLite.mmdb/releases】
  2. 使用"github.com/maxmind/mmdbwriter"工具生成自定义的CustomGeoLite.mmdb数据库，具体操作如下：
   1. 新建一个导入ip信息的txt文件【geoIpInfo.txt】，文件内容格式如下：
      ```   
       Ip, Country.En, Country.Zh, Region.En, Region.Zh, City.En, City.Zh, Latitude, Longitude
       172.16.58.0/24,China,中国,Guangdong,广东,Shenzhen,深圳,23.1,113.3
      ```
   2. 运行命令：
      ```
       go run .\cmd\main.go -readGeoIpPath="./etc/geoIpInfo.txt" -writeGeoIpPath="./etc/CustomGeoLite.mmdb"

      ```
  3. 运行成功后，会在当前目录下生成CustomGeoLite.mmdb文件。
  4. 使用go-mmdb解析CustomGeoLite.mmdb数据库，具体操作如下：
   1. 运行命令：
      ```
        go run .\main.go  -ip="202.201.48.42" -mmdbPath="./etc/CustomGeoLite.mmdb" -customMmdbPath="./etc/CustomGeoLite.mmdb"

      ```
   2. 运行成功后，会输出如下信息：
      ```
      广东省-深圳市

      ```
      说明：ip地址202.201.48.42属于广东省-深圳市。