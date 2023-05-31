package config

type Config struct {
	Zap   Zap   `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`
	Bot   Bot   `mapstructure:"bot" json:"bot" yaml:"bot"`
	Canal Canal `mapstructure:"canal" json:"canal" yaml:"canal"`
}
