package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Log      logConfig      `toml:"log" mapstructure:"log" json:"log" yaml:"log"`
	Database databaseConfig `toml:"database" mapstructure:"database" json:"database" yaml:"database"`
	Dest     destConfig     `toml:"dest" mapstructure:"dest" json:"dest" yaml:"dest"`
}

type destConfig struct {
	Type        string            `toml:"type" mapstructure:"type" json:"type" yaml:"type"`
	Meilisearch meilisearchConfig `toml:"meilisearch" mapstructure:"meilisearch" json:"meilisearch" yaml:"meilisearch"`
}

type meilisearchConfig struct {
	Host     string                    `toml:"host" mapstructure:"host" json:"host" yaml:"host"`
	Key      string                    `toml:"key" mapstructure:"key" json:"key" yaml:"key"`
	Index    string                    `toml:"index" mapstructure:"index" json:"index" yaml:"index"`
	Embedder meilisearchEmbedderConfig `toml:"embedder" mapstructure:"embedder" json:"embedder" yaml:"embedder"`
}

type meilisearchEmbedderConfig struct {
	Name                     string `toml:"name" mapstructure:"name" json:"name" yaml:"name"`
	Source                   string `toml:"source" mapstructure:"source" json:"source" yaml:"source"`
	Model                    string `toml:"model" mapstructure:"model" json:"model" yaml:"model"`
	APIKey                   string `toml:"api_key" mapstructure:"api_key" json:"api_key" yaml:"api_key"`
	DocumentTemplate         string `toml:"document_template" mapstructure:"document_template" json:"document_template" yaml:"document_template"`
	URL                      string `toml:"url" mapstructure:"url" json:"url" yaml:"url"`
	Dimensions              int    `toml:"dimensions" mapstructure:"dimensions" json:"dimensions" yaml:"dimensions"`
}

var Cfg *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/github.com/ManyACG/lolidump/")
	viper.SetConfigType("toml")
	viper.SetEnvPrefix("github.com/ManyACG/lolidump")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("log.level", "TRACE")
	viper.SetDefault("log.file_path", "logs/ManyACG.log")
	viper.SetDefault("log.backup_num", 7)

	viper.SetDefault("database.database", "manyacg")
	viper.SetDefault("database.max_staleness", 120)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error when reading config: %s\n", err)
		os.Exit(1)
	}
	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		fmt.Printf("error when unmarshal config: %s\n", err)
		os.Exit(1)
	}
}
