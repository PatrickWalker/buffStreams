package helpers

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config is a system wide configuration object
type Config struct {
	DB       DBConfig
	PageSize int32
}

//DBConfig Object represents things needed to open a connection
type DBConfig struct {
	//root:password@tcp(127.0.0.1:3310)/cynalytica?multiStatements=true
	ConnectionString string
}

// GetConfig intializes the config using files and environemnt variables
// The precedece order is. Defaults in this function -> Config File -> Environment Variable --> Docker Secrets
func GetConfig() Config {
	v := viper.New()
	//Set some good defaults
	v.SetDefault("DBUser", "root")
	v.SetDefault("DBPassword", "password")
	v.SetDefault("DBHost", "127.0.0.1:3310")
	v.SetDefault("DBSchmea", "buffup")
	v.SetDefault("DBArgs", "multiStatements=true&parseTime=true")

	v.SetDefault("PageSize", 10)

	//Go and read env variables and override any default config type values
	// for example API_DBUser will work
	v.SetEnvPrefix("api")
	v.AutomaticEnv()

	//read a yaml config from the current folder
	v.SetConfigName("config") // name of config file (without extension)
	v.SetConfigType("json")   // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("Unable to find Config File Locally")
		} else {
			// Config file was found but another error was produced
			log.Printf("Unable to read config file. Continuing with defaults and env variables \n err : %v \n", err.Error())
		}
	}
	conn := buildDBConnectionString(v)
	// Create the DBConfig Connection String
	return Config{
		DB: DBConfig{
			ConnectionString: conn,
		},
		PageSize: v.GetInt32("PageSize"),
	}
}

// buildDBConnectionString returns a mysql connection string
func buildDBConnectionString(v *viper.Viper) string {
	connString := fmt.Sprintf("%v:%v@tcp(%v)/%v?%v", v.Get("DBUser"), v.Get("DBPassword"), v.Get("DBHost"), v.Get("DBSchmea"), v.Get("DBArgs"))
	return connString
}
