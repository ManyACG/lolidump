package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/ManyACG/lolidump/common"
	"github.com/ManyACG/lolidump/config"
	"github.com/ManyACG/lolidump/dao"
	"github.com/ManyACG/lolidump/service"
	"github.com/bytedance/sonic"
	"github.com/meilisearch/meilisearch-go"

	_ "net/http/pprof"

	"gopkg.in/osteele/liquid.v1"
)

func Run() {
	config.InitConfig()
	common.Init()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	common.Logger.Info("Start lolidump...")
	dao.InitDB(ctx)
	defer func() {
		if err := dao.Client.Disconnect(ctx); err != nil {
			common.Logger.Panic(err)
			os.Exit(1)
		}
	}()

	switch config.Cfg.Dest.Type {
	case "meilisearch":
		dumpToMeilisearch(ctx)
	default:
		common.Logger.Panic("Unsupported destination type")
	}
}

func testTemplate() {
	cfg := config.Cfg.Dest.Meilisearch
	docTemplateStr := cfg.Embedder.DocumentTemplate
	engine := liquid.NewEngine()
	tmpl, err := engine.ParseString(docTemplateStr)
	if err != nil {
		common.Logger.Panic(err)
	}

	testData := map[string]any{
		"ID":          "1234567890",
		"title":       "Test Title",
		"artist":      "Test Artist",
		"tags":        "tag1,tag2,tag3",
		"description": "",
		"r18":         true,
	}
	result, err := tmpl.Render(testData)
	if err != nil {
		common.Logger.Panic(err)
	}
	common.Logger.Infof("Template test result: %s", result)
}

func dumpToMeilisearch(ctx context.Context) {
	testTemplate()
	cfg := config.Cfg.Dest.Meilisearch

	client := meilisearch.New(config.Cfg.Dest.Meilisearch.Host, meilisearch.WithAPIKey(cfg.Key))
	_, err := client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
		Uid:        cfg.Index,
		PrimaryKey: "id",
	})
	if err != nil {
		common.Logger.Panic(err)
	}

	index := client.Index(cfg.Index)
	_, err = index.UpdateEmbeddersWithContext(ctx, map[string]meilisearch.Embedder{
		cfg.Embedder.Name: {
			Source:           cfg.Embedder.Source,
			APIKey:           cfg.Embedder.APIKey,
			Model:            cfg.Embedder.Model,
			URL:              cfg.Embedder.URL,
			DocumentTemplate: cfg.Embedder.DocumentTemplate,
			Dimensions:       cfg.Embedder.Dimensions,
		},
	})
	if err != nil {
		common.Logger.Panic(err)
	}

	artowrks, err := service.DumpArtworkToSearch(ctx)
	if err != nil {
		common.Logger.Panic(err)
	}
	artworksJSON, err := sonic.Marshal(artowrks)
	if err != nil {
		common.Logger.Panic(err)
	}
	_, err = index.AddDocuments(artworksJSON)
	if err != nil {
		common.Logger.Panic(err)
	}
	common.Logger.Infof("Dumped %d artworks to Meilisearch", len(artowrks))
}
