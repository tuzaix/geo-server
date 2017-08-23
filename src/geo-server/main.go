package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	//"runtime"
)

var (
	conf = flag.String("conf", "../conf/geo.conf", "geo toml config")
)

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// 加载配置
	LoadConfig(*conf)

	// 启动pprof
	go func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", Config.GeoServer.PprofPort), nil)
	}()

	// 初始化日志
	InitLogger(Config.LogConfig)

	// 初始化lnglat
	BuildLngLat(Config.GeoServer.LocFile)
	//
	// 启动 服务监听
	StartGeoServer(Config.GeoServer)
}
