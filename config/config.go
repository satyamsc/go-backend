package config

import "github.com/spf13/viper"

type Config struct {
    DBPath     string
    ServerAddr string
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.SetDefault("DB_PATH", "./data/devices.db")
    viper.SetDefault("SERVER_ADDR", ":8080")
    viper.AutomaticEnv()
    _ = viper.ReadInConfig()
    return &Config{
        DBPath:     viper.GetString("DB_PATH"),
        ServerAddr: viper.GetString("SERVER_ADDR"),
    }, nil
}

