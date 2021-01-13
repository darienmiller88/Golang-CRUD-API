package controller

import (
	"encoding/json"
	"fmt"
	"muxCRUDAPI/models"
	"muxCRUDAPI/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UserController struct...
type UserController struct{
	router *mux.Router
	carTypes []string
}

//Init - initalize routes
func (u *UserController) Init(newRouter *mux.Router){
	u.router = newRouter
	u.carTypes = []string{"SUV", "Van", "Sedan", "Coupe", "Pickup", "Wagon", "Hatchback"}

	// static := http.FileServer(http.Dir("./static"))
	// http.Handle("/",static)
	u.router.HandleFunc("/users/carType/{carType}", u.getUsersByFavoriteCarType).Methods("GET")
	u.router.HandleFunc("/users/{id}", u.getUserByID).Methods("GET")
	u.router.HandleFunc("/users/{id}", u.updateUser).Methods("PUT")
	u.router.HandleFunc("/users/{id}", u.deleteUser).Methods("DELETE")
	u.router.HandleFunc("/users", u.createUser).Methods("POST")	
	u.router.HandleFunc("/users", u.getAllUsers).Methods("GET")	

}

//CREATE - Method to create a user
func (u *UserController) createUser(res http.ResponseWriter, req *http.Request){
	user := model.User{}

	//Decode the request body into the user body
	json.NewDecoder(req.Body).Decode(&user)	

	//First, validate the user's preferred car type
	if !u.validateCarType(user.FavoriteCarType){
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": "Car type is not valid", "types allowed": u.carTypes})
		return
	}

	//Also validate their name to make sure it's not empty
	if len(strings.ReplaceAll(user.Name, " ", "")) == 0{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": "You must include a name field!"})
		return
	}

	//Add the user object to the database 
	err := mgm.Coll(&user).Create(&user)

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err})
		return
	}

	fmt.Println(user)

	//Finally, return the document that was just recently added
	utilities.SendJSON(http.StatusCreated, res, user)
}

//READ - Get all the users that have been created
func (u *UserController) getAllUsers(res http.ResponseWriter, req *http.Request){
	documents, err := mgm.Coll(&model.User{}).Find(mgm.Ctx(), bson.D{})
	users := []model.User{}

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err})
		return
	}

	errHandle(documents.All(mgm.Ctx(), &users))
	utilities.SendJSON(200, res, users)
}

//READ - Get a user by their id
func (u *UserController) getUserByID(res http.ResponseWriter, req *http.Request){
	id := mux.Vars(req)["id"]
	user := model.User{}

	//find all the user in the "users" collection who have "id" as their favrite car type
	result := mgm.Coll(&model.User{}).FindByID(id, &user)

	if result != nil{
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message":  result.Error()})
		return
	}

	utilities.SendJSON(200, res, user)
}

func (u *UserController) getUsersByFavoriteCarType(res http.ResponseWriter, req *http.Request){
	carType := mux.Vars(req)["carType"]
	users := []model.User{}

	//find all the users in the "users" collection who have "carType" as their favrite car type
	documents, err := mgm.Coll(&model.User{}).Find(mgm.Ctx(), bson.M{"favorite_car_type": carType})

	if !u.validateCarType(carType){
		utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message": "Car type is not valid", "types allowed": u.carTypes})
		return
	}

	if err != nil{
		utilities.SendJSON(http.StatusInternalServerError, res, bson.M{"message": err})
		return
	}

	//Unmarshall documents into array
	documents.All(mgm.Ctx(), &users)
	utilities.SendJSON(200, res, users)
}

func (u *UserController) updateUser(res http.ResponseWriter, req *http.Request){
	objectID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	user := bson.M{}

	if err != nil{
	   utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message":  "invalid id"})
	  return
	}

	json.NewDecoder(req.Body).Decode(&user)
	result, err := mgm.Coll(&model.User{}).UpdateOne(mgm.Ctx(), bson.M{"_id": objectID}, bson.M{"$set": user})

	if err != nil{
		utilities.SendJSON(http.StatusNotFound, res, bson.M{"message":  err})
		return
	}

	utilities.SendJSON(200, res, utilities.M{"result": result})
}

func (u *UserController) deleteUser(res http.ResponseWriter, req *http.Request){
	objectID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])

	if err != nil{
	   utilities.SendJSON(http.StatusBadRequest, res, bson.M{"message":  "invalid id: " + objectID.String()})
	  return
	}
	
	result := mgm.Coll(&model.User{}).FindOneAndDelete(mgm.Ctx(), bson.M{"_id": objectID})
	document := &model.User{}

	result.Decode(&document)
	utilities.SendJSON(200, res, bson.M{"deleted": document})
}

func (u *UserController) validateCarType(carType string) bool{
	lower := strings.ToLower(carType)

	for _, t := range u.carTypes{
		if lower == strings.ToLower(t) {
			return true
		}
	}
	return false
}

func errHandle(err error){
	if err != nil {
		panic(err)
	}
}