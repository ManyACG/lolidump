package service

import (
	"context"

	"github.com/ManyACG/lolidump/dao"
	"github.com/krau/ManyACG/types"
)

func GetRandomTags(ctx context.Context, limit int) ([]string, error) {
	tags, err := dao.GetRandomTags(ctx, limit)
	if err != nil {
		return nil, err
	}
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}

func GetRandomTagModels(ctx context.Context, limit int) ([]*types.TagModel, error) {
	tags, err := dao.GetRandomTags(ctx, limit)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagByName(ctx context.Context, name string) (*types.TagModel, error) {
	return dao.GetTagByName(ctx, name)
}
