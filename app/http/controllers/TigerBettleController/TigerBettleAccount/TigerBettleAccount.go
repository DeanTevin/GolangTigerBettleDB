package TigerBettleAccount

import (
	tbService "goravel/app/service/TigerBettle"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBettleAccount struct {
	//Dependent services
}

func TigerBettleAccountAction() *TigerBettleAccount {
	return &TigerBettleAccount{
		//Inject services
	}
}

func (r *TigerBettleAccount) CreateUserHistory(ctx http.Context) ([]map[string]string, error) {
	id, _ := strconv.Atoi(ctx.Request().Input("id"))
	ledger, _ := strconv.Atoi(ctx.Request().Input("ledger"))
	code, _ := strconv.Atoi(ctx.Request().Input("code"))
	UUID, err := tbService.NewTigerBettleService().ConvertUUIDString(ctx.Request().Input("uuid"))

	if err != nil {
		return nil, err
	}

	payloadData := []tbTypes.Account{
		{
			ID:          tbTypes.ToUint128(uint64(id)),
			Ledger:      uint32(ledger),
			Code:        uint16(code),
			UserData128: tbTypes.BytesToUint128(UUID),
			Flags:       uint16(8), //History | This is binary shit idk WTF favours for this shit.
		},
	}

	return tbService.NewTigerBettleService().CreateAccounts(payloadData)
}
