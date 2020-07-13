/**
* @Author: Dylan
* @Date: 2020/6/9 10:20
 */

package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var MediaConf = &Config{}

var ReverseConfig *viper.Viper

var once sync.Once

const _configPath = "src/helpers/config"

func InitConfig() {
	once.Do(func() {
		ReverseConfig = viper.New()
		ReverseConfig.SetConfigName("config")
		ReverseConfig.SetConfigType("yaml")
		ReverseConfig.AddConfigPath(_configPath)
		if err := ReverseConfig.ReadInConfig(); err != nil {
			log.Fatal(err)
		}

		ReverseConfig.Unmarshal(MediaConf)

	})
}
