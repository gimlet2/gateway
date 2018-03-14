package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
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
	Name     string     `json:"name"`
	Upstream []Upstream `json:"upstream"`
	Match    Match      `json:"match"`
}

type Upstream struct {
	Uri    string  `json:"uri"`
	Weight float32 `json:"weight"`
}

type Match struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Query   map[string][]string `json:"query,omitempty"`
	Json    map[string][]string `json:"json,omitempty"`
}

func (route Route) Matching(r *http.Request) bool {
	return route.Match.Matching(r)
}

func (route Route) Drop() Upstream {
	total := float32(0.0)
	for _, u := range route.Upstream {
		total += u.Weight
	}
	index := 0
	r := rand.Float32()
	for total >= r {
		total -= r
		r = rand.Float32()
		index++
	}
	return route.Upstream[index]
}

func (u Upstream) Forward(w http.ResponseWriter, r *http.Request) {
	client := http.DefaultClient
	var newRequest, _ = http.NewRequest(r.Method, u.Uri, r.Body)
	resp, e := client.Do(newRequest)
	if e != nil {
		log.Printf("Failed to make upstream request %e", e)
		utils.WriteErrorWithCode(w, 500, "Failed to make upstream request")
	}
	w.WriteHeader(resp.StatusCode)

	for n, h := range resp.Header {
		for _, s := range h {
			w.Header().Add(n, s)
		}
	}
	buff, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Printf("Failed to make upstream request %e", e)
		utils.WriteErrorWithCode(w, 500, "Failed to make upstream request")
	}
	w.Write(buff)

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
