package main

import (
	"database/sql"
	"events/pkg/handlers"
	"events/pkg/middleware"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("can't zap.NewProduction()")
		return
	}
	defer func() {
		err = zapLogger.Sync()
		fmt.Println("can't  zapLogger.Sync()")
	}()
	//logger := zapLogger.Sugar()

	r := mux.NewRouter()

	//dbPort:="8889"
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

	EventsRrepo := handlers.NewEventsHand(db)
	userRepo := handlers.NewUserHand(db)
	TitlesRepo := handlers.NewTitlesHand(db)
	TaskRepo := handlers.NewTaskHandler(db)

	r.HandleFunc("/api/login", userRepo.Login).Methods("POST")
	r.HandleFunc("/api/signUp", userRepo.SingUp).Methods("POST")
	r.HandleFunc("/api/checkCode", userRepo.CheckCode).Methods("POST")

	r.HandleFunc("/api/newEvent", EventsRrepo.AddEvent).Methods("POST")
	r.HandleFunc("/api/getEvents", EventsRrepo.GetEvents).Methods("GET")
	r.HandleFunc("/api/changeEvent", EventsRrepo.ChangeEvent).Methods("PUT")
	r.HandleFunc("/api/deleteEvent", EventsRrepo.DeleteEvent).Methods("DELETE")

	r.HandleFunc("/api/newTitle", TitlesRepo.AddTitle).Methods("POST")
	r.HandleFunc("/api/getTitles", TitlesRepo.GetTitles).Methods("GET")
	r.HandleFunc("/api/deleteTitle", TitlesRepo.DeleteTitle).Methods("DELETE")

	r.HandleFunc("/api/newTask", TaskRepo.AddTask).Methods("POST")
	r.HandleFunc("/api/getTasks", TaskRepo.GetTasks).Methods("GET")
	r.HandleFunc("/api/changeTask", TaskRepo.ChangeTask).Methods("PUT")
	r.HandleFunc("/api/deleteTask", TaskRepo.DeleteTask).Methods("DELETE")

	r.HandleFunc("/api/getTaskPieces", TaskRepo.GetTasksPiece).Methods("POST") // post - becuase get have no body
	r.HandleFunc("/api/deleteTaskPiece", TaskRepo.DeleteTaskPiece).Methods("DELETE")

	spa := handlers.SpaHandler{StaticPath: ".././static/dist", IndexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	handler := middleware.Auth(r)

	port := "8080"
	fmt.Println("start serv on port " + port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("can't Listen and server")
		return
	}

}
