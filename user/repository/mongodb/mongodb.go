package mongodb

import (
	"context"
	"log"

	"user/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbUserRepository struct {
	db      *mongo.Client
	collect *mongo.Collection
}

// NewMongodbUserRepository ...
func NewMongodbUserRepository(db *mongo.Client) domain.UserRepository {
	collect := db.Database("chatroom").Collection("users")
	return &mongodbUserRepository{
		db:      db,
		collect: collect,
	}
}

func (m *mongodbUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var result domain.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objectID,
	}
	err = m.collect.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
func (m *mongodbUserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	var result domain.User

	filter := bson.M{
		"name": name,
	}
	err := m.collect.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
func (m *mongodbUserRepository) Store(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := m.collect.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}
func (m *mongodbUserRepository) All(ctx context.Context) ([]*domain.User, error) {
	var result []*domain.User
	filter := bson.M{}
	cursor, err := m.collect.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each message
	if err = cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	return result, nil
}
