package config

type DumpConfig struct {
	// Will override Databases, tables is in database table_db
	Tables  []string `mapstructure:"tables" json:"tables" yaml:"tables"`
	TableDB string   `mapstructure:"table_db" json:"table_db" yaml:"table_db"`
	// If true, discard error msg, else, output to stderr
	DiscardErr bool `mapstructure:"discard_err" json:"discard_err" yaml:"discard_err"`

	// Set true to skip --master-data if we have no privilege to do
	// 'FLUSH TABLES WITH READ LOCK'
	SkipMasterData bool `mapstructure:"skip_master_data" json:"skip_master_data" yaml:"skip_master_data"`
}

type Canal struct {
	Addr     string     `mapstructure:"addr" json:"addr" yaml:"addr"`
	User     string     `mapstructure:"user" json:"user" yaml:"user"`
	Password string     `mapstructure:"password" json:"password" yaml:"password"`
	Charset  string     `mapstructure:"charset" json:"charset" yaml:"charset"`
	ServerID uint32     `mapstructure:"server_id" json:"server_id" yaml:"server_id"`
	Flavor   string     `mapstructure:"flavor" json:"flavor" yaml:"flavor"`
	Dump     DumpConfig `mapstructure:"dump" json:"dump" yaml:"dump"`
}
