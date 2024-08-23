package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetBalanceRequest struct {
	ID int `form:"account_id" json:"account_id"`
}

func (r *GetBalanceRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetBalanceRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		// The keys are consistent with the incoming keys.
		"account_id": "required|number",
	}
}

func (r *GetBalanceRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetBalanceRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetBalanceRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	fields := []string{"account_id"}

	for _, field := range fields {
		if value, exists := data.Get(field); exists {
			if floatValue, ok := value.(float64); ok {
				data.Set(field, int(floatValue))
			}
		}
	}

	return nil
}
