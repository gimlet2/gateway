package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gimlet2/gateway/config"
	"github.com/gimlet2/gateway/utils"
)

func main() {

	log.Print("Start Server")
	conf := config.Load(os.Getenv("config"))
	get("/config", func(w http.ResponseWriter, r *http.Request) {
		j, _ := json.Marshal(conf)
		w.Write(j)
	})
	for _, endpoint := range conf.Endpoints {
		http.HandleFunc(conf.Path+endpoint.GetPath(), func(w http.ResponseWriter, r *http.Request) {
			if len(endpoint.Methods) != 0 && !utils.Contains(endpoint.Methods, r.Method) {
				utils.WriteErrorWithCode(w, http.StatusMethodNotAllowed, "Method is not allowed")
				return
			}
			for _, route := range endpoint.Routes {
				if route.Matching(r) {
					route.Drop().Forward(w, r)
					break
				}
			}

		})
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func post(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handler(w, r)
		} else {
			utils.WriteErrorWithCode(w, http.StatusMethodNotAllowed, "Method is not allowed")
		}
	})
}

func get(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handler(w, r)
		} else {
			utils.WriteErrorWithCode(w, http.StatusMethodNotAllowed, "Method is not allowed")
		}
	})
}
