package main

import (
	"mcp/dao"
	"mcp/dao/tool"

	"gorm.io/gen"
)

func main() {
	config := tool.DBConfig

	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./dao/query", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(dao.MustNewDB(config))

	g.ApplyBasic(
		tool.AllModels...,
	)

	g.Execute()
}
