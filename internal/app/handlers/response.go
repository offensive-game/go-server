package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondOK(res *http.ResponseWriter, body interface{}) {
	(*res).Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	(*res).WriteHeader(200)
	_, err = (*res).Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func respondNotAuthorized(res *http.ResponseWriter) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(401)
}

func respondBadRequest(res *http.ResponseWriter, message string) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(400)

	_, err := (*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	if err != nil {
		panic(err)
	}
}

func respondServerError(res *http.ResponseWriter, message string) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(500)

	_, err := (*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	if err != nil {
		panic(err)
	}
}
