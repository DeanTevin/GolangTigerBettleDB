package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("tigerbettle", map[string]any{
		"address":        config.Env("TB_ADDRESS", "3000"),
		"account_flags":  config.Env("TB_ACCOUNT_FLAGS", []string{"History", "CreditsMustNotExceedDebits", "DebitsMustNotExceedCredits", "Linked"}),
		"utc_timeformat": config.Env("TB_UTC_TIMEFORMAT", 7),
	})
}
