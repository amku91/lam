package order

import (
	"encoding/json"
	"github.com/amku91/lam/api/app/order/delivery"
	"github.com/amku91/lam/api/app/order/entities"
	"github.com/amku91/lam/api/app/order/repository"
	"github.com/go-chi/chi"
	"github.com/go-ozzo/ozzo-validation"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/amku91/lam/api/app/common"
	"github.com/go-ozzo/ozzo-validation/is"
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"github.com/amku91/lam/api/config"
)

// Controller controller for order
type Controller struct {
}

// Routes routes for order
func (rs Controller) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.PlaceOrder)
	r.Put("/{id}", rs.TakeOrder)
	r.Get("/", rs.OrderList)

	return r
}

//Place Order
func (rs Controller) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var (
		order entities.Order
		err   error
	)

	//Set the response type to JSON

	w.Header().Set("Content-Type", "application/json")

	// Decode the post data onto the order struct
	err = json.NewDecoder(r.Body).Decode(&order)

	if err != nil {

		common.HandleAPIError(500, "Decode error. Please pass the required request parameters.", w)

		return
	}

	// Check if the struct is valid
	missingParams := order.IsEmpty()

	if len(missingParams) != 0 {

		common.HandleAPIError(500, "Invalid request body. Invalid parameters: "+strings.Join(missingParams, ","), w)

		return
	}

	// Assign to latitude and longitude
	originset := entities.Origin{order.Origin[0], order.Origin[1]}
	destinationset := entities.Destination{order.Destination[0], order.Destination[1]}

	// Check if request contains valid latitude and longitude
	err = validation.Errors{
		"Origin Latitude": validation.Validate(originset.Latitude, validation.Required.Error("Latitude is required"),
			is.Latitude.Error("Invalid latitude format")),
		"Origin Longitude": validation.Validate(originset.Longitude, validation.Required.Error("Longitude is required"),
			is.Longitude.Error("Invalid longitude format")),
		"Destination Latitude": validation.Validate(destinationset.Latitude, validation.Required.Error("Latitude is required"),
			is.Latitude.Error("Invalid latitude format")),
		"Destination Longitude": validation.Validate(destinationset.Longitude, validation.Required.Error("Longitude is required"),
			is.Longitude.Error("Invalid longitude format")),
	}.Filter()

	if err != nil {

		common.HandleAPIError(500, err.Error(), w)
		return
	}

	// Add some static values
	order.CreatedAt = time.Now().Unix()

	//Build latitude,longitude according to google maps client library
	origin := [] string{originset.Latitude + "," + originset.Longitude}

	destination := [] string{destinationset.Latitude + "," + destinationset.Longitude}

	connection, err := maps.NewClient(maps.WithAPIKey(config.GOOGLE_API_KEY))

	if err != nil {
		common.HandleAPIError(500, "Unable to connect with google maps", w)
		return
	}

	dmr := &maps.DistanceMatrixRequest{
		Origins:      origin,
		Destinations: destination,
		Mode:         "driving",
		Language:     config.GOOGLE_LANGUAGE,
	}
	// Init google maps distance matrix API
	distance, err := connection.DistanceMatrix(context.Background(), dmr)

	if err != nil {
		common.HandleAPIError(500, "Something went wrong with google maps API", w)
		return
	}

	ObjectID := bson.NewObjectId()
	order.ID = &ObjectID

	order.Status = config.UNASSIGN_ORDER_STATUS
	order.Distance = distance.Rows[0].Elements[0].Distance.Meters

	//Get row Count

	rowCount, err := repository.GetRowCount()

	order.OrderID = rowCount + 1

	if err != nil {

		common.HandleAPIError(500, "Something went wrong while generating order id. Please try again later."+err.Error(), w)

		return
	}

	// Place Order DB Call.
	_, err = repository.PlaceOrder(order)

	//if error, return the error response
	if err != nil {

		common.HandleAPIError(500, "Something went wrong while placing order. Please try again later.", w)

		return
	}

	// initialize the response struct
	productresponse := delivery.PlaceOrderResponse{}
	productresponse = delivery.PlaceOrderResponse{order.OrderID, order.Distance, order.Status}

	//return the success response
	resp := productresponse

	respBody, marshError := json.MarshalIndent(resp, "", "  ")

	if marshError != nil {

		log.Println("OrderController: PlaceOrder: Something went wrong while marshalling the response: " + marshError.Error())

		_, respErr := w.Write([]byte("Fatal error!"))

		if respErr != nil {

			log.Println("OrderController: PlaceOrder: Something went wrong while writing the response object: " + marshError.Error())
		}
		return
	}

	_, respErr := w.Write(respBody)

	if respErr != nil {

		log.Println("OrderController: PlaceOrder: Something went wrong while writing the response object: " + marshError.Error())
	}
}

