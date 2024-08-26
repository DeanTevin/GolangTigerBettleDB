package TigerBettleAccount

import (
	tbRequests "goravel/app/requests/TigerBettleRequest"
	tbService "goravel/app/service/TigerBettle"
	"strconv"

	"github.com/goravel/framework/facades"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBettleTransfer struct {
	//Dependent services
}

func TigerBettleTransferAction() *TigerBettleTransfer {
	return &TigerBettleTransfer{
		//Inject services
	}
}

func QueryTransferSanitizeInput(request tbRequests.QueryTransferRequest) (tbTypes.QueryFilter, error) {
	var uuid [16]byte
	var err error
	if request.UUID != "" {
		uuid, err = tbService.NewTigerBettleService().ConvertUUIDString(request.UUID)

		if err != nil {
			return tbTypes.QueryFilter{}, err
		}
	}

	payloadData := tbTypes.QueryFilter{

		Ledger: uint32(request.Ledger),
		Code:   uint16(request.Code),
		Flags: tbTypes.QueryFilterFlags{
			Reversed: request.Reversed,
		}.ToUint32(),
		UserData128:  tbTypes.BytesToUint128(uuid),
		Limit:        uint32(100),
		TimestampMin: tbService.NewTigerBettleService().ConvertTimestampString(request.TimestampMin),
		TimestampMax: tbService.NewTigerBettleService().ConvertTimestampString(request.TimestampMax),
	}

	facades.Log().Debug(payloadData)

	return payloadData, nil
}

func (r *TigerBettleTransfer) QueryTransfer(request tbRequests.QueryTransferRequest) ([]map[string]string, error) {

	client, err := tbService.NewTigerBettleService().GetClient()

	if err != nil {
		return nil, err
	}

	QueryFilter, _ := QueryTransferSanitizeInput(request)

	transfers, err := tbService.NewTigerBettleService().QueryTransfer(QueryFilter, client)

	facades.Log().Debug(transfers)
	if err != nil {
		return nil, err
	}

	facades.Log().Debug(transfers)

	var result []map[string]string

	//mapping result from query account
	for _, transfer := range transfers {

		transactionMap := map[string]string{
			"transfer_id":        strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(transfer.ID.String()), 10),
			"credits_account_id": strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(transfer.CreditAccountID.String()), 10),
			"debits_account_id":  strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(transfer.DebitAccountID.String()), 10),
			"user_data_128":      strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(transfer.UserData128.String()), 10),
			"user_data_64":       strconv.FormatUint(transfer.UserData64, 10),
			"user_data_32":       strconv.FormatUint(uint64(transfer.UserData32), 10),
			"timeout":            strconv.FormatUint(uint64(transfer.Timeout), 10),
			"ledger":             strconv.FormatUint(uint64(transfer.Ledger), 10),
			"code":               strconv.FormatUint(uint64(transfer.Code), 10),
			"flags":              strconv.FormatUint(uint64(transfer.Flags), 10),
			"amount":             strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(transfer.Amount.String()), 10),
			"timestamp":          tbService.NewTigerBettleService().ConvertTimestampToString(transfer.Timestamp),
		}

		// Append the map to the result slice
		result = append(result, transactionMap)
	}

	//close client connection
	client.Close()

	return result, nil
}
