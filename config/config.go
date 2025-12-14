package config

type Config struct {
    DBPath     string
    ServerAddr string
}

func Load() (*Config, error) {
    return &Config{
        DBPath:     "data/app.db",
        ServerAddr: ":8080",
    }, nil
}
