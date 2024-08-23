package TigerBettleAccount

import (
	tbRequests "goravel/app/requests/TigerBettleRequest"
	tbService "goravel/app/service/TigerBettle"
	"strconv"

	"github.com/goravel/framework/facades"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBettleBalance struct {
	//Dependent services
}

func TigerBettleBalanceAction() *TigerBettleAccount {
	return &TigerBettleAccount{
		//Inject services
	}
}

func AccountBalanceSanitizeInput(request tbRequests.GetBalanceRequest) (tbTypes.AccountFilter, error) {

	id, err := tbTypes.HexStringToUint128(strconv.Itoa(request.ID))
	if err != nil {
		return tbTypes.AccountFilter{}, err
	}

	payloadData := tbTypes.AccountFilter{

		AccountID: id,
		Limit:     100,
		Flags: tbTypes.AccountFilterFlags{
			Debits:   true, // Include transfer from the debit side.
			Credits:  true, // Include transfer from the credit side.
			Reversed: true, // Sort by timestamp in reverse-chronological order.
		}.ToUint32(),
	}

	return payloadData, nil
}

func (r *TigerBettleAccount) AccountBalance(request tbRequests.GetBalanceRequest) ([]map[string]string, error) {

	client, err := tbService.NewTigerBettleService().GetClient()

	if err != nil {
		return nil, err
	}

	AccountFilter, _ := AccountBalanceSanitizeInput(request)

	accounts, err := tbService.NewTigerBettleService().AccountBalances(AccountFilter, client)

	if err != nil {
		return nil, err
	}

	facades.Log().Debug(accounts)

	var result []map[string]string

	//mapping result from query account
	for _, account := range accounts {
		accountMap := map[string]string{
			"credits_pending": strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.CreditsPending.String()), 10),
			"credits_posted":  strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.CreditsPosted.String()), 10),
			"debits_pending":  strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.DebitsPending.String()), 10),
			"debits_posted":   strconv.FormatUint(tbService.NewTigerBettleService().HexStringToUint(account.DebitsPosted.String()), 10),
		}

		// Append the map to the result slice
		result = append(result, accountMap)
	}

	//close client connection
	client.Close()

	return result, nil
}
