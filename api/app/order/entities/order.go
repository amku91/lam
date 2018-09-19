package entities

import (
	"gopkg.in/mgo.v2/bson"
)

type Origin struct {
	Latitude  string
	Longitude string
}

type Destination struct {
	Latitude  string
	Longitude string
}

type Order struct {
	ID          *bson.ObjectId `json:"id" bson:"_id"`
	OrderID     int            `json:"-" bson:"orderid"`
	Origin      []string       `json:"origin" bson:"origin"`
	Destination []string       `json:"destination" bson:"destination"`
	Distance    int            `json:"-" bson:"distance"`
	Status      string         `json:"-" bson:"status"`
	CreatedAt   int64          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt   int64          `json:"updated_at,omitempty" bson:"updated_at"`
}

type Status struct {
	Status    string `json:"status" bson:"status"`
	UpdatedAt int64  `json:"updated_at,omitempty" bson:"updated_at"`
}

type OrderList struct {
	ID       int    `json:"id" bson:"orderid"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

func (o Order) IsEmpty() []string {
	var missingParams []string

	if len(o.Origin) != 2 {
		missingParams = append(missingParams, "Origin")
	}

	if len(o.Destination) != 2 {
		missingParams = append(missingParams, "Destination")
	}

	return missingParams
}
