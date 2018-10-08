package handlers

import (
	"github.com/sirupsen/logrus"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"net/http"
	"time"
)

type SignoutMessage struct {
	Username string `json:"username"`
}

type Signout struct {
	appContext middleware.AppContext
	log        *logrus.Entry
}

func (s *Signout) SetAppContext(appContext middleware.AppContext) {
	s.appContext = appContext
}

func (s *Signout) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.appContext.Logger.Info("Logging out")
	tx := utils.GetTransactionFromContext(req)
	user := utils.GetUserFromContext(req)

	s.log = s.appContext.Logger.WithFields(logrus.Fields{
		"userId": user.Id,
	})

	statement, err := tx.Prepare(`
	  DELETE FROM sessions WHERE userId = $1
	`)
	if err != nil {
		s.log.Error("Unable to prepare statement to delete session")
		panic(err)
	}

	_, err = statement.Exec(user.Id)
	if err != nil {
		s.log.Error("Unable to delete session")
		panic(err)
	}

	deleteCookie := &http.Cookie{
		Name:     "offensive-login",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(res, deleteCookie)
	utils.RespondOK(&res, nil)
}
