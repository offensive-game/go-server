package handlers
import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type SignupMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
}

type Signup struct {
	Db *sql.DB
}

func (s Signup) ServeHTTP (res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body := SignupMessage{}
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			panic(err)
		}

		exists, err := userWithCredentialsExist(s, body.Username, body.Password)
		if err != nil {
			panic(err)
		}

		if exists {
			respondBadRequest(&res, "User already exists")
		} else {
			createNewUser(s, body)
			respondOK(&res, body)
		}
	}
}

func userWithCredentialsExist (s Signup, username string, password string) (bool, error) {
	queryStatement, err := s.Db.Prepare("SELECT COUNT(*) from users WHERE username = $1 OR email = $2")
	if err != nil {
		return false, err
	}

	var count int
	queryStatement.QueryRow(username, password).Scan(&count)

	return count != 0, nil
}

func createNewUser(s Signup, signupMessage SignupMessage) bool {
	query, err := s.Db.Prepare("INSERT INTO users (username, password, email) VALUES($1, $2, $3)")
	if err != nil {
		return false
	}

	_, err = query.Exec(signupMessage.Username, signupMessage.Password, signupMessage.Email)
	if err != nil {
		return false
	}

	return true
}