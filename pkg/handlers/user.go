package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"events/pkg/items"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"
)

type UserHand struct {
	Repo items.UserRepoInterface
}

type LoginForm struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserHand(db *sql.DB) *UserHand {
	us := new(UserHand)
	us.Repo = items.NewUserRepo(db)
	return us
}

//func (h *UserHand) GetLogin(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("hi get login")
//
//	http.ServeFile(w, r, filepath.Join(".././static/dist", "index.html"))
//}

func (h *UserHand) Login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hi login")

	body, err3 := ioutil.ReadAll(r.Body)
	if err3 != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fd := &LoginForm{}

	err := json.Unmarshal(body, fd)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	fmt.Println("fd = ", fd)

	user, err1 := h.Repo.GetUserByEmail(fd.Email)

	fmt.Println("user = ", user)
	fmt.Println("user errr = ", err1)

	if !err1 {
		JSONError(w, http.StatusBadRequest, "bad login")
		return
	}

	shaPassword := h.makeSha1(fd.Email + fd.Password + string(items.ExampleTokenSecret))

	fmt.Println("string = ", fd.Email+fd.Password+string(items.ExampleTokenSecret))
	fmt.Println("shaPassword = ", shaPassword)

	if !err1 || shaPassword != user.Password {
		JSONError(w, http.StatusUnauthorized, "bad login or password")
		fmt.Println("bad login or password")
		return
	}

	//удаляю сессию которая была до у этого пользователя
	tokenPaylod := items.GetTokenPayload(w, r)

	if tokenPaylod != nil {
		if tokenPaylod == nil {
			JSONError(w, http.StatusBadRequest, "bad access token")
			return
		}

		lastSessionId, ok := tokenPaylod["user"].(map[string]interface{})["sessionId"]
		if !ok {
			JSONError(w, http.StatusBadRequest, "bad access token without sessionId")
			return
		}

		if lastSessionId != nil {
			err := h.Repo.DeleteSession(int(lastSessionId.(float64)))
			if err != nil {
				JSONError(w, http.StatusBadRequest, "bad with delete session")
				return
			}
		}
	}

	refreshToken, err3 := user.NewRefreshToken()

	newSession := items.Session{
		UserId:    user.ID,
		Expires:   time.Now().Add(time.Hour * 24 * 30).Unix(),
		Token:     refreshToken,
		UserAgent: r.Header.Get("User-Agent"),
	}

	sessionId, err3 := h.Repo.NewSession(newSession)

	if err3 != nil {
		fmt.Println("some bad new session - ", err3)
		return
	}

	token := user.NewUserToken(fd.Email, sessionId)

	resp, err := user.MakeUserTokenResp(w, token)
	if err != nil {
		return
	}

	fmt.Println(token)
	fmt.Println(string(resp))

	tokenString, err := token.SignedString(items.ExampleTokenSecret)

	expiration := time.Now().Add(time.Hour * 24 * 30)

	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expiration,
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	fmt.Println("refreshToken = ", refreshToken)

	cookie = http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	//fmt.Println("cookies = ",r.Cookies())

	//fmt.Println("redirect")
	//http.Redirect(w, r, "/login", 303)
	//fmt.Println("redirect2")

	_, err2 := w.Write(resp)
	if err2 != nil {
		JSONError(w, http.StatusBadRequest, "bad write response")
	}

	//fmt.Println("url = ",  r.Header.Get("Referer"))
	//
	////http.Redirect(w,r,"http://localhost:8085/static/pages/manager.html", 302)
	//
	//w.Header().Set("Location","http://localhost:8085/static/pages/manager.html")
	//w.Header().Add("Set-Cookie", " path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT")
	//w.WriteHeader(303)
	//
	//fmt.Println("\nend login\n")

}

func (h *UserHand) RefreshTokens(w http.ResponseWriter, r *http.Request) {

	tokenPaylod := items.GetTokenPayload(w, r)

	if tokenPaylod == nil {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	userId, ok := tokenPaylod["user"].(map[string]interface{})["id"]
	if !ok {
		http.Redirect(w, r, "/login", 302)

		JSONError(w, http.StatusBadRequest, "bad access token")
		return
	}

	sessionId, ok := tokenPaylod["user"].(map[string]interface{})["sessionId"]
	if !ok {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad access token without sessionId")
		return
	}

	fmt.Println("\n\nuserID = ", userId)

	oldRefreshToken, err := r.Cookie("refreshToken")

	if err != nil {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad refresg in cookie")
		return
	}
	fmt.Println("oldRefresh = ", oldRefreshToken.Value)

	user, err1 := h.Repo.GetUserById(int(userId.(float64)))
	fmt.Println("user = ", user)
	fmt.Println("user errr = ", err1)
	if !err1 {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad user in base from jwt")
		return
	}

	session, err1 := h.Repo.GetSession(int(sessionId.(float64)))

	fmt.Println("session = ", session)

	if session.Token != oldRefreshToken.Value {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad refresh tokens are not equal")
		return
	}

	if session.Expires <= time.Now().Unix() {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "bad refreshtoken is expire")
		return
	}

	newTokenString, err := user.NewRefreshToken()
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		JSONError(w, http.StatusBadRequest, "can't creat new refresh")
		return
	}
	//newRefreshToken := items.RefreshToken{
	//	Token: newTokenString,
	//	Exp: time.Now().Unix(),
	//}

	err = h.Repo.UpdateSession(int(sessionId.(float64)), newTokenString)

	if err != nil {
		fmt.Println("err with update sessions sql, err = ", err)
		JSONError(w, http.StatusBadRequest, "err with update sessions sq")
		return
	}

	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    newTokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	token := user.NewUserToken(user.Email, int64(sessionId.(float64)))

	//resp, err := user.MakeUserTokenResp(w, token)
	//if err != nil {
	//	return
	//}

	tokenString, err := token.SignedString(items.ExampleTokenSecret)

	expiration := time.Now().Add(time.Hour * 24 * 30)

	cookie = http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expiration,
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	//_, err2 := w.Write(resp)
	//if err2 != nil {
	//	JSONError(w, http.StatusBadRequest, "bad write response")
	//}

}

func (h *UserHand) SingUp(w http.ResponseWriter, r *http.Request) {

	body, err3 := ioutil.ReadAll(r.Body)
	if err3 != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
	}
	fd := &LoginForm{}

	fmt.Println(string(body))

	err := json.Unmarshal(body, fd)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	fmt.Println("fd = ", fd)

	//isUser := h.Repo.FindUser(fd.Login)
	//if isUser {
	//	JSONError(w, http.StatusBadRequest, "this username is already exist")
	//	return
	//}

	newUser := items.User{
		//ID:       h.Repo.GetUsersLen() + 1,
		Login:    fd.Login,
		Email:    fd.Email,
		Password: h.makeSha1(fd.Email + fd.Password + string(items.ExampleTokenSecret)),
	}

	if h.Repo.CheckUserEmail(newUser) {
		JSONError(w, 425, "this Email is already exist")
		return
	}

	if h.Repo.CheckUserLogin(newUser) {
		JSONError(w, 426, "this Login is already exist")
		return
	}

	code := h.generCode()
	fmt.Println("code = ", code)

	h.sendEmail(newUser.Email, code)

	idCodeForm, err := newUser.NewRefreshToken()

	if err != nil {
		fmt.Println("err = ", err.Error())
		JSONError(w, http.StatusBadRequest, "some bad in generation id for code form")
		return
	}

	newCodeForm := items.CodeForm{
		Token:    idCodeForm,
		Code:     code,
		Email:    newUser.Email,
		Login:    newUser.Login,
		Password: newUser.Password,
		Expires:  time.Now().Add(time.Minute * 2).Unix(),
	}

	err = h.Repo.NewCodeForm(newCodeForm)

	fmt.Println("id = ", idCodeForm)
	if err != nil {
		fmt.Println("err = ", err.Error())
		JSONError(w, http.StatusBadRequest, "some bad in new code form insert")
		return
	}

	cookie := http.Cookie{
		Name:     "codeId",
		Value:    idCodeForm,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	fmt.Println("ok!!!")

}

