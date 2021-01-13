package utilities

import(
	"net/http"
	"encoding/json"
)

//M - Type def for map[string]interface{}
type M map[string]interface{} 

//SendJSON - Utility function to send JSON
func SendJSON(statusCode int, res http.ResponseWriter, body interface{}){
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(body)
	
	return
}