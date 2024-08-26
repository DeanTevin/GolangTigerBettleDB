package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type QueryTransferRequest struct {
	Ledger       int      `form:"ledger" json:"ledger"`
	Code         int      `form:"code" json:"code"`
	UUID         string   `form:"uuid" json:"uuid"`
	Flags        []string `form:"flags" json:"flags"`
	TimestampMin string   `form:"timestamp_min" json:"timestamp_min"`
	TimestampMax string   `form:"timestamp_max" json:"timestamp_max"`
	Reversed     bool     `form:"reversed" json:"reversed"`
}

func (r *QueryTransferRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *QueryTransferRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		// The keys are consistent with the incoming keys.
		"ledger":        "number",
		"code":          "number",
		"uuid":          "string",
		"flags":         "array",
		"timestamp_min": "date|lte_field:timestamp_max",
		"timestamp_max": "date|gte_field:timestamp_min",
		"reversed":      "bool",
	}
}

func (r *QueryTransferRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *QueryTransferRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *QueryTransferRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
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
