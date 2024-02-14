package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Models struct {
	Db       string
	Coll     string
	Client   *mongo.Client
	LogEntry LogEntry
}

func New(client *mongo.Client, dbName string) *Models {
	return &Models{
		Db:       dbName,
		Coll:     "logs",
		Client:   client,
		LogEntry: LogEntry{},
	}
}

func (m *Models) Insert(entry LogEntry) error {
	coll := m.Client.Database(m.Db).Collection(m.Coll)

	_, err := coll.InsertOne(context.Background(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs: ", err)
		return err
	}
	return nil
}

func (m *Models) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	coll := m.Client.Database(m.Db).Collection(m.Coll)
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("error finding all documents: ", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("error finding all documents,error closing cursor failed: ", err)
			panic(err)
		}
	}(cursor, ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("error decoding logs into slice: ", err)
			return nil, err
		}
		logs = append(logs, &item)
	}

	return logs, nil
}

func (m *Models) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	coll := m.Client.Database(m.Db).Collection(m.Coll)
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var logEntry LogEntry
	err = coll.FindOne(ctx, bson.M{"_id": docID}).Decode(&logEntry)
	if err != nil {
		return nil, err
	}

	return &logEntry, nil
}

func (m *Models) Drop() error {
	log.Println("-------- Dropping collection 'logs' --------")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err := m.Client.Database(m.Db).Collection(m.Coll).Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *Models) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	coll := m.Client.Database(m.Db).Collection(m.Coll)
	docID, err := primitive.ObjectIDFromHex(m.LogEntry.ID)
	if err != nil {
		return nil, err
	}
	updateValue := bson.D{{
		"$set",
		bson.D{
			{"name", m.LogEntry.Name},
			{"data", m.LogEntry.Data},
			{"updated_at", time.Now()},
		},
	}}
	results, err := coll.UpdateOne(ctx, bson.D{{"id", docID}}, updateValue)
	if err != nil {
		return nil, err
	}
	return results, nil
}
