package helpers

import (
	"github.com/goravel/framework/contracts/http"
)

type RequestValidationHelper struct {
	//Dependent services
}

func (r *RequestValidationHelper) TestValidateRequest(request http.FormRequest, ctx http.Context) http.Json {
	errors, err := ctx.Request().ValidateRequest(request)

	if err != nil {
		return http.Json{
			"error": err.Error(),
		}
	}

	if errors != nil {
		return http.Json{
			"error": errors.All(),
		}
	}

	return nil
}
