package config

import (
	"encoding/json"
	"os"

	"github.com/alexanderlazar/my-rss-aggregator/internal/database"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	URL  string `json:"db_url"`
	User string `json:"current_user_name"`
}

type State struct {
	Db        *database.Queries
	Configptr *Config
}

func Read() (Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(homedir + "/" + configFileName)
	if err != nil {
		return Config{}, err
	}
	var conf Config
	json.Unmarshal(data, &conf)
	return conf, nil
}

func (c Config) SetUser(user string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.User = user
	updated, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(homedir+"/"+configFileName, updated, 0644)
	if err != nil {
		return err
	}
	return nil
}
