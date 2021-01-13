package model

import (
	"github.com/kamva/mgm/v3"
)

//User struct...
type User struct{
	mgm.DefaultModel `bson:",inline"`
	Name            string `json:"name" bson:"name"`
	FavoriteCarType string `json:"favorite_car_type" bson:"favorite_car_type"`
	FavoriteCar     string `json:"favorite_car" bson:"favorite_car"`
	FavoriteCars    []Car  `json:"favorite_cars" bson:"favorite_cars"`
}