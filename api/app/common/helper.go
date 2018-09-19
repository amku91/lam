package common

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amku91/lam/api/app/common/delivery"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func HandleAPIError(code int, errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(code)
	resp, marshError := GetErrorResponse(code, errorMessage)
	if marshError != nil {
		if marshError != nil {

			log.Println("Helper: HandleAPIError: Something went wrong while marshalling the response: " + marshError.Error())

			_, respErr := w.Write([]byte("Fatal error!"))

			if respErr != nil {

				log.Println("Helper: HandleAPIError: Something went wrong while writing the response object: " + marshError.Error())
			}
			return
		}

	}

	_, respErr := w.Write(resp)

	if respErr != nil {

		log.Println("Helper: HandleAPIError: Something went wrong while writing the response object: " + marshError.Error())
	}
	return
}

// GetErrorResponse function to return the  error response
func GetErrorResponse(code int, errorMessage string) (resp []byte, err error) {

	result := delivery.ErrorResponse{Error: errorMessage}

	respBody, marshError := json.MarshalIndent(result, "", "  ")

	if marshError != nil {
		log.Println("Helper: GetErrorResponse: Error while marshaling the response: " + marshError.Error())
	}

	return respBody, marshError
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func ConvertToObjectIDHex(id string) (result bson.ObjectId, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unable to convert %v to object id", id)
		}
	}()

	return bson.ObjectIdHex(id), err
}
