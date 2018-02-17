package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gimlet2/gateway/utils"
)

type API struct {
	Path      string     `json:"path"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Sticky  bool     `json:"sticky"`
	Routes  []Route  `json:"routes"`
	Methods []string `json:"methods"`
}

func (e Endpoint) GetPath() string {
	if e.Path == "" {
		return e.Name
	}
	return e.Path
}

type Route struct {
	Name     string   `json:"name"`
	Upstream Upstream `json:"upstream"`
	Match    Match    `json:"match"`
}

func (route Route) Matching(r *http.Request) bool {
	return route.Match.Matching(r)
}

type Upstream struct {
	Uri string `json:"uri"`
}

func (u Upstream) Forward(w http.ResponseWriter, r *http.Request) {
	client := http.DefaultClient
	var newRequest, _ = http.NewRequest(r.Method, u.Uri, r.Body)
	resp, e := client.Do(newRequest)
	if e != nil {
		log.Printf("Error %e", e)
	}
	w.WriteHeader(resp.StatusCode)

	for n, h := range resp.Header {
		for _, s := range h {
			w.Header().Add(n, s)
		}
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	w.Write(buff)

}

type Match struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Query   map[string][]string `json:"query,omitempty"`
	Json    map[string][]string `json:"json,omitempty"`
}

func (m Match) Matching(r *http.Request) bool {
	for n, v := range m.Headers {
		if !utils.Contains(v, r.Header.Get(n)) {
			return false
		}
	}
	for n, v := range m.Query {
		if !utils.ContainsOne(v, map[string][]string(r.URL.Query())[n]) {
			return false
		}
	}
	// check json here
	// for n, v := range m.Headers {
	// 	if !utils.Contains(v, r.Header.Get(n)) {
	// 		return false
	// 	}
	// }
	return true
}

func Load(configFile string) API {
	var api API
	content, _ := ioutil.ReadFile(configFile)
	json.Unmarshal(content, &api)
	return api
}
