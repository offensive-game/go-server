package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token string `json:"token"`
}

type Login struct {
	appContext middleware.AppContext
	tx *sql.Tx
}

func (l *Login) SetAppContext(appContext middleware.AppContext) {
	l.appContext = appContext
}

func (l *Login) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body := LoginMessage{}
	l.tx = utils.GetTransactionFromContext(req)

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		utils.RespondBadRequest(&res, "Bad request")
		return
	}

	id, err := l.userIdForCredentials(body.Username, body.Password)
	if err != nil {
		utils.RespondBadRequest(&res, "Invalid username or password")
		return
	}

	session, err := l.createSession(id)

	if err != nil {
		utils.RespondBadRequest(&res, "Can't create session")
		return
	}

	res.Header().Set("Content-Type", "application/json")
	body.Token = session

	jsonData, err := json.Marshal(body)
	if err != nil {
		utils.RespondServerError(&res, "Server error")
		return
	}

	cookie := http.Cookie{
		Name: "offensive-login",
		Expires: time.Now().UTC().AddDate(0, 1, 0),
		Value: session,
		Path: "/",
	}

	http.SetCookie(res, &cookie)
	res.WriteHeader(200)
	_, err = res.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func (l Login) userIdForCredentials(username string, password string) (int64, error) {
	statement, err := l.tx.Prepare("SELECT id, password FROM users WHERE username = $1")
	if err != nil {
		return 0, err
	}

	var id int64
	var dbPassword string
	row := statement.QueryRow(username)
	err = row.Scan(&id, &dbPassword)
	if err != nil {
		return 0, err
	}

	success := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if success == nil {
		return id, nil
	}
	return 0, success
}

func (l Login) createSession(userId int64) (string, error) {
	// delete old sessions
	deleteStmt, err := l.tx.Prepare("DELETE FROM sessions WHERE userId = $1")
	if err != nil {
		return "", err
	}

	_, err = deleteStmt.Exec(userId)
	if err != nil {
		return "", err
	}

	// create new session
	createStmt, err := l.tx.Prepare("INSERT INTO sessions (userId, token) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}

	tokenBase := fmt.Sprintf("%d%d", userId, time.Now().UTC().Nanosecond())
	hash, err := bcrypt.GenerateFromPassword([]byte(tokenBase), bcrypt.MinCost)
	hashString := string(hash)

	_, err = createStmt.Exec(userId, hashString)
	if err != nil {
		return "", err
	}

	return hashString, nil
}
