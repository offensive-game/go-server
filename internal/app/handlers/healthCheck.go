package handlers

import (
	"fmt"
	"go-server/internal/app/middleware"
	"go-server/internal/app/models"
	"net/http"
)

type HealthCheck struct {
	appContext middleware.AppContext
}

func (hc HealthCheck) SetAppContext (appContext middleware.AppContext) {
	hc.appContext = appContext
}

func (hc HealthCheck) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(models.User)
	_, err := res.Write([]byte(fmt.Sprintf("Alive %s", user.Username)))
	if err != nil {
		panic(err)
	}
}
