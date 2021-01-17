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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CarController struct...
type CarController struct{
	router *mux.Router
	carTypes []string
}  

//Init - initialize car routes
func (c *CarController) Init(newRouter *mux.Router){
	c.router = newRouter
	c.carTypes = []string{"SUV", "Van", "Sedan", "Coupe", "Pickup", "Wagon", "Hatchback"}

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

	if !c.validateCarType(car.CarType){
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{
			"message": "Car type not supported: " + car.CarType,
			"Types supported: ": c.carTypes, 
		})
		return
	}	

	if len(strings.ReplaceAll(car.CarName, " ", "")) == 0 {
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": "car name required!"})
		return
	}

	mgm.Coll(&car).Create(&car)
	utilities.SendJSON(http.StatusCreated, res, bson.M{"message": car})
}

func (c *CarController) updateCar(res http.ResponseWriter, req *http.Request){
	objectID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	data := bson.M{}
	
	if err != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": err.Error()})
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": err.Error()})
		return
	}

	result, err := mgm.Coll(&model.Car{}).UpdateOne(mgm.Ctx(), bson.M{"_id": objectID}, bson.M{"$set": data})

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err.Error()})
		return
	}

	utilities.SendJSON(200, res, bson.M{"message": result})
}

func (c *CarController) deleteCar(res http.ResponseWriter, req *http.Request){
	objectID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	
	if err != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": err.Error()})
		return
	}

	doc, err := mgm.Coll(&model.Car{}).DeleteOne(mgm.Ctx(), bson.M{"_id": objectID})

	if err != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": err.Error()})
		return
	}

	utilities.SendJSON(200, res, bson.M{"result": doc})
}

func (c *CarController) validateCarType(carType string) bool{
	lower := strings.ToLower(carType)

	for _, t := range c.carTypes{
		if lower == strings.ToLower(t) {
			return true
		}
	}
	return false
}

func (c *CarController) validateCar(car model.Car) bool{
	return len(strings.ReplaceAll(car.Make, " ", "")) == 0
}