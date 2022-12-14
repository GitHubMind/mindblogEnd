package config

type Server struct {
	MysqlList []GeneralDB `mapstructure:"mysqllist" json:"mysqllist" yaml:"mysqllist"`
	GinConfig GinConfig   `mapstructure:"gin" json:"gin" yaml:"gin"`
	System    System      `mapstructure:"system" json:"system" yaml:"system"`
	Mysql     Mysql       `mapstructure:"mysql" json:"mysql" yaml:"mysqldefault"`
	JWT       JWT         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Captcha   Captcha     `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Redis     Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	//可白名单配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	Local Local `mapstructure:"local" json:"local" yaml:"local"`

	RabbitMQ RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`

	//AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
}
