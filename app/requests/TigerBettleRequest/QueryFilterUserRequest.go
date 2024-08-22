package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type QueryFilterUserRequest struct {
	Ledger int      `form:"ledger" json:"ledger"`
	Code   int      `form:"code" json:"code"`
	UUID   string   `form:"uuid" json:"uuid"`
	Flags  []string `form:"flags" json:"flags"`
}

func (r *QueryFilterUserRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *QueryFilterUserRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		// The keys are consistent with the incoming keys.
		"ledger": "number",
		"code":   "number",
		"uuid":   "string",
		"flags":  "array",
	}
}

func (r *QueryFilterUserRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *QueryFilterUserRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *QueryFilterUserRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	fields := []string{"id", "ledger", "code", "uuid"}

	for _, field := range fields {
		if value, exists := data.Get(field); exists {
			if floatValue, ok := value.(float64); ok {
				data.Set(field, int(floatValue))
			}
		}
	}

	return nil
}
