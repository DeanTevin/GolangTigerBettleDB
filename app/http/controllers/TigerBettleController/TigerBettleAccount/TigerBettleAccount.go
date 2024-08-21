package TigerBettleAccount

import (
	tbRequests "goravel/app/requests/TigerBettleRequest"
	tbService "goravel/app/service/TigerBettle"
	"strconv"

	"github.com/goravel/framework/facades"
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

func SanitizeInput(request tbRequests.CreateUserHistoryRequest) ([]tbTypes.Account, error) {
	uuid, err := tbService.NewTigerBettleService().ConvertUUIDString(request.UUID)

	if err != nil {
		return nil, err
	}

	payloadData := []tbTypes.Account{
		{
			ID:          tbTypes.ToUint128(uint64(request.ID)),
			Ledger:      uint32(request.Ledger),
			Code:        uint16(request.Code),
			Flags:       MapAccountFlags(request.Flags),
			UserData128: tbTypes.BytesToUint128(uuid),
		},
	}

	return payloadData, nil
}

func MapAccountFlags(requestFlags []string) uint16 {

	values := facades.Config().Get("tigerbettle.account_flags")
	validation, _ := values.([]string) //assert as array slice

	// Create a binary representation based on the presence of each value
	binaryString := ""
	for _, checkValue := range validation {
		if contains(requestFlags, checkValue) {
			binaryString += "1"
		} else {
			binaryString += "0"
		}
	}

	result, _ := strconv.ParseInt(binaryString, 2, 64)
	return uint16(result)
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func (r *TigerBettleAccount) CreateUserHistory(request tbRequests.CreateUserHistoryRequest) ([]map[string]string, error) {

	payloadData, err := SanitizeInput(request)

	if err != nil {
		return nil, err
	}

	return tbService.NewTigerBettleService().CreateAccounts(payloadData)
}
