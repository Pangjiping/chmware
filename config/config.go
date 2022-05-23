package config

type Config struct {
}

//func LoadConfig(path string) (Config, error) {
//	var config Config
//	viper.AddConfigPath(path)
//	viper.SetConfigName("app")
//	viper.SetConfigType("env")
//
//	viper.AutomaticEnv()
//
//	if err = viper.ReadInConfig(); err != nil {
//		return err
//	}
//	err = viper.Unmarshal(&config)
//	return err
//}
