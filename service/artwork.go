package service

import (
	"github.com/ManyACG/lolidump/adapter"
	"github.com/ManyACG/lolidump/dao"

	"context"

	"github.com/krau/ManyACG/types"

	"github.com/duke-git/lancet/v2/slice"
)

func DumpArtworkToSearch(ctx context.Context) ([]*types.ArtworkSearchDocument, error) {
	cursor, err := dao.GetArtworks(ctx)
	if err != nil {
		return nil, err
	}
	var artworkModels []*types.ArtworkModel
	if err = cursor.All(ctx, &artworkModels); err != nil {
		return nil, err
	}
	artworks, err := adapter.ConvertToArtworks(ctx, artworkModels)
	if err != nil {
		return nil, err
	}
	results := make([]*types.ArtworkSearchDocument, len(artworks))
	for i, artwork := range artworks {
		results[i] = &types.ArtworkSearchDocument{
			ID:          artwork.ID,
			Title:       artwork.Title,
			Artist:      artwork.Artist.Name + " " + artwork.Artist.Username,
			Tags:        slice.Join(artwork.Tags, ","),
			Description: artwork.Description,
			R18:         artwork.R18,
		}
	}
	return results, nil
}
