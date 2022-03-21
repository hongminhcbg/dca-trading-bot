package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ApiKey    string `json:"api_key" yaml:"api_key" mapstructure:"api_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	IsTestnet bool   `json:"is_testnet" yaml:"is_testnet" mapstructure:"is_testnet"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	// debug, fucking bullshit viper
	keys := viper.AllKeys()
	fmt.Println("all key viper read = ", keys)
	for _, k := range keys {
		fmt.Printf("[DB]    k = %s, value = %v\n", k, viper.Get(k))
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
