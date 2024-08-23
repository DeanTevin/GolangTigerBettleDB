package TigerBettleController

import (
	"github.com/goravel/framework/contracts/http"

	helper "goravel/app/helpers"
	tbAccount "goravel/app/http/controllers/TigerBettleController/TigerBettleAccount"
	tbRequests "goravel/app/requests/TigerBettleRequest"
)

type TigerBettleController struct {
	//Dependent services
	validationHelper *helper.RequestValidationHelper
}

func NewTigerBettleController() *TigerBettleController {
	return &TigerBettleController{
		//Inject services
	}
}

func (r *TigerBettleController) PostCreateUserTB(ctx http.Context) http.Response {

	var request tbRequests.CreateUserHistoryRequest
	err := r.validationHelper.TestValidateRequest(&request, ctx)
	if err != nil {
		return ctx.Response().Status(406).Json(err)
	}

	result, errors := tbAccount.TigerBettleAccountAction().CreateUserAccount(request)

	if errors != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": errors.Error(),
		})
	}

	return ctx.Response().Success().Json(result)
}

func (r *TigerBettleController) QueryUserTB(ctx http.Context) http.Response {

	var request tbRequests.QueryFilterUserRequest
	errorval := r.validationHelper.TestValidateRequest(&request, ctx)
	if errorval != nil {
		return ctx.Response().Status(406).Json(errorval)
	}

	result, err := tbAccount.TigerBettleAccountAction().QueryUserAccounts(request)

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	return ctx.Response().Success().Json(result)
}

func (r *TigerBettleController) AccountBalances(ctx http.Context) http.Response {
	var request tbRequests.GetBalanceRequest
	errorval := r.validationHelper.TestValidateRequest(&request, ctx)
	if errorval != nil {
		return ctx.Response().Status(406).Json(errorval)
	}

	result, err := tbAccount.TigerBettleBalanceAction().AccountBalance(request)

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	return ctx.Response().Success().Json(result)
}
