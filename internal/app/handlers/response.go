package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondOK(res *http.ResponseWriter, body interface{}) error {
	(*res).Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	(*res).WriteHeader(200)
	(*res).Write(jsonData)

	return nil
}

func respondNotAuthorized(res *http.ResponseWriter) error {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(401)

	return nil
}

func respondBadRequest(res *http.ResponseWriter, message string) error {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(400)

	(*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	return nil
}

func respondServerError(res *http.ResponseWriter, message string) error {
	(*res).Header().Set("Content-Type", "application/json")
	(*res).WriteHeader(500)

	(*res).Write([]byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	return nil
}
