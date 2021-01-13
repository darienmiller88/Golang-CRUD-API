package model

import "github.com/kamva/mgm/v3"

//Car model
type Car struct{
	mgm.DefaultModel `bson:",inline"`
	CarName   string  `json:"car_name"   bson:"car_name"`
	Make      string  `json:"make"       bson:"make"`
	ModelYear int     `json:"model_year" bson:"model_year"`
	CarType   string  `json:"car_type"   bson:"car_type"`
}