package viperLib

import (
	"blog/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitViper() (config.Server, *viper.Viper) {
	//viper.SetConfigFile("./config/config.yaml") // 指定配置文件路径
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")

	viper.AddConfigPath("config")
	viper.AddConfigPath("config.yaml")
	err := viper.ReadInConfig()
	// 读取配置信息
	if err != nil { // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	//如果发生了变化，我可以重新加载项目吗？ 12.14

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)

	})
	var server config.Server
	if err = viper.Unmarshal(&server); err != nil {
		fmt.Println("viper bug", err)
	}
	//初始化存储，并且设定过期时间

	return server, viper.GetViper()
}
