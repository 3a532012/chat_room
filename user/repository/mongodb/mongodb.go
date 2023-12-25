package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"user/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (m *mongodbUserRepository) AddFriend(ctx context.Context, userID string, friendID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	friendObjectID, err := primitive.ObjectIDFromHex(friendID)
	if err != nil {
		return err
	}
	match := bson.M{"_id": userObjectID}
	addToSet := bson.M{"$addToSet": bson.M{
		"friends": friendObjectID,
	}}
	result, err := m.collect.UpdateOne(ctx, match, addToSet)
	if err != nil {
		log.Fatal(err)
	}
	if result.ModifiedCount == 0 {
		fmt.Println("Element already exists in the array.")
	} else {
		fmt.Println("Element added to the array.")
	}
	return nil
}

func (m *mongodbUserRepository) RemoveFriend(ctx context.Context, userID string, friendID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	friendObjectID, err := primitive.ObjectIDFromHex(friendID)
	if err != nil {
		return err
	}
	_filter := make([]bson.M, 0)
	match := bson.M{"$match": bson.M{}}
	addToSet := bson.M{"$pull": bson.M{}}
	match["$match"].(bson.M)["_id"] = userObjectID
	addToSet["$pull"].(bson.M)["friends"] = friendObjectID
	_filter = append(_filter, match)
	_filter = append(_filter, addToSet)
	if _, err := m.collect.Aggregate(ctx, _filter); err != nil {
		return err
	}

	return nil
}

func (m *mongodbUserRepository) FindByTagID(ctx context.Context, tagID string) (*domain.User, error) {
	var result domain.User
	filter := bson.M{
		"tag-id": tagID,
	}

	if err := m.collect.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
func (m *mongodbUserRepository) UpdateTagID(ctx context.Context, userID, tagID string) (*domain.User, error) {
	exist, err := m.TagIDExists(ctx, tagID)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("tag id already used")
	}
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objectID,
	}
	update := bson.M{
		"$set": bson.M{
			"tag-id": tagID,
		},
	}
	var updatedUser domain.User
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	m.collect.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	return &updatedUser, nil
}
func (m *mongodbUserRepository) TagIDExists(ctx context.Context, tagID string) (bool, error) {
	filter := bson.M{
		"tag-id": tagID,
	}

	count, err := m.collect.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
