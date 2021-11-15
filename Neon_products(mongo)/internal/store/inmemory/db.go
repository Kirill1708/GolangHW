package inmemory

import (
	"context"
	"errors"
	"fmt"
	"log"
	"neon_products/internal/models"
	"neon_products/internal/store"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	collection *mongo.Collection
	mu         *sync.RWMutex
}

func Init() store.Store {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	return &DB{
		collection: client.Database("neonProducts").Collection("products"),
		mu:         new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, neon *models.Neon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.collection.InsertOne(ctx, neon)
	return err
}

func (db *DB) All(ctx context.Context) ([]*models.Neon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	filter := bson.D{{}}

	return db.filterTasks(ctx, filter)
}

func (db *DB) ByID(ctx context.Context, id string) (*models.Neon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("conversion of id from string to ObjectID: %s", id)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: idObj}}
	u := &models.Neon{}
	ok := db.collection.FindOne(ctx, filter).Decode(u)
	if ok != nil {
		return nil, fmt.Errorf("no neon with id %s", id)
	}
	return u, nil
}

func (db *DB) Update(ctx context.Context, neon *models.Neon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	filter := bson.D{primitive.E{Key: "_id", Value: neon.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "title", Value: neon.Title},
		primitive.E{Key: "sizes", Value: neon.Sizes},
		primitive.E{Key: "colors", Value: neon.Colors},
		primitive.E{Key: "cost", Value: neon.Cost},
		primitive.E{Key: "theme", Value: neon.Theme},
	}}}

	u := &models.Neon{}

	return db.collection.FindOneAndUpdate(ctx, filter, update).Decode(u)
}

func (db *DB) Delete(ctx context.Context, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("conversion of id from string to ObjectID: %s", id)
	}
	filter := bson.D{primitive.E{Key: "_id", Value: idObj}}

	res, err := db.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no tasks were deleted")
	}

	return nil
}

func (db *DB) filterTasks(ctx context.Context, filter interface{}) ([]*models.Neon, error) {
	var neons []*models.Neon

	cur, err := db.collection.Find(ctx, filter)
	if err != nil {
		return neons, err
	}

	for cur.Next(ctx) {
		var t models.Neon
		err := cur.Decode(&t)
		if err != nil {
			return neons, err
		}

		neons = append(neons, &t)
	}

	if err := cur.Err(); err != nil {
		return neons, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	if len(neons) == 0 {
		return neons, mongo.ErrNoDocuments
	}

	return neons, nil
}
