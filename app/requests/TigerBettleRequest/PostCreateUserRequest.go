package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CreateUserHistoryRequest struct {
	ID     int      `form:"id" json:"id"`
	Ledger int      `form:"ledger" json:"ledger"`
	Code   int      `form:"code" json:"code"`
	UUID   string   `form:"uuid" json:"uuid"`
	Flags  []string `form:"flags" json:"flags"`
}

func (r *CreateUserHistoryRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateUserHistoryRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		// The keys are consistent with the incoming keys.
		"id":     "required|max_len:255|number",
		"ledger": "required|number",
		"code":   "required|number",
		"uuid":   "required|string",
		"flags":  "array",
	}
}

func (r *CreateUserHistoryRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateUserHistoryRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateUserHistoryRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
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
