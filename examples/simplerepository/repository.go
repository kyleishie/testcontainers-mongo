package simplerepository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

var (
	errorFirstNameRequired = errors.New("user.firstName is required")
	errorLastNameRequired  = errors.New("user.lastName is required")
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

type User struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

func (r *Repository) CreateUser(ctx context.Context, user User) (id primitive.ObjectID, err error) {

	if user.FirstName == "" {
		return id, errorFirstNameRequired
	}

	if user.LastName == "" {
		return id, errorLastNameRequired
	}

	result, err := r.db.Collection(usersCollection).InsertOne(ctx, user)
	if err != nil {
		return id, err
	}

	id = result.InsertedID.(primitive.ObjectID)
	return
}
