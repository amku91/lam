package repository

import (
	"github.com/amku91/lam/api/app/order/entities"
	"github.com/amku91/lam/api/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Place Order
func PlaceOrder(l entities.Order) (entities.Order, error) {
	c := newOrderCollection()

	defer c.Close()
	return l, c.Session.Insert(l)
}

func TakeOrder(orderID int, o entities.Status) error {
	c := newOrderCollection()
	defer c.Close()

	err := c.Session.Update(bson.M{
		"orderid": orderID,
	}, bson.M{
		"$set": bson.M{
			"status":     o.Status,
			"updated_at": o.UpdatedAt,
		},
	})
	return err
}

// Get One Order
func GetOneOrder(orderid int) (entities.Status, error) {

	var (
		order entities.Status
	)

	c := newOrderCollection()

	defer c.Close()

	return order, c.Session.Find(bson.M{"orderid": orderid}).One(&order)

}

// Get All Order
func GetRowCount() (int, error) {

	c := newOrderCollection()

	defer c.Close()

	return c.Session.Find(nil).Count()

}

// Get All Order
func GetAllOrders(pageNumber int, limit int) ([]entities.OrderList, error) {

	var (
		order []entities.OrderList
	)

	c := newOrderCollection()

	defer c.Close()

	return order, c.Session.Find(nil).Skip(pageNumber*limit).Limit(limit).Sort("_id").All(&order)

}

func newOrderCollection() *mongo.Collection {
	return mongo.NewCollectionSession("order")
}
