package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Response struct {
	Status int    `json:"status"`
	Data   *Event `json:"data"`
}

func StartGeoServer(gs GeoServer) {
	http.HandleFunc("/getlocation", LocationHandler)
	Info("start geo-server at ", gs.Port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", gs.Port), nil)
}

func wrap(event *Event) string {
	var resp *Response
	if event != nil {
		resp = &Response{
			Status: 200,
			Data:   event,
		}
	} else {
		resp = &Response{
			Status: -1001,
			Data:   nil,
		}
	}
	if bytes, err := json.Marshal(resp); err == nil {
		return string(bytes)
	}
	return ""
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	lng, _ := strconv.ParseFloat(r.FormValue("longitude"), 64)
	lat, _ := strconv.ParseFloat(r.FormValue("latitude"), 64)
	event := Searcher.Get(lat, lng)
	io.WriteString(w, wrap(event))
	var (
		pro, city, county string
	)
	if event != nil {
		pro, city, county = event.Province, event.City, event.County
	}
	Infof("request (%f, %f) belong to %s%s%s", lng, lat, pro, city, county)
}
