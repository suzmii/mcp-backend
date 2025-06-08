package tool

import "mcp/dao"

var DBConfig = dao.DBConfig{
	Host:         "localhost",
	Port:         3306,
	Username:     "root",
	Password:     `oppofindx2`,
	DatabaseName: "mcp",
	DriverName:   "mariadb",
	AutoCreateDB: true,
}