func (h *UserHand) CheckCode(w http.ResponseWriter, r *http.Request) {

	body, err3 := ioutil.ReadAll(r.Body)
	if err3 != nil {
		JSONError(w, http.StatusBadRequest, "can't read request body")
		return
	}
	fd := &items.CodeUser{}

	fmt.Println(string(body))

	err := json.Unmarshal(body, fd)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	fmt.Println("fd = ", fd)

	codesForm, err1 := h.Repo.GetCode(fd.Token)

	fmt.Println("codesForm = ", codesForm)

	if !err1 {
		fmt.Println("err = ", err1)
		JSONError(w, http.StatusBadRequest, "some bad in codesForm get from database")
		return
	}

	if codesForm.Code != fd.Code {
		fmt.Println("some bad in check code - codes aren't equal")
		JSONError(w, http.StatusBadRequest, "some bad in check code - codes aren't equal")
		return
	}
	fmt.Println("hi2")
	newUser := items.User{
		Login:    codesForm.Login,
		Email:    codesForm.Email,
		Password: codesForm.Password,
	}

	id, err := h.Repo.NewUser(newUser)
	newUser.ID = id
	fmt.Println("hi3")

	if err != nil {
		fmt.Println("error = ", err.Error())
		JSONError(w, 500, err.Error())
		return
	}

	fmt.Println("hi4")
	refreshToken, err3 := newUser.NewRefreshToken()

	newSession := items.Session{
		UserId:    newUser.ID,
		Expires:   time.Now().Add(time.Hour * 24 * 30).Unix(),
		Token:     refreshToken,
		UserAgent: r.Header.Get("User-Agent"),
	}

	sessionId, err3 := h.Repo.NewSession(newSession)

	token := newUser.NewUserToken(newUser.Email, sessionId)

	resp, err := newUser.MakeUserTokenResp(w, token)

	fmt.Println("hi5")

	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	fmt.Println(token)
	fmt.Println(string(resp))

	tokenString, err := token.SignedString(items.ExampleTokenSecret)

	expiration := time.Now().Add(time.Hour * 24 * 30)

	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expiration,
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	fmt.Println("refreshToken = ", refreshToken)

	cookie = http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: false, // httponly = true + secure
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	_, err2 := w.Write(resp)
	if err2 != nil {
		JSONError(w, http.StatusBadRequest, "bad write response")
	}

}

func (user *UserHand) sendEmail(toEmail, code string) {
	fmt.Println("say hi")
	from := "sasha@khromov.su"
	password := "19ui3jfq3in19"

	toEmailAddress := toEmail
	to := []string{toEmailAddress}

	host := "smtp.yandex.ru"
	port := "25"
	address := host + ":" + port

	subject := "<h2>подтверждение почты</h2>\n"
	body := "ваш код : " + code
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)

	if err != nil {
		fmt.Println("some bad in send email = ", err)
		panic(err)
	}

	fmt.Println("email is enter")

}

func GetUserId(w http.ResponseWriter, r *http.Request) (int, bool) { // (userId, Ok)
	tokenPaylod := items.GetTokenPayload(w, r)

	if tokenPaylod == nil {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return -1, false
	}

	userId, ok := tokenPaylod["user"].(map[string]interface{})["id"]
	if !ok {
		JSONError(w, http.StatusBadRequest, "bad access token")
		return -1, false
	}
	return int(userId.(float64)), true
}

func (user *UserHand) generCode() string {
	rand.Seed(time.Now().UnixNano())
	var digits = "0123456789"
	var str = ""
	for i := 0; i < 4; i++ {
		str += string(digits[rand.Intn(10)])
	}
	return str
}

func (h *UserHand) makeSha1(s string) string {
	ha := sha256.New()
	ha.Write([]byte(s))
	sha1_hash := hex.EncodeToString(ha.Sum(nil))

	fmt.Println(s, sha1_hash)
	return sha1_hash
}
