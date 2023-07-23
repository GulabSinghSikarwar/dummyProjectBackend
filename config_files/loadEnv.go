package ConfigFiles

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	JwtExpiresIn time.Duration `mapstructure:"JWT_EXPIRED_IN"`
	JwtMaxage    int           `mapstructure:"JWT_MAXAGE"`
	JwtSecret    string        `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile("app")
	viper.SetConfigType("env")

	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return

}
