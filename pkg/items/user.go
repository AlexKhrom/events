package items

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type User struct {
	ID           int64
	Login        string
	Email        string
	Fullname     string
	Password     string
	RefreshToken RefreshToken
}

type CodeUser struct {
	Token string `json:"token"`
	Code  string `json:"code"`
}

type CodeForm struct {
	Id       string
	Token    string
	Email    string
	Login    string
	Code     string
	Password string
	Expires  int64
	Times    string
}

type Session struct {
	Id        int64
	UserId    int64
	Expires   int64
	Token     string
	UserAgent string
}

type UserRepoInterface interface {
	FindUser(userLogin string) bool
	GetUserById(userId int) (User, bool)
	GetUserByEmail(userEmail string) (User, bool)
	NewUser(newUser User) (int64, error)
	NewCodeForm(form CodeForm) error
	NewSession(session Session) (int64, error)
	GetSession(sessionId int) (Session, bool)
	DeleteSession(sessionId int) error
	UpdateSession(sessionId int, refresh string) error
	GetCode(id string) (CodeForm, bool)
	CheckUserLogin(newUser User) bool
	CheckUserEmail(newUser User) bool
	GetUsersLen() int
	//SetRefreshToken(token string, exp int64, userLogin string) bool
}

var (
	ExampleTokenSecret = []byte("super secret string")
)

type UserRepo struct {
	Users map[string]User
	DB    *sql.DB
	Mut   sync.Mutex
}

type RefreshToken struct {
	Token string
	Exp   int64
}

func NewUserRepo(db *sql.DB) *UserRepo {
	repo := new(UserRepo)
	repo.DB = db
	//repo.Users = map[string]User{
	//	"user": {1, "User Userov", "", "d7316a3074d562269cf4302e4eed46369b523687", RefreshToken{
	//		"",
	//		0,
	//	}},
	//	"alex": {2, "Alexander Khromov", "", "731990ec145624822eee97d6bedb0a79efb28ccb", RefreshToken{
	//		"",
	//		0,
	//	}},
	//}
	return repo
}

func (repo *UserRepo) SetRefreshToken(token string, exp int64, email string) bool {

	if !repo.FindUser(email) {
		return true
	}

	expString := "{\"Token\":\"" + token + "\",\"Exp\":" + strconv.Itoa(int(exp)) + "}"

	_, err := repo.DB.Exec(
		"UPDATE `users` SET `refresh` = '" + expString + "' WHERE `users`.`email` = '" + Ecran(email) + "';",
	)

	if err != nil {
		fmt.Println("setfresh error = ", err)
		return true
	}
	return false
}

func (repo *UserRepo) FindUser(userEmail string) bool { // все кончено нужно делаит через id а не user-login

	rows, err := repo.DB.Query("SELECT email FROM users WHERE `users`.`email`='" + Ecran(userEmail) + "';")
	if err != nil {
		fmt.Println("find user error = ", err)
		return false
	}
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			fmt.Println("find user error = ", err)
			return false
		}

		if email == userEmail {
			return true
		}
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return false

}

