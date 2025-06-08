package tool

import "mcp/dao/models"

var AllModels []any

func init() {
	AllModels = append(AllModels,
		new(models.User),
		new(models.Message),
		new(models.Session),
	)
}
