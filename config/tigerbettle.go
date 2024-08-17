package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("tigerbettle", map[string]any{
		"address": config.Env("TB_ADDRESS", "3000"),
	})
}
