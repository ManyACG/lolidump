package service

import (
	"context"

	"github.com/ManyACG/lolidump/dao"
	"github.com/ManyACG/lolidump/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateFavorite(ctx context.Context, userID, artworkID primitive.ObjectID) (*types.FavoriteModel, error) {
	res, err := dao.CreateFavorite(ctx, userID, artworkID)
	if err != nil {
		return nil, err
	}
	return &types.FavoriteModel{
		ID:        res.InsertedID.(primitive.ObjectID),
		UserID:    userID,
		ArtworkID: artworkID,
	}, nil
}

func GetFavorite(ctx context.Context, userID, artworkID primitive.ObjectID) (*types.FavoriteModel, error) {
	return dao.GetFavorite(ctx, userID, artworkID)
}

func DeleteFavorite(ctx context.Context, userID, artworkID primitive.ObjectID) error {
	_, err := dao.DeleteFavorite(ctx, userID, artworkID)
	return err
}
