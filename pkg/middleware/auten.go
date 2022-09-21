package middleware

import (
	"database/sql"
	"events/pkg/handlers"
	"events/pkg/items"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var (
	answerCode = 302
	noAuthUrls = map[string]string{
		"/api/posts": "POST",
		"/api/post/": "DELETE_POST",
		"/upvote":    "GET",
		"/downvote":  "GET",
		"/unvote":    "GET",
	}
	// noSessUrls = map[string]struct{}{
	//	"/": struct{}{},
	// }
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//next.ServeHTTP(w, r)

		fmt.Println("\n\n\n=====================")
		fmt.Println("time: ", time.Now().String()[:27])
		fmt.Println("r.Method: ", r.Method)
		//fmt.Println("vers method = ", r.Proto)
		fmt.Println("url: ", r.URL.Path)
		fmt.Println("\n=====================")
		
		rMethod, ok := noAuthUrls[r.URL.Path]

		if (strings.Contains(r.URL.Path, "/login") && r.Method == "GET") ||
			strings.Contains(r.URL.Path, "/signUp") && r.Method == "GET" ||
			strings.Contains(r.URL.Path, "/css/") && r.Method == "GET" ||
			strings.Contains(r.URL.Path, "/js/") && r.Method == "GET" ||
			strings.Contains(r.URL.Path, "/favicon.ico") && r.Method == "GET" ||
			strings.Contains(r.URL.Path, "/api/login") && r.Method == "POST" ||
			strings.Contains(r.URL.Path, "/api/signUp") && r.Method == "POST" ||
			strings.Contains(r.URL.Path, "/api/checkCode") && r.Method == "POST" {

			next.ServeHTTP(w, r)
			return

		} else if !ok || r.Method != rMethod {

			tokenCookie, err := r.Cookie("token")
			if err != nil {
				fmt.Println("\nno cookie token")
				http.Redirect(w, r, "/login", answerCode)
				return
			}

			//здесь проверять строку из токена

			if tokenCookie == nil {
				fmt.Println("\nno auth1")
				http.Redirect(w, r, "/login", answerCode)
				return
			}

			fmt.Println("cookie = ", tokenCookie.Value)

			inToken := tokenCookie.Value

			if inToken == "" {
				fmt.Println("\nno auth2")
				http.Redirect(w, r, "/login", answerCode)
				return
			}

			fmt.Println("token = ", inToken)

			hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
				method, ok1 := token.Method.(*jwt.SigningMethodHMAC)
				if !ok1 || method.Alg() != "HS256" {
					return nil, fmt.Errorf("bad sign method")
				}
				return items.ExampleTokenSecret, nil
			}

			token, err := jwt.Parse(inToken, hashSecretGetter)

			if err != nil || !token.Valid {
				fmt.Println("jwt is expire", err, "/////", token.Valid)

				dsn := "test:root@tcp(localhost:3306)/events?"
				// указываем кодировку
				dsn += "&charset=utf8mb4"
				// отказываемся от prapared statements
				// параметры подставляются сразу
				dsn += "&interpolateParams=true"

				db, err := sql.Open("mysql", dsn)
				if err != nil {
					fmt.Println("errpr!!!!", err)
					return
				}

				db.SetMaxOpenConns(10)

				err = db.Ping() // вот тут будет первое подключение к базе
				if err != nil {
					panic(err)
				}

				userRepo := handlers.NewUserHand(db)

				userRepo.RefreshTokens(w, r)

				next.ServeHTTP(w, r)
				return
			}

			_, ok = token.Claims.(jwt.MapClaims)
			if !ok {
				fmt.Println("\nno auth3")
				http.Redirect(w, r, "/login", answerCode)
				return
			}

			fmt.Println("auth ok!!!")
			next.ServeHTTP(w, r)
		}
	})
}
