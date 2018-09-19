package config

import (
	"strconv"
)


const(

	GOOGLE_API_KEY =  "AIzaSyB1H_NgAH9e-PWEyqywbTwR5uzIoY3jZ0c"

	GOOGLE_LANGUAGE = "en-EN"

	UNASSIGN_ORDER_STATUS =  "UNASSIGN"

	TAKEN_ORDER_STATUS_CL = "TAKEN"

	TAKEN_ORDER_STATUS_SL = "taken"

	RECORDS_SAFETY_LIMIT = 100

	MONGO_MAX_POOL = 4000

	MONGO_DSN = "mongodb://admin:admin123@lam_mongodb_1/lam"

	MONGO_DATABASE = "lam"
)


func OverrideAtoI(input string) int{
	number, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return number
}