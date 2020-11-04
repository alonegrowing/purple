package sql

type SQLGroupConfig struct {
	Name   string   `toml:"name"`
	Master string   `toml:"master"`
	Slaves []string `toml:"slaves"`
}
