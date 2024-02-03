package systems

import "github.com/spf13/viper"

type AppConfig struct {
	Port    string
	Address string

	MangoDBConfig
}

type MangoDBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func GetAndSetupConfig() (*AppConfig, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return getConfig()
}

func getConfig() (*AppConfig, error) {
	config := AppConfig{
		Address: viper.GetString("address"),
		Port:    viper.GetString("port"),
		MangoDBConfig: MangoDBConfig{
			Host: viper.GetString("mangodb.host"),
			Port: viper.GetString("mangodb.port"),
			User: viper.GetString("mangodb.user"),
			Pass: viper.GetString("mangodb.pass"),
			Name: viper.GetString("mangodb.name"),
		},
	}
	return &config, nil
}
