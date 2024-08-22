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

func CreateUserAccountSanitizeInput(request tbRequests.CreateUserHistoryRequest) ([]tbTypes.Account, error) {
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

func QueryAccountsSanitizeInput(request tbRequests.QueryFilterUserRequest) (tbTypes.QueryFilter, error) {
	uuid, err := tbService.NewTigerBettleService().ConvertUUIDString(request.UUID)

	if err != nil {
		return tbTypes.QueryFilter{}, err
	}

	payloadData := tbTypes.QueryFilter{

		Ledger:      uint32(request.Ledger),
		Code:        uint16(request.Code),
		Flags:       uint32(MapAccountFlags(request.Flags)),
		UserData128: tbTypes.BytesToUint128(uuid),
		Limit:       uint32(100),
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

func (r *TigerBettleAccount) CreateUserAccount(request tbRequests.CreateUserHistoryRequest) ([]map[string]string, error) {

	//sanitize input first
	payloadData, err := CreateUserAccountSanitizeInput(request)
	if err != nil {
		return nil, err
	}

	//connect to client
	client, err := tbService.NewTigerBettleService().GetClient()
	if err != nil {
		return nil, err
	}

	//create account
	_, err = tbService.NewTigerBettleService().CreateAccounts(payloadData, client)

	if err != nil {
		return nil, err
	}

	//lookup created account
	accounts, err := tbService.NewTigerBettleService().LookupAccounts([]tbTypes.Uint128{
		payloadData[0].ID,
	}, client)

	if err != nil {
		return nil, err
	}

	var result []map[string]string

	// mapping result from lookup account
	for _, account := range accounts {
		accountMap := map[string]string{
			"id":     strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.ID.String()), 10),
			"ledger": strconv.FormatUint(uint64(account.Ledger), 10),
			"code":   strconv.FormatUint(uint64(account.Code), 10),
			"uuid":   tbService.NewTigerBettleService().ConvertBytesToUUIDString(account.UserData128.Bytes()),
		}

		// Append the map to the result slice
		result = append(result, accountMap)
	}

	// client close
	client.Close()

	return result, nil
}

func (r *TigerBettleAccount) QueryUserAccounts(request tbRequests.QueryFilterUserRequest) ([]map[string]string, error) {

	client, err := tbService.NewTigerBettleService().GetClient()

	if err != nil {
		return nil, err
	}

	QueryFilter, _ := QueryAccountsSanitizeInput(request)
	QueryFilter.Limit = 100

	accounts, err := tbService.NewTigerBettleService().QueryAccounts(QueryFilter, client)

	if err != nil {
		return nil, err
	}

	var result []map[string]string

	//mapping result from query account
	for _, account := range accounts {
		accountMap := map[string]string{
			"id":     strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.ID.String()), 10),
			"ledger": strconv.FormatUint(uint64(account.Ledger), 10),
			"code":   strconv.FormatUint(uint64(account.Code), 10),
			"uuid":   tbService.NewTigerBettleService().ConvertBytesToUUIDString(account.UserData128.Bytes()),
		}

		// Append the map to the result slice
		result = append(result, accountMap)
	}

	//close client connection
	client.Close()

	return result, nil
}
