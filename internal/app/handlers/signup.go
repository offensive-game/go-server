package handlers

import (
	"database/sql"
	"encoding/json"
	"go-server/internal/app/middleware"
	"go-server/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignupMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Signup struct {
	appContext middleware.AppContext
	tx *sql.Tx
}

func (s *Signup) SetAppContext(appContext middleware.AppContext) {
	s.appContext = appContext
}

func (s *Signup) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.tx = utils.GetTransactionFromContext(req)
	body := SignupMessage{}

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		panic(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	body.Password = string(hash)

	exists, err := s.userWithCredentialsExist(body.Username, body.Password)
	if err != nil {
		panic(err)
	}

	if exists {
		utils.RespondBadRequest(&res, "User already exists")
	} else {
		ok := s.createNewUser(body)
		if ok {
			body.Password = ""
			utils.RespondOK(&res, body)
		} else {
			utils.RespondServerError(&res, "Unable to create a new user")
		}
	}
}

func (s Signup) userWithCredentialsExist(username string, email string) (bool, error) {
	queryStatement, err := s.tx.Prepare("SELECT COUNT(*) from users WHERE username = $1 OR email = $2")
	if err != nil {
		return false, err
	}

	var count int
	row := queryStatement.QueryRow(username, email)
	err = row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (s Signup) createNewUser(signupMessage SignupMessage) bool {
	query, err := s.tx.Prepare("INSERT INTO users (username, password, email) VALUES($1, $2, $3)")
	if err != nil {
		return false
	}

	_, err = query.Exec(signupMessage.Username, signupMessage.Password, signupMessage.Email)
	if err != nil {
		return false
	}

	return true
}
