package config

type GinConfig struct {
	Port         int    `mapstructure:"prot" json:"prot" yaml:"prot"` //接口
	RunMode      string `mapstructure:"run-mode" json:"run-mode" yaml:"run-mode"`
	IsOpenSwager bool   `mapstructure:"is-open-swager" json:"is-open-swager" yaml:"is-open-swager"`
}
