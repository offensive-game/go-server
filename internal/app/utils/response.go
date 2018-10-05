package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondOK(res *http.ResponseWriter, body interface{}) {
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

func RespondNotAuthorized(res *http.ResponseWriter) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(401)
}

func RespondBadRequest(res *http.ResponseWriter, message string) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(400)

	_, err := (*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	if err != nil {
		panic(err)
	}
}

func RespondServerError(res *http.ResponseWriter, message string) {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(500)

	_, err := (*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	if err != nil {
		panic(err)
	}
}
