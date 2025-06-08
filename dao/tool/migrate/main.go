package main

import (
	"mcp/dao"
	"mcp/dao/tool"

	"github.com/sirupsen/logrus"
)

func main() {
	db := dao.MustNewDB(tool.DBConfig)
	if err := db.AutoMigrate(
		tool.AllModels...,
	); err != nil {
		logrus.Fatal(err)
	}
}
