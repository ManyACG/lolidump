package dao

import (
	"context"
	"time"

	"github.com/krau/ManyACG/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var deletedCollection *mongo.Collection

func GetDeletedByURL(ctx context.Context, sourceURL string) (*types.DeletedModel, error) {
	var deleted types.DeletedModel
	err := deletedCollection.FindOne(ctx, bson.M{"source_url": sourceURL}).Decode(&deleted)
	if err != nil {
		return nil, err
	}
	return &deleted, err
}

func CheckDeletedByURL(ctx context.Context, sourceURL string) bool {
	deleted, err := GetDeletedByURL(ctx, sourceURL)
	if err != nil {
		return false
	}
	return deleted != nil
}

func CreateDeleted(ctx context.Context, deleted *types.DeletedModel) (*mongo.InsertOneResult, error) {
	deleted.DeletedAt = primitive.NewDateTimeFromTime(time.Now())
	return deletedCollection.InsertOne(ctx, deleted)
}

func DeleteDeletedByURL(ctx context.Context, sourceURL string) (*mongo.DeleteResult, error) {
	return deletedCollection.DeleteOne(ctx, bson.M{"source_url": sourceURL})
}
