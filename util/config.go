package util

import "github.com/spf13/viper"

// Stores all configs for the applications.
type Config struct {
	// Mapstructure is an Unmarshalling commanded used to map env vars to the struct vars.
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadCOnfig will read configs from a file or environment variables.
//  It takes in a path as a type string and outputs a config object or an error of type error.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	//Matches the name of the app in app.env
	viper.SetConfigName("app")
	// Matches app.env
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
