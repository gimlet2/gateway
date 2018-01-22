package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"github.com/gimlet2/gateway/config"
)

func main() {

	log.Print("Start Server")
	conf := config.Load(os.Getenv("config"))
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		j, _ := json.Marshal(conf)
		w.Write(j)
	})
	http.HandleFunc("/decision", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			//decoder := json.NewDecoder(r.Body)
			//var request DecisionRequest
			//decoder.Decode(&request)
			//log.Print(string(request.Data.Type))
			//j, _ := json.Marshal(request)
			//w.Write(j)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Error"))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

