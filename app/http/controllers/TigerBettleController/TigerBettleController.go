package TigerBettleController

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	helper "goravel/app/helpers"
	tbAccount "goravel/app/http/controllers/TigerBettleController/TigerBettleAccount"
	tbRequests "goravel/app/requests/TigerBettleRequest"
	tbService "goravel/app/service/TigerBettle"

	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
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

	result, errors := tbAccount.TigerBettleAccountAction().CreateUserHistory(request)

	if errors != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": errors.Error(),
		})
	}

	return ctx.Response().Success().Json(result)
}

func (r *TigerBettleController) QueryUserTB(ctx http.Context) http.Response {

	client, err := tbService.NewTigerBettleService().GetClient()

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	// code, _ := strconv.Atoi(ctx.Request().Input("code"))
	uuid, _ := tbService.NewTigerBettleService().ConvertUUIDString(ctx.Request().Input("code"))

	accounts, err := client.QueryAccounts(tbTypes.QueryFilter{
		UserData128: tbTypes.BytesToUint128(uuid),
		Limit:       uint32(1),
	})

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	var result []map[string]string

	for _, account := range accounts {
		accountMap := map[string]string{
			"id":     account.ID.String(),
			"ledger": strconv.FormatUint(uint64(account.Ledger), 10),
			"code":   strconv.FormatUint(uint64(account.Code), 10),
			"uuid":   tbService.NewTigerBettleService().ConvertBytesToUUIDString(account.UserData128.Bytes()),
		}

		// Append the map to the result slice
		result = append(result, accountMap)
	}

	client.Close() //close connection

	return ctx.Response().Success().Json(result)

}
