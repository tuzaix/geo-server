package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

var (
	Config GeoConfig
)

type GeoConfig struct {
	GeoServer GeoServer `toml:"http-server"`
	LogConfig LogConfig `toml:"log"`
}

type GeoServer struct {
	Port      int    `toml:"port"`
	PprofPort int    `toml:"pprof_port"`
	LocFile   string `toml:"loction_file"`
}

type LogConfig struct {
	LogLevel     string `toml:"level"`
	LogConsole   int    `toml:"console"`
	LogDir       string `toml:"dir"`
	LogFilename  string `toml:"filename"`
	LogCount     int    `toml:"count"`
	LogSuffix    string `toml:"suffix"`
	LogColorfull int    `toml:"colorfull"`
}

// 初始化全局配置文件
func LoadConfig(filename string) {
	var (
		data []byte
		err  error
	)

	data, err = ioutil.ReadFile(filename)

	if err != nil {
		panic("read configuration file failed " + err.Error())
	}

	if _, err = toml.Decode(string(data), &Config); err != nil {
		panic("toml decode failed " + err.Error())
	}
	return
}
