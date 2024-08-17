package TigerBettleController

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"

	tbService "goravel/app/service/TigerBettle"

	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBettleController struct {
	//Dependent services
}

func NewTigerBettleController() *TigerBettleController {
	return &TigerBettleController{
		//Inject services
	}
}

func (r *TigerBettleController) PostCreateUserTB(ctx http.Context) http.Response {
	client, err := tbService.NewTigerBettleService().GetClient()

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	id, _ := strconv.Atoi(ctx.Request().Input("id"))
	ledger, _ := strconv.Atoi(ctx.Request().Input("ledger"))
	code, _ := strconv.Atoi(ctx.Request().Input("code"))
	UUID, err := tbService.NewTigerBettleService().ConvertUUIDString(ctx.Request().Input("uuid"))

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	res, err := client.CreateAccounts([]tbTypes.Account{
		{
			ID:          tbTypes.ToUint128(uint64(id)),
			Ledger:      uint32(ledger),
			Code:        uint16(code),
			UserData128: tbTypes.BytesToUint128(UUID),
			Flags:       uint16(8), //History | This is binary shit idk WTF favours for this shit.
		},
	})

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"Error": err.Error(),
		})
	}

	for _, err := range res {
		return ctx.Response().Status(500).Json(http.Json{
			"data": err.Result.String(),
		})
	}

	accounts, err := client.LookupAccounts([]tbTypes.Uint128{
		tbTypes.ToUint128(uint64(id)),
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
			"uuid":   account.UserData128.String(),
		}

		// Append the map to the result slice
		result = append(result, accountMap)
	}

	client.Close() //close connection

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
