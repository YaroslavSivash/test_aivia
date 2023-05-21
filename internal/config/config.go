package config

import "github.com/joho/godotenv"

type Config struct {
	AppPort string
}

func LoadConfig(confPath string) (*Config, error) {
	var myEnvs map[string]string
	myEnvs, err := godotenv.Read(confPath)
	if err != nil {
		return nil, err
	}
	return &Config{
		AppPort: myEnvs["app_port"],
	}, nil
}
