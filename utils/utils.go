package utils

import (
	"fmt"
	"net/http"
)

func Contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func ContainsAll(slice []string, vals []string) bool {
	for _, v := range slice {
		res := true
		for _, val := range vals {
			if v != val {
				res = false
				break
			}
			if res {
				return true
			}
		}
	}
	return false
}

func ContainsOne(slice []string, vals []string) bool {
	for _, v := range slice {
		for _, val := range vals {
			if v == val {
				return true
			}
		}
	}
	return false
}

func WriteError(w http.ResponseWriter, message string) {
	w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", message)))
}

func WriteErrorWithCode(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", message)))
}
