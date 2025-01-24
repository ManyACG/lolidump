package cmd

import (
	"context"
	"os"
	"time"

	"github.com/ManyACG/lolidump/common"
	"github.com/ManyACG/lolidump/config"
	"github.com/ManyACG/lolidump/dao"
	"github.com/ManyACG/lolidump/service"

	_ "net/http/pprof"

	"github.com/bytedance/sonic"
	"github.com/meilisearch/meilisearch-go"
)

func Run() {
	config.InitConfig()
	common.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	common.Logger.Info("Start github.com/ManyACG/lolidump...")
	dao.InitDB(ctx)
	defer func() {
		if err := dao.Client.Disconnect(ctx); err != nil {
			common.Logger.Fatal(err)
			os.Exit(1)
		}
	}()

	switch config.Cfg.Dest.Type {
	case "meilisearch":
		dumpToMeilisearch()
	default:
		common.Logger.Fatal("Unsupported destination type")
	}

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// sig := <-quit
	// common.Logger.Info(sig, " Exiting...")
	// defer common.Logger.Info("Exited.")
	// if err := service.Cleanup(context.TODO()); err != nil {
	// 	common.Logger.Error(err)
	// }

}

func dumpToMeilisearch() {
	client := meilisearch.New(config.Cfg.Dest.Meilisearch.Host, meilisearch.WithAPIKey(config.Cfg.Dest.Meilisearch.Key))
	index := client.Index("manyacg")
	artowrks, err := service.DumpArtworks(context.Background())
	if err != nil {
		common.Logger.Fatal(err)
	}
	artworksJSON, err := sonic.Marshal(artowrks)
	if err != nil {
		common.Logger.Fatal(err)
	}
	_, err = index.AddDocuments(artworksJSON)
	if err != nil {
		common.Logger.Fatal(err)
	}
	common.Logger.Infof("Dumped %d artworks to Meilisearch", len(artowrks))
}
