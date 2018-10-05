package handlers

import (
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
)

type Me struct {
	appContext middleware.AppContext
}

func (me Me) SetAppContext(appContext middleware.AppContext) {
	me.appContext = appContext
}

func (me Me) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	user := utils.GetuserFromContext(req)
	utils.RespondOK(&res, user)
}
