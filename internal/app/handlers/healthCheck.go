package handlers

import (
	"go-server/internal/app/middleware"
	"net/http"
)

type HealthCheck struct {
	appContext middleware.AppContext
}

func (hc *HealthCheck) SetAppContext (appContext middleware.AppContext) {
	hc.appContext = appContext
}

func (hc *HealthCheck) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	hc.appContext.Logger.Info("Handling Health Check")
	_, err := res.Write([]byte("Alive"))
	if err != nil {
		panic(err)
	}
}
