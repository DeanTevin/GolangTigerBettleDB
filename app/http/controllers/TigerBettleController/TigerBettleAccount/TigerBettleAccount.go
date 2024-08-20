package TigerBettleAccount

import (
	tbRequests "goravel/app/requests/TigerBettleRequest"
	tbService "goravel/app/service/TigerBettle"

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

func (r *TigerBettleAccount) SanitizeInput(request tbRequests.CreateUserHistoryRequest) (ID int, Ledger int, Code int, UUID [16]byte, errs error) {
	uuid, err := tbService.NewTigerBettleService().ConvertUUIDString(request.UUID)

	if err != nil {
		return 0, 0, 0, [16]byte{}, err
	}

	return request.ID, request.Ledger, request.Code, uuid, nil
}

func (r *TigerBettleAccount) CreateUserHistory(request tbRequests.CreateUserHistoryRequest) ([]map[string]string, error) {

	id, ledger, code, UUID, err := r.SanitizeInput(request)

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
