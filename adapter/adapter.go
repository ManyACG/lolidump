package adapter

import (
	"context"

	"github.com/ManyACG/lolidump/dao"
	"github.com/krau/ManyACG/types"
)

func GetArtworkModelTags(ctx context.Context, artworkModel *types.ArtworkModel) ([]string, error) {
	tags := make([]string, len(artworkModel.Tags))
	for i, tagID := range artworkModel.Tags {
		tagModel, err := dao.GetTagByID(ctx, tagID)
		if err != nil {
			return nil, err
		}
		tags[i] = tagModel.Name
	}
	return tags, nil
}

func GetArtworkModelPictures(ctx context.Context, artworkModel *types.ArtworkModel) ([]*types.Picture, error) {
	pictures := make([]*types.Picture, len(artworkModel.Pictures))
	for i, pictureID := range artworkModel.Pictures {
		pictureModel, err := dao.GetPictureByID(ctx, pictureID)
		if err != nil {
			return nil, err
		}
		pictures[i] = pictureModel.ToPicture()
	}
	return pictures, nil
}

func GetArtworkModelIndexPicture(ctx context.Context, artworkModel *types.ArtworkModel) (*types.Picture, error) {
	pictureModel, err := dao.GetPictureByID(ctx, artworkModel.Pictures[0])
	if err != nil {
		return nil, err
	}
	return pictureModel.ToPicture(), nil
}

func ConvertToArtwork(ctx context.Context, artworkModel *types.ArtworkModel, opts ...*types.AdapterOption) (*types.Artwork, error) {
	tags := make([]string, 0)
	var pictures []*types.Picture
	var artist *types.Artist
	var err error
	if len(opts) == 0 {
		tags, err = GetArtworkModelTags(ctx, artworkModel)
		if err != nil {
			return nil, err
		}
		pictures, err = GetArtworkModelPictures(ctx, artworkModel)
		if err != nil {
			return nil, err
		}
		artistModel, err := dao.GetArtistByID(ctx, artworkModel.ArtistID)
		if err != nil {
			return nil, err
		}
		artist = artistModel.ToArtist()
		return &types.Artwork{
			ID:          artworkModel.ID.Hex(),
			Title:       artworkModel.Title,
			Description: artworkModel.Description,
			R18:         artworkModel.R18,
			LikeCount:   artworkModel.LikeCount,
			CreatedAt:   artworkModel.CreatedAt.Time(),
			SourceType:  artworkModel.SourceType,
			SourceURL:   artworkModel.SourceURL,
			Artist:      artist,
			Tags:        tags,
			Pictures:    pictures,
		}, nil
	}
	option := MergeOptions(opts...)
	if option.LoadTag {
		tags, err = GetArtworkModelTags(ctx, artworkModel)
		if err != nil {
			return nil, err
		}
	}
	if option.LoadPicture {
		pictures, err = GetArtworkModelPictures(ctx, artworkModel)
		if err != nil {
			return nil, err
		}
	}
	if option.LoadArtist {
		artistModel, err := dao.GetArtistByID(ctx, artworkModel.ArtistID)
		if err != nil {
			return nil, err
		}
		artist = artistModel.ToArtist()
	}
	if option.OnlyIndexPicture && !option.LoadPicture {
		indexPicture, err := GetArtworkModelIndexPicture(ctx, artworkModel)
		if err != nil {
			return nil, err
		}
		pictures = []*types.Picture{indexPicture}
	}
	return &types.Artwork{
		ID:          artworkModel.ID.Hex(),
		Title:       artworkModel.Title,
		Description: artworkModel.Description,
		R18:         artworkModel.R18,
		CreatedAt:   artworkModel.CreatedAt.Time(),
		SourceType:  artworkModel.SourceType,
		SourceURL:   artworkModel.SourceURL,
		Artist:      artist,
		Tags:        tags,
		Pictures:    pictures,
	}, nil
}

func ConvertToArtworks(ctx context.Context, artworkModels []*types.ArtworkModel, opts ...*types.AdapterOption) ([]*types.Artwork, error) {
	if len(artworkModels) == 1 {
		artwork, err := ConvertToArtwork(ctx, artworkModels[0])
		if err != nil {
			return nil, err
		}
		return []*types.Artwork{artwork}, nil
	}
	artworks := make([]*types.Artwork, len(artworkModels))
	errCh := make(chan error, len(artworkModels))
	for i, artworkModel := range artworkModels {
		go func(i int, artworkModel *types.ArtworkModel) {
			artwork, err := ConvertToArtwork(ctx, artworkModel, opts...)
			if err != nil {
				errCh <- err
				return
			}
			artworks[i] = artwork
			errCh <- nil
		}(i, artworkModel)
	}
	for range artworkModels {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}
	return artworks, nil
}

func OnlyLoadTag() *types.AdapterOption {
	return &types.AdapterOption{
		LoadTag: true,
	}
}

func OnlyLoadArtist() *types.AdapterOption {
	return &types.AdapterOption{
		LoadArtist: true,
	}
}

func OnlyLoadPicture() *types.AdapterOption {
	return &types.AdapterOption{
		LoadPicture: true,
	}
}

func LoadAll() *types.AdapterOption {
	return &types.AdapterOption{
		LoadTag:     true,
		LoadArtist:  true,
		LoadPicture: true,
	}
}

func LoadNone() *types.AdapterOption {
	return &types.AdapterOption{}
}

func MergeOptions(opts ...*types.AdapterOption) *types.AdapterOption {
	result := &types.AdapterOption{}
	for _, opt := range opts {
		if opt.LoadTag {
			result.LoadTag = true
		}
		if opt.LoadArtist {
			result.LoadArtist = true
		}
		if opt.LoadPicture {
			result.LoadPicture = true
		}
		if opt.OnlyIndexPicture {
			result.OnlyIndexPicture = true
		}
	}
	return result
}
