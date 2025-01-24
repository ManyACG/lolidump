package dao

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ManyACG/lolidump/common"
	"github.com/ManyACG/lolidump/config"
	"github.com/ManyACG/lolidump/dao/collections"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var DB *mongo.Database

func InitDB(ctx context.Context) {
	common.Logger.Info("Initializing database...")
	uri := config.Cfg.Database.URI
	if uri == "" {
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d",
			config.Cfg.Database.User,
			config.Cfg.Database.Password,
			config.Cfg.Database.Host,
			config.Cfg.Database.Port,
		)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri),
		options.Client().SetReadPreference(readpref.Nearest(readpref.WithMaxStaleness(time.Duration(config.Cfg.Database.MaxStaleness)*time.Second))))
	if err != nil {
		common.Logger.Fatal(err)
		os.Exit(1)
	}
	if err = client.Ping(ctx, nil); err != nil {
		common.Logger.Fatal(err)
		os.Exit(1)
	}
	Client = client
	DB = Client.Database(config.Cfg.Database.Database)
	if DB == nil {
		common.Logger.Fatal("Failed to get database")
		os.Exit(1)
	}
	createCollection(ctx)

	common.Logger.Info("Database initialized")
}

func createCollection(ctx context.Context) {
	for _, collection := range collections.AllCollections {
		DB.CreateCollection(ctx, collection)
	}

	artworkCollection = DB.Collection(collections.Artworks)

	tagCollection = DB.Collection(collections.Tags)

	pictureCollection = DB.Collection(collections.Pictures)

	artistCollection = DB.Collection(collections.Artists)

	adminCollection = DB.Collection(collections.Admins)

	deletedCollection = DB.Collection(collections.Deleted)

	callbackDataCollection = DB.Collection(collections.CallbackData)

	cachedArtworkCollection = DB.Collection(collections.CachedArtworks)

	etcDataCollection = DB.Collection(collections.EtcData)

	userCollection = DB.Collection(collections.Users)

	likeCollection = DB.Collection(collections.Likes)

	favoriteCollection = DB.Collection(collections.Favorites)

	unauthUserCollection = DB.Collection(collections.UnauthUser)

	apiKeyCollection = DB.Collection(collections.ApiKeys)
}