func (repo *UserRepo) GetUserByEmail(userEmail string) (User, bool) {

	//fmt.Println("request sql = ", "SELECT * FROM users WHERE email ='"+strconv.Quote(userEmail)+"'")
	rows, err := repo.DB.Query("SELECT * FROM users WHERE email ='" + Ecran(userEmail) + "'")
	if err != nil {
		fmt.Println("find user error = ", err)
		return User{}, false
	}

	user := User{}
	refreshString := ""
	//refresh := RefreshToken{}

	for rows.Next() {

		err = rows.Scan(&user.ID, &user.Login, &user.Email, &user.Password, &refreshString)
		if err != nil {
			fmt.Println("find user error = ", err)
			return User{}, false
		}
		fmt.Println("user = ", user, "======================================")
		//
		//fmt.Println("resresh = ", refreshString)
		//
		//err = json.Unmarshal([]byte(refreshString), &refresh)
		//if err != nil {
		//	fmt.Println("get user error = ", err)
		//	return User{}, false
		//}
		//
		//user.RefreshToken = refresh

		return user, true
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()

	return User{}, false

}

func (repo *UserRepo) GetUserById(userId int) (User, bool) {

	stringUserId := strconv.Itoa(userId)

	rows, err := repo.DB.Query("SELECT * FROM users WHERE id =" + stringUserId)
	if err != nil {
		fmt.Println("find user error = ", err)
		return User{}, false
	}

	user := User{}
	refreshString := ""
	//refresh := RefreshToken{}

	for rows.Next() {

		err = rows.Scan(&user.ID, &user.Login, &user.Email, &user.Password, &refreshString)
		if err != nil {
			fmt.Println("find user error = ", err)
			return User{}, false
		}
		//
		//fmt.Println("resresh = ", refreshString)
		//
		//err = json.Unmarshal([]byte(refreshString), &refresh)
		//if err != nil {
		//	fmt.Println("get user error = ", err)
		//	return User{}, false
		//}
		//
		//user.RefreshToken = refresh

		return user, true
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()

	return User{}, false

}

func (repo *UserRepo) NewUser(newUser User) (int64, error) {
	repo.Mut.Lock()
	defer repo.Mut.Unlock()
	result, err := repo.DB.Exec(
		"INSERT INTO users (`email`,`login`,`password`,`refresh`) VALUES (?,?,?,?)",
		newUser.Email,
		newUser.Login,
		newUser.Password,
		"{}",
	)
	if err != nil {
		fmt.Println("some bad in newUser, err = ", err)
		return -1, err
	}
	id, err1 := result.LastInsertId()

	if err1 != nil {
		fmt.Println("some bad in newUser lastInsertId")
		return -1, err1
	}

	return id, err
}

func (repo *UserRepo) NewSession(session Session) (int64, error) {
	repo.Mut.Lock()
	defer repo.Mut.Unlock()
	result, err := repo.DB.Exec(
		"INSERT INTO sessions (`user_id`,`expires`,`refresh`,`user_agent`) VALUES (?,?,?,?)",
		session.UserId,
		session.Expires,
		session.Token,
		session.UserAgent,
	)
	if err != nil {
		fmt.Println("some bad in new session, err = ", err)
		return -1, err
	}
	id, err1 := result.LastInsertId()

	if err1 != nil {
		fmt.Println("some bad in newUser lastInsertId")
		return -1, err1
	}

	return id, err
}

func (repo *UserRepo) GetSession(sessionId int) (Session, bool) {
	repo.Mut.Lock()
	defer repo.Mut.Unlock()
	rows, err := repo.DB.Query("SELECT * FROM sessions WHERE id ='" + strconv.Itoa(sessionId) + "'")
	if err != nil {
		fmt.Println("find user error = ", err)
		return Session{}, false
	}

	session := Session{}

	for rows.Next() {

		err = rows.Scan(&session.Id, &session.UserId, &session.Expires, &session.Token, &session.UserAgent)
		if err != nil {
			fmt.Println("find user error = ", err)
			return Session{}, false
		}

		return session, true
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return Session{}, false
}

func (repo *UserRepo) DeleteSession(sessionId int) error {
	_, err := repo.DB.Exec(
		"DELETE FROM `sessions`  WHERE `sessions`.`id` = '" + strconv.Itoa(sessionId) + "';",
	)

	return err
}

func (repo *UserRepo) UpdateSession(sessionId int, refresh string) error {
	_, err := repo.DB.Exec(
		"UPDATE `sessions` SET `sessions`.`refresh` = '" + refresh + "' WHERE `sessions`.`id` = '" + strconv.Itoa(sessionId) + "';",
	)

	return err
}

func (repo *UserRepo) NewCodeForm(form CodeForm) error {
	repo.Mut.Lock()
	_, err := repo.DB.Exec(
		"INSERT INTO codes (`token`,`email`,`login`,`password`,`code`,`expires`,`times`) VALUES (?,?,?,?,?,?,?)",
		form.Token,
		form.Email,
		form.Login,
		form.Password,
		form.Code,
		form.Expires,
		0,
	)

	repo.Mut.Unlock()
	return err
}

func (repo *UserRepo) GetCode(token string) (CodeForm, bool) {
	repo.Mut.Lock()
	defer repo.Mut.Unlock()
	rows, err := repo.DB.Query("SELECT * FROM codes WHERE token ='" + Ecran(token) + "'")
	if err != nil {
		fmt.Println("find user error = ", err)
		return CodeForm{}, false
	}

	codesForm := CodeForm{}

	for rows.Next() {

		err = rows.Scan(&codesForm.Id, &codesForm.Token, &codesForm.Email, &codesForm.Login, &codesForm.Password, &codesForm.Code, &codesForm.Expires, &codesForm.Times)
		if err != nil {
			fmt.Println("find user error = ", err)
			return CodeForm{}, false
		}

		return codesForm, true
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return CodeForm{}, false

}

func (repo *UserRepo) CheckUserLogin(newUser User) bool {
	rows, err := repo.DB.Query("SELECT id FROM users WHERE `login`='" + Ecran(newUser.Login) + "'")
	if err != nil {
		fmt.Println("CheckUserLogin error = ", err)
		return true
	}
	for rows.Next() {
		return true
	}
	return false
}

func (repo *UserRepo) CheckUserEmail(newUser User) bool {
	rows, err := repo.DB.Query("SELECT id FROM users WHERE `email`='" + Ecran(newUser.Email) + "'")
	if err != nil {
		fmt.Println("CheckUserLogin error = ", err)
		return true
	}
	for rows.Next() {
		return true
	}
	return false
}

func (repo *UserRepo) GetUsersLen() int {
	rows, err := repo.DB.Query("SELECT id FROM users")
	if err != nil {
		fmt.Println("find user error = ", err)
		return 0
	}
	c := 0
	for rows.Next() {
		c++
	}
	return c
}

func (u User) NewUserToken(email string, sessionId int64) *jwt.Token {

	timeNow := time.Now().Unix()
	exprTime := time.Now().Add(time.Minute * 15).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"id":        u.ID,
			"sessionId": sessionId,
			"email":     email,
		},
		"iat": timeNow,
		"exp": exprTime,
	})

	return token
}

func (u User) MakeUserTokenResp(w http.ResponseWriter, token *jwt.Token) ([]byte, error) {
	tokenString, err := token.SignedString(ExampleTokenSecret)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	resp, err := json.Marshal(map[string]interface{}{
		"status": http.StatusOK,
		"token":  tokenString,
	})

	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant pack payload")
		return nil, nil
	}

	return resp, nil
}

func (u User) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func GetTokenPayload(w http.ResponseWriter, r *http.Request) jwt.MapClaims {

	inTokenCookie, err1 := r.Cookie("token")

	if err1 != nil || inTokenCookie == nil {
		fmt.Println("bad token in cookie gettokenpayload")
		return nil
	}

	inToken := inTokenCookie.Value

	fmt.Println("token = ", inToken)

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return ExampleTokenSecret, nil
	}

	token, _ := jwt.Parse(inToken, hashSecretGetter)
	//if err != nil {
	//	fmt.Println("bad token")
	//	jsonError(w, http.StatusUnauthorized, "bad token")
	//	return nil
	//}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		jsonError(w, http.StatusUnauthorized, "no payload")
	}
	return payload
}

func jsonError(w io.Writer, status int, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	if err != nil {
		fmt.Println("bad in jsonError")
		return
	}
	_, err2 := w.Write(resp)
	if err2 != nil {
		jsonError(w, http.StatusBadRequest, "bad write response")
	}
}

func Ecran(str string) string {
	return strconv.Quote(str)[1 : len(strconv.Quote(str))-1]
}
