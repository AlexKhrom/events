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

	//bdPort:="8889"
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
	//BagRepo := handlers.BagHandlers{}
	//ProductRepo := handlers.NewProductHand(db)
	userRepo := handlers.NewUserHand(db)
	TitlesRepo := handlers.NewTitlesHand(db)
	TaskRepo := handlers.NewTaskHandler(db)
	//OrdersRepo := handlers.NewOrdersHand(db)
	////SearchPageRepo := handlers.NewSearchPageHand(db)
	//
	//r.HandleFunc("/api/putInBag/{PRODUCT_ID}/{NUM_PRODUCT}", BagRepo.AddToBag).Methods("POST")
	//
	//r.HandleFunc("/api/order", OrdersRepo.GetOrder).Methods("GET")                  // get all orders
	//r.HandleFunc("/api/order", OrdersRepo.NewOrder).Methods("POST")                 // creat new order
	//r.HandleFunc("/api/order", OrdersRepo.ChangeOrderStatus).Methods("PUT")         // change order-status
	//r.HandleFunc("/api/order/{ORDER_ID}", OrdersRepo.DeleteOrder).Methods("DELETE") // delete order
	//r.HandleFunc("/api/orders", OrdersRepo.GetOrders).Methods("POST")               // get some orders
	//
	//r.HandleFunc("/api/giveProductUrl/{PRODUCT_ID}", ProductRepo.GetProduct).Methods("POST")
	//

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
	r.HandleFunc("/api/changeTask", EventsRrepo.ChangeEvent).Methods("PUT")
	r.HandleFunc("/api/deleteTask", TitlesRepo.DeleteTitle).Methods("DELETE")

	spa := handlers.SpaHandler{StaticPath: ".././static/dist", IndexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	//fs := http.FileServer(http.Dir(".././static/dist"))
	//http.Handle("/", fs)

	//handler := middleware.Auth(r)
	handler := middleware.Auth(r)
	//handler = middleware.AccessLog(logger, r)

	// put words to db - wordToId - for test
	//items.PutWords(db)
	//items.PrintWords(db)
	//items.GetOneWord("legwyldebf", db)

	port := "8085"
	fmt.Println("start serv on port " + port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("can't Listen and server")
		return
	}

}
