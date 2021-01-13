package controller

import (
	//	"encoding/json"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"muxCRUDAPI/models"
	"muxCRUDAPI/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

//CarController struct...
type CarController struct{
	router *mux.Router
}  

//Init - initialize car routes
func (c *CarController) Init(newRouter *mux.Router){
	c.router = newRouter

	c.router.HandleFunc("/cars/make/{make}", c.getCarsByMake).Methods("GET")
	c.router.HandleFunc("/cars/type/{type}", c.getCarsByType).Methods("GET")
	c.router.HandleFunc("/cars/{id}", c.getCarByID).Methods("GET")
	c.router.HandleFunc("/cars/{id}", c.updateCar).Methods("PUT")
	c.router.HandleFunc("/cars/{id}", c.deleteCar).Methods("DELETE")
	c.router.HandleFunc("/cars/fill", c.addInitialCars).Methods("POST")
	c.router.HandleFunc("/cars", c.getAllCars).Methods("GET")
	c.router.HandleFunc("/cars", c.createCar).Methods("POST")
}

func (c *CarController) getAllCars(res http.ResponseWriter, req *http.Request){
	result, err := mgm.Coll(&model.Car{}).Find(mgm.Ctx(), bson.D{})
	cars := []model.Car{}

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err.Error()})
		return
	}

	result.All(mgm.Ctx(), &cars)
	utilities.SendJSON(200, res, cars)
}

func (c *CarController) getCarsByType(res http.ResponseWriter, req *http.Request){
	carType := mux.Vars(req)["carType"]
	result, err := mgm.Coll(&model.Car{}).Find(mgm.Ctx(), bson.M{"car_type": carType})
	cars := []model.Car{}

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err.Error()})
		return
	}

	result.All(mgm.Ctx(), &cars)
	utilities.SendJSON(200, res, cars)
}

func (c *CarController) getCarsByMake(res http.ResponseWriter, req *http.Request){
	make := mux.Vars(req)["make"]	
	result, err := mgm.Coll(&model.Car{}).Find(mgm.Ctx(), bson.M{"make": make})
	cars := []model.Car{}

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err.Error()})
		return
	}

	result.All(mgm.Ctx(), &cars)
	utilities.SendJSON(http.StatusOK, res, cars)
}

func (c*CarController) addInitialCars(res http.ResponseWriter, req *http.Request){
	jsonFile, err := os.Open("cars.json")

	if err != nil{
		panic(err)
	}

	fmt.Println("Opened json file!")
	defer jsonFile.Close()

	data, _ := ioutil.ReadAll(jsonFile)

	var cars []model.Car
	json.Unmarshal(data, &cars)

	for i := 0; i < len(cars); i++ {
		mgm.Coll(&model.Car{}).Create(&cars[i])
	}

	utilities.SendJSON(http.StatusCreated, res, bson.M{"cars":cars})
}

func (c *CarController) getCarByID(res http.ResponseWriter, req *http.Request){
	id := mux.Vars(req)["id"]
	car := model.Car{}

	err := mgm.Coll(&model.Car{}).FindByID(id, &car)

	if err != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": err.Error()})
		return
	}

	utilities.SendJSON(http.StatusOK, res, car)
}

func (c *CarController) createCar(res http.ResponseWriter, req *http.Request){
	car := model.Car{}
	json.NewDecoder(req.Body).Decode(&car)

	


	mgm.Coll(&car).Create(&car)
	utilities.SendJSON(200, res, bson.M{"message": "to be implemented"})
}

func (c *CarController) updateCar(res http.ResponseWriter, req *http.Request){
	utilities.SendJSON(200, res, bson.M{"message": "to be implemented"})
}

func (c *CarController) deleteCar(res http.ResponseWriter, req *http.Request){
	utilities.SendJSON(200, res, bson.M{"message": "to be implemented"})
}

func (c *CarController) validateCar(car model.Car) bool{
	return  len(strings.ReplaceAll(car.Make, " ", "")) == 0
}