package testing

import (
	"gopkg.in/h2non/baloo.v3"
	"testing"
	"fmt"
	"encoding/json"
	"strconv"
)

// test stores the HTTP testing client preconfigured
var server = baloo.New("http://localhost:8080")

type OrderID struct {
	ID int `json:"id"`
}


//Place Order Testing

func Test_PlaceOrder_Success_Status_200(t *testing.T) {
	server.Post("/order").
		JSON(map[string][]string{"origin": {"12.9678661", "78.65"}, "destination": {"12.9726869", "78.63316544"}}).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func Test_PlaceOrder_Failed_Without_Geocode_ResponseCode_500(t *testing.T) {
	server.Post("/order").
		JSON([]string{"origin", "destination"}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_PlaceOrder_Failed_GeoCode_KeyMismatch_ResponseCode_500(t *testing.T) {
	server.Post("/order").
		JSON(map[string][]string{"originBB": {"12.9678661", "78.65BBB"}, "destinationCC": {"12.9726869", "78.63316544CC"}}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_PlaceOrder_Failed_Invalid_GeoCode_ResponseCode_500(t *testing.T) {
	server.Post("/order").
		JSON(map[string][]string{"origin": {"12.9678661", "78.65BBB"}, "destination": {"12.9726869", "78.63316544CC"}}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_PlaceOrder_Failed_Invalid_GeoCode_Range_ResponseCode_500(t *testing.T) {
	server.Post("/order").
		JSON(map[string][]string{"origin": {"12.9678661", "190.65"}, "destination": {"12.9726869", "-278.63316"}}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}


//Take Order Testing

//Get Order ID function

func runPlaceOrder() int {

	var (
		response OrderID
	)
	a, err := server.Post("/order").
		JSON(map[string][]string{"origin": {"12.9678661", "78.65"}, "destination": {"12.9726869", "78.63316"}}).
		Send()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = json.Unmarshal(a.Bytes(), &response)
	if err != nil {
		fmt.Println(err.Error())
	}
	return response.ID
}


func Test_TakeOrder_Success_Status_200(t *testing.T) {
	//Run Post first
	orderID := runPlaceOrder()
	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status": "taken"}).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func Test_TakeOrder_Success_With_Schema_Validation_ResponseCode_200(t *testing.T) {
	//Run Place Order first
	orderID := runPlaceOrder()
	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status": "taken"}).
		Expect(t).
		Status(200).
		Type("json").
		JSON(map[string]string{"status": "SUCCESS"}).//Check response json
		Done()
}

func Test_TakeOrder_Failed_Already_Taken_Status_Value_ResponseCode_409(t *testing.T) {
	//First place an order and take also
	orderID := runPlaceOrder()

	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status": "taken"}).
		Expect(t).
		Status(200).
		Type("json").
		JSON(map[string]string{"status": "SUCCESS"}).//Check response json
		Done()

	//Now fire for same request
	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status":"taken"}).
		Expect(t).
		Status(409).
		Type("json").
		Done()
}

func Test_TakeOrder_Failed_Without_Status_Value_ResponseCode_500(t *testing.T) {
	orderID := runPlaceOrder()

	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status":""}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_TakeOrder_Failed_Method_Not_Allowed_Value_ResponseCode_405(t *testing.T) {

	server.Put("/order/").
		JSON(map[string]string{"status":"taken"}).
		Expect(t).
		Status(405).
		Done()
}

func Test_TakeOrder_Failed_Invalid_Status_ResponseCode_500(t *testing.T) {
	orderID := runPlaceOrder()

	server.Put("/order/" + strconv.Itoa(orderID)).
		JSON(map[string]string{"status":"AAABBB"}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_TakeOrder_Failed_Invalid_Order_ID_ResponseCode_500(t *testing.T) {

	server.Put("/order/987654" ).//Just for testing purpose
		JSON(map[string]string{"status":"taken"}).
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

//Get Order List Testing

func Test_OrderList_Success_Status_200(t *testing.T) {
	server.Get("/order").
		AddQuery("page", "1").
		AddQuery("limit", "1").
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Safety_Limit_Error_Status_500(t *testing.T) {
	server.Get("/order").
		AddQuery("page", "1").
		AddQuery("limit", "120").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Zero_Limit_Status_500(t *testing.T) {
	server.Get("/order").
		AddQuery("page", "1").
		AddQuery("limit", "0").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Zero_PageNumber_Status_500(t *testing.T) {
	server.Get("/order").
		AddQuery("page", "0").
		AddQuery("limit", "0").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Without_PageNumber_Status_500(t *testing.T) {
	server.Get("/order").
		AddQuery("limit", "0").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Without_Limit_Status_500(t *testing.T) {
	server.Get("/order").
		AddQuery("page", "0").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}

func Test_OrderList_Failed_Without_Query_Params_Status_500(t *testing.T) {
	server.Get("/order").
		Expect(t).
		Status(500).
		Type("json").
		Done()
}