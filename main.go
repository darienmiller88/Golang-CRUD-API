package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"muxCRUDAPI/controllers"

	"github.com/gorilla/mux"
	//"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main(){
	//Load env file
	godotenv.Load()
	openDatabase()
	
	router := mux.NewRouter()
	userController := controller.UserController{}
	carController := controller.CarController{}

	router.Use(mux.CORSMethodMiddleware(router))
	userController.Init(router.PathPrefix("/api").Subrouter())
	carController.Init(router.PathPrefix("/api").Subrouter())

	fmt.Println("listening on port", os.Getenv("PORT"), "...")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router))
}

func openDatabase(){
	err := mgm.SetDefaultConfig(nil, os.Getenv("DB_NAME"), options.Client().ApplyURI(os.Getenv("DB_CONNECTION")))

	if err != nil{
		panic(err)
	}
}