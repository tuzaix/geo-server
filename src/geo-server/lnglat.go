package main

import (
	"fmt"
	"github.com/gansidui/geohash"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var (
	Searcher Search
	Rad      = math.Pi / 180.0
)

const (
	Precision = 12
	Radius    = 6371000 //6378137
)

type Event struct {
	Province string  `json:"province"`
	City     string  `json:"city"`
	County   string  `json:"county"`
	Lng      float64 `json:"longitude"`
	Lat      float64 `json:"latitude"`
	GeoHash  string  `json:"geohash"`
}

type Map struct {
	Map map[string][]Event
}

type Search struct {
	Search []Map
}

func (this Search) Set(precision int, mapping Map) {
	this.Search[precision] = mapping
}

func (this Search) Get(lat, lng float64) (event *Event) {
	hash, _ := geohash.Encode(lat, lng, Precision)
	sizeCache := make(map[string]float64)
	var (
		size, new_size float64
		prefix, key    string
		rets           []Event
		ret            Event
		ok             bool
	)

	for n := Precision - 1; n > 1; n-- {
		prefix = string(hash[:n+1])
		if rets, ok = this.Search[n].Map[prefix]; ok {
			event = &(rets[0])
			if len(rets) == 1 {
				return
			} else {
				key = fmt.Sprintf("%f_%f_%f_%f", lat, lng, event.Lat, event.Lng)
				if size, ok = sizeCache[key]; !ok {
					size = EarthDistance(lat, lng, event.Lat, event.Lng)
					sizeCache[key] = size
				}
				for _, ret = range rets {
					key = fmt.Sprintf("%f_%f_%f_%f", lat, lng, ret.Lat, ret.Lng)
					if new_size, ok = sizeCache[key]; !ok {
						new_size = EarthDistance(lat, lng, ret.Lat, ret.Lng)
						sizeCache[key] = new_size
					}
					if size > new_size {
						size = new_size
						event = &ret
					}
				}
			}
			return
		}
	}
	return
}

func EarthDistance(src_lat, src_lng, dist_lat, dist_lng float64) float64 {
	src_lat = src_lat * Rad
	src_lng = src_lng * Rad
	dist_lat = dist_lat * Rad
	dist_lng = dist_lng * Rad
	theta := dist_lng - src_lng
	dist := math.Acos(math.Sin(src_lat)*math.Sin(dist_lat) + math.Cos(src_lat)*math.Cos(src_lat)*math.Cos(theta))
	return dist * Radius
}

// function 定义
func BuildLngLat(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		Error("load geo meta data error ", err.Error())
	}
	sdata := string(data)
	parts := strings.Split(sdata, "\n")

	var (
		events []Event = make([]Event, len(parts)-1)
		lng    float64
		lat    float64
		hash   string
	)
	for idx, part := range parts {
		subParts := strings.Split(part, "|")
		if len(subParts) != 7 {
			Warning("geo meta format invalid")
			continue
		}
		lng, err = strconv.ParseFloat(subParts[5], 64)
		if err != nil {
			Error("conv longtitude str to float error", err.Error())
			continue
		}
		lat, err = strconv.ParseFloat(subParts[6], 64)
		if err != nil {
			Error("conv latitude str to float error", err.Error())
			continue
		}
		hash, _ = geohash.Encode(lat, lng, Precision)
		events[idx] = Event{
			Province: subParts[0],
			City:     subParts[1],
			County:   subParts[2],
			Lng:      lng,
			Lat:      lat,
			GeoHash:  hash,
		}
	}
	Info("load geo meta success")
	Searcher = Search{Search: make([]Map, Precision)}
	for n := 0; n < Precision; n++ {
		mapping := make(map[string][]Event)
		subfix := n + 1
		for _, e := range events {
			prefix := e.GeoHash[:subfix]
			if _, ok := mapping[prefix]; !ok {
				mapping[prefix] = make([]Event, 0)
			}
			mapping[prefix] = append(mapping[prefix], e)
		}
		Searcher.Set(n, Map{Map: mapping})
	}
	Info("build geo index success")
}
