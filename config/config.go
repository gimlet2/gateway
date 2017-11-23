package config

import (
	"encoding/json"
	"io/ioutil"
)

type API struct {
	Path      string     `json:"path"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name   string  `json:"name"`
	Path   string  `json:"path"`
	Sticky bool    `json:"sticky"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Name     string   `json:"name"`
	Upstream Upstream `json:"upstream"`
	Match    Match    `json:match`
}

type Upstream struct {
	Uri string `json:"uri"`
}

type Match struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Query   map[string][]string `json:"query,omitempty"`
	Json    map[string][]string `json:"json,omitempty"`
}

func Load(configFile string) API {
	var api API
	content, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(content, &api)
	return api
}
