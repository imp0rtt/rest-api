package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rest-api/internal/apperror"
	"rest-api/internal/user"
	"rest-api/pkg/logging"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	//nCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}
	d.logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. oid: %s", oid)
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error:%v", err)
	}

	if err := cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error: %v", err)
	}

	return u, nil
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID. hex: %s", id)
	}
	// mongo.getDataBase("test").getCollection("docs").find({})
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrorNoFound
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error:%v", id, err)
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from db: user(%s) due to error:%v", id, err)
	}

	return u, nil
}

func (d *db) FindByEmail(ctx context.Context, email string) (u user.User, err error) {
	oid := email
	//if err != nil {
	//	return u, fmt.Errorf("failed to convert hex to ObjectID. hex: %s", email)
	//}
	// mongo.getDataBase("test").getCollection("docs").find({})
	filter := bson.M{"email": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrorNoFound
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error:%v", email, err)
	}
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from db: user(%s) due to error:%v", email, err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectId, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID, ID=%s", objectId)
	}
	filter := bson.M{"_id": objectId}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal userBytes. error: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrorNoFound
	}

	d.logger.Tracef("Matched %d documnets and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID, ID=%s", objectId)
	}
	filter := bson.M{"_id": objectId}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}

	if result.DeletedCount == 0 {
		return apperror.ErrorNoFound
	}

	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
