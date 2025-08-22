package main

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Id          string `mapstructure:"id"`
	IsLeader    bool   `mapstructure:"is_leader"`
	BaseDataDir string `mapstructure:"base_data_dir"`
	Address     string `mapstructure:"address"`
	LeaderAddr  string `mapstructure:"leader_addr"`
	Advertise   string `mapstructure:"advertise"`
	ServerAddr  string `mapstructure:"server_addr"`
	Peers       []Peer `mapstructure:"-"`
}

func NewConfig() *Config {
	viper.SetDefault("id", "1")
	viper.SetDefault("is_leader", true)
	viper.SetDefault("base_data_dir", "C:\\Drafts\\thorac\\data")
	viper.SetDefault("address", "0.0.0.0:8080")
	viper.SetDefault("advertise", "0.0.0.0:8080")
	viper.SetDefault("server_addr", "0.0.0.0:8080")
	viper.SetDefault("leader_addr", "0.0.0.0:8080")
	viper.SetDefault("peers", []Peer{})

	viper.AutomaticEnv()

	viper.SetEnvPrefix("THORAC")
	viper.AllowEmptyEnv(true)

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if peerStr := viper.GetString("peers"); peerStr != "" {
		var peers []Peer
		if err := json.Unmarshal([]byte(peerStr), &peers); err != nil {
			log.Fatalf("Invalid JSON in THORAC_PEERS: %v", err)
		}
		config.Peers = peers
	}

	return config
}
