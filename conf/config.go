package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ApiKey              string  `json:"api_key" yaml:"api_key" mapstructure:"api_key"`
	SecretKey           string  `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	IsTestnet           bool    `json:"is_testnet" yaml:"is_testnet" mapstructure:"is_testnet"`
	PriceWillTookProfit float64 `json:"price_will_took_profit" yaml:"price_will_took_profit" mapstructure:"price_will_took_profit"`
	Bock                int64   `json:"bock" yaml:"bock" mapstructure:"bock"`
	IntervalCheckTp     int64   `json:"interval_check_tp" yaml:"interval_check_tp" mapstructure:"interval_check_tp"`
	IntervalCheckBuy    int64   `json:"interval_check_buy" yaml:"interval_check_buy" mapstructure:"interval_check_buy"`
	AmountUsdtEachBlock string  `json:"amount_usdt_each_block" yaml:"amount_usdt_each_block" mapstructure:"amount_usdt_each_block"`
	Noti                *Noti   `json:"noti" yaml:"noti" mapstructure:"noti"`
	MysqlDsn            string  `yaml:"mysql_dsn" json:"mysql_dsn" mapstructure:"mysql_dsn"`
}

type Noti struct {
	Url    string `yaml:"url" json:"url" mapstructure:"url"`
	ChatId string `yaml:"chat_id" json:"chat_id" mapstructure:"chat_id"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
