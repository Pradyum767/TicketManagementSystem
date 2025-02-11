package Repository

import (
	dataModel "myapp/DataModel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryInterface interface {
	GetTicket(id string) (bson.M, error)
	Allticket(status string) (*mongo.Cursor, error)
	Myticket(email string) (*mongo.Cursor, error)
	AdminResponse(ticket dataModel.Ticket, id string) (dataModel.Ticket, *mongo.UpdateResult, error)
	ReOpenTicket(id string) (dataModel.Ticket, *mongo.UpdateResult, error)
	CloseTicket(id string) (dataModel.Ticket, *mongo.UpdateResult, error)
	Queryas(Ticket dataModel.Ticket) (*mongo.InsertOneResult, error)
	GetUser(email string) (dataModel.User, error)
}
