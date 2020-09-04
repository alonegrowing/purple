package resource

import (
	"purple/pkg/config"
	"purple/stone/sql"
)

var DefaultDB *sql.Group

func init() {
	NewMysqlGroup(config.ServiceConfig.Database)

	DefaultDB = sql.SQLGroupManager.Get("intersting")
}

func NewMysqlGroup(database []sql.SQLGroupConfig) error {
	if len(database) == 0 {
		return nil
	}
	for _, d := range database {
		g, err := sql.NewGroup(d.Name, d.Master, d.Slaves)
		if err != nil {
			return err
		}
		err = sql.SQLGroupManager.Add(d.Name, g)
		if err != nil {
			return err
		}
	}
	return nil
}
