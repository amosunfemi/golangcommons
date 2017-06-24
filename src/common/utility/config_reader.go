package utility

import (
	"fmt"

	viper "github.com/spf13/viper"
)

// ConfigReader ...
type ConfigReader struct {
	ConfigType       string
	FileConfigReader FileConfig
}

// FileConfig ... Not instantiated directly
type FileConfig struct {
	FileName     string
	FilePath     string
	ConfigDetail map[string]interface{}
}

// InitConfig ...
func (fileConfig *FileConfig) InitConfig(configfilename string, configfilepath string) error {
	viper.SetConfigName(configfilename)
	viper.AddConfigPath(configfilepath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(err)
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return err
	}
	fileConfig.ConfigDetail = viper.AllSettings()
	return nil
}
