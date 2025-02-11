package Repository

import (
	"context"
	"fmt"
	"log"
	dataModel "myapp/DataModel"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection  *mongo.Collection
	collection2 *mongo.Collection
	Client      *mongo.Client
}

func NewRepository() Repository {
	mongocollection, mongocollection2, mongoclient := getMongoConnection()
	return Repository{
		collection:  mongocollection,
		collection2: mongocollection2,
		Client:      mongoclient,
	}
}

func (r Repository) GetTicket(id string) (bson.M, error) {

	var ticket bson.M

	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&ticket)

	return ticket, err

}

func (r Repository) Allticket(status string) (*mongo.Cursor, error) {

	cursor, err := r.collection.Find(context.TODO(), bson.M{"status": status})

	return cursor, err
}

func (r Repository) Myticket(email string) (*mongo.Cursor, error) {

	cursor, err := r.collection.Find(context.TODO(), bson.M{"email": email})

	return cursor, err
}

func (r Repository) AdminResponse(ticket dataModel.Ticket, id string) (dataModel.Ticket, *mongo.UpdateResult, error) {

	var ticket1 dataModel.Ticket
	var result *mongo.UpdateResult = nil
	var err error = nil

	if err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&ticket1); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(ticket1)

		if ticket.Status == "Close" {
			ticket1.Status = "Closed"
		} else if ticket.Status == "Resolve this issue" {
			ticket1.Status = "Resolved"
		} else {
			ticket1.Status = ticket.Status
		}
		ticket1.Response = ticket.Response

		ticket1.Admincontact = ticket.Admincontact
		ticket1.ID = id
		result, err = r.collection.UpdateOne(
			context.TODO(), bson.M{"id": ticket1.ID}, bson.D{
				{"$set", ticket1}})
		if err != nil {

			log.Fatal(err)

		}
	}
	return ticket1, result, err
}

func (r Repository) ReOpenTicket(id string) (dataModel.Ticket, *mongo.UpdateResult, error) {

	var ticket dataModel.Ticket
	var result *mongo.UpdateResult
	result = nil
	var err error

	if err = r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&ticket); err != nil {
		log.Fatal(err)

	} else {
		ticket.Status = "Pending"
		ticket.Response = ""
		ticket.Admincontact = ""

		result, err = r.collection.UpdateOne(
			context.TODO(), bson.M{"id": id}, bson.D{
				{"$set", ticket}})

		if err != nil {

			log.Fatal(err)

		}
	}
	return ticket, result, err
}

func (r Repository) CloseTicket(id string) (dataModel.Ticket, *mongo.UpdateResult, error) {
	var ticket dataModel.Ticket
	var result *mongo.UpdateResult = nil
	var err error

	if err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&ticket); err != nil {
		log.Fatal(err)
	} else {

		ticket.Status = "Closed"

		result, err = r.collection.UpdateOne(
			context.TODO(), bson.M{"id": id}, bson.D{
				{"$set", ticket}})
	}
	return ticket, result, err
}
func (r Repository) Queryas(ticket dataModel.Ticket) (*mongo.InsertOneResult, error) {

	insertResult, err := r.collection.InsertOne(context.TODO(), ticket)
	return insertResult, err
}

func (r Repository) GetUser(email string) (dataModel.User, error) {

	var user dataModel.User
	err := r.collection2.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	fmt.Println(user, err)

	return user, err
}

func getMongoConnection() (*mongo.Collection, *mongo.Collection, *mongo.Client) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {

		log.Fatal(err)

	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {

		log.Fatal(err)

	}

	collection := client.Database("Ticket-Management").Collection("Ticket")
	collection2 := client.Database("Ticket-Management").Collection("user")
	return collection, collection2, client
}
