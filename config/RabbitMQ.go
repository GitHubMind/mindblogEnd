package config

type RabbitMQ struct {
	User string `"mapstructure:"user" json:"user" yaml:"user"`
	Pass string `"mapstructure:"pass" json:"pass" yaml:"pass"`
	Port string `"mapstructure:"prot" json:"prot" yaml:"prot"`
}
