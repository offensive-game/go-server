package handlers

import (
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
)

type Me struct {
	appContext middleware.AppContext
	log *logrus.Entry
}

func (me *Me) SetAppContext(appContext middleware.AppContext) {
	me.appContext = appContext
}

func (me *Me) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	user := utils.GetUserFromContext(req)
	me.log = me.appContext.Logger.WithFields(logrus.Fields{"userId": user.Id})

	me.log.Info("Handling /Me request")
	utils.RespondOK(&res, user)
}