//Take Order
func (rs Controller) TakeOrder(w http.ResponseWriter, r *http.Request) {
	var (
		order   entities.Status
		orderID int
		err     error
	)

	//Set the response type to JSON

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	orderID, err = strconv.Atoi(id)

	//If error, return error response
	if err != nil {

		common.HandleAPIError(500, "Invalid order id format", w)

		return
	}

	order, err = repository.GetOneOrder(orderID)

	if err != nil {

		common.HandleAPIError(500, "No order found", w)

		return
	}

	if order.Status == config.TAKEN_ORDER_STATUS_CL {

		common.HandleAPIError(409, "ORDER_ALREADY_BEEN_TAKEN", w)

		return
	}
	// Decode the post data onto the order struct
	err = json.NewDecoder(r.Body).Decode(&order)

	if err != nil {

		common.HandleAPIError(500, "Decode error. Please pass the required request parameters.", w)

		return
	}

	// Check if request contains valid status
	err = validation.Errors{
		"Status": validation.Validate(order.Status, validation.Required.Error("Status is required")),
	}.Filter()

	if err != nil {

		common.HandleAPIError(500, err.Error(), w)
		return
	}

	if order.Status != config.TAKEN_ORDER_STATUS_SL {
		common.HandleAPIError(500, "Order status is invalid", w)
		return
	}

	// Add some static values
	order.UpdatedAt = time.Now().Unix()
	order.Status = config.TAKEN_ORDER_STATUS_CL

	// Place Order DB Call.
	err = repository.TakeOrder(orderID, order)

	//if error, return the error response
	if err != nil {

		common.HandleAPIError(500, "Something went wrong while taking order. Please try again later.", w)

		return
	}

	// initialize the response struct
	productresponse := delivery.TakeOrderResponse{}
	productresponse = delivery.TakeOrderResponse{"SUCCESS"}

	//return the success response
	resp := productresponse

	respBody, marshError := json.MarshalIndent(resp, "", "  ")

	if marshError != nil {

		log.Println("OrderController: TakeOrder: Something went wrong while marshalling the response: " + marshError.Error())

		_, respErr := w.Write([]byte("Fatal error!"))

		if respErr != nil {

			log.Println("OrderController: TakeOrder: Something went wrong while writing the response object: " + marshError.Error())
		}
		return
	}

	_, respErr := w.Write(respBody)

	if respErr != nil {

		log.Println("OrderController: TakeOrder: Something went wrong while writing the response object: " + marshError.Error())
	}
}

//Order List
func (rs Controller) OrderList(w http.ResponseWriter, r *http.Request) {
	var (
		order      []entities.OrderList
		pageNumber int
		limit      int
		err        error
		resp       interface{}
	)

	//Set the response type to JSON

	w.Header().Set("Content-Type", "application/json")

	pageno := r.FormValue("page")

	lim := r.FormValue("limit")

	safetyLimit := config.RECORDS_SAFETY_LIMIT

	pageNumber, err = strconv.Atoi(pageno)

	if err != nil {

		common.HandleAPIError(500, "Invalid page number", w)
		return
	}

	limit, err = strconv.Atoi(lim)

	if err != nil {

		common.HandleAPIError(500, "Invalid limit", w)
		return
	}

	if pageNumber < 1 {

		common.HandleAPIError(500, "Page number should be minimum 1", w)
		return
	}

	if limit < 1 {

		common.HandleAPIError(500, "Limit should be minimum 1", w)
		return
	}

	if limit > safetyLimit {
		common.HandleAPIError(500, "Limit is too high. Max limit is : "+strconv.Itoa(safetyLimit), w)
		return
	}

	//If page is first then start from zero entry
	if pageNumber == 1 {
		pageNumber = 0
	}

	order, err = repository.GetAllOrders(pageNumber, limit)

	if err != nil {

		common.HandleAPIError(500, "No order found", w)

		return
	}
	//Initilize empty array if no order left
	emptyResponse := make([]string, 0)

	if order == nil {
		resp = emptyResponse
	} else {
		resp = order
	}

	respBody, marshError := json.MarshalIndent(resp, "", "  ")

	if marshError != nil {

		log.Println("OrderController: OrderList: Something went wrong while marshalling the response: " + marshError.Error())

		_, respErr := w.Write([]byte("Fatal error!"))

		if respErr != nil {

			log.Println("OrderController: OrderList: Something went wrong while writing the response object: " + marshError.Error())
		}
		return
	}
	_, respErr := w.Write(respBody)

	if respErr != nil {

		log.Println("OrderController: OrderList: Something went wrong while writing the response object: " + marshError.Error())
	}
}
