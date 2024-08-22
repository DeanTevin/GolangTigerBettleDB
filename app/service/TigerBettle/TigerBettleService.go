package TigerBettle

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbTypes "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TigerBettleService struct {
	//Dependent services

	TB_ADDRESS    string
	TB_CONNECTION bool
}

func NewTigerBettleService() *TigerBettleService {
	return &TigerBettleService{
		//Inject services
		TB_ADDRESS:    facades.Config().GetString("tigerbettle.address"),
		TB_CONNECTION: checkConnection(facades.Config().GetString("tigerbettle.address")),
	}
}

func checkConnection(address string) bool {
	// Attempt to establish a connection to the address
	conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
	if err != nil {
		// Connection failed
		return false
	}
	// Close the connection once it's established
	conn.Close()
	return true
}

func (r *TigerBettleService) GetClient() (tb.Client, error) {

	if !r.TB_CONNECTION {
		return nil, errors.New("error connection closed")
	}

	client, err := tb.NewClient(tbTypes.ToUint128(0), []string{r.TB_ADDRESS})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *TigerBettleService) ConvertUUIDString(uuidStr string) ([16]byte, error) {
	// Remove hyphens from UUID string
	cleanedUUIDStr := strings.ReplaceAll(uuidStr, "-", "")

	// Convert hex string to byte array
	byteArray, err := hex.DecodeString(cleanedUUIDStr)
	if err != nil {
		return [16]byte{}, errors.New("error decoding hex string")
	}

	// Ensure the byte array has 16 bytes
	if len(byteArray) != 16 {
		return [16]byte{}, errors.New("decoded byte array does not have 16 bytes")
	}

	// Convert byte array to [16]byte
	var uuid [16]byte
	copy(uuid[:], byteArray)

	return uuid, nil
}

func (r *TigerBettleService) ConvertBytesToUUIDString(uuid [16]byte) string {
	// Convert [16]byte to a hexadecimal string
	hexStr := hex.EncodeToString(uuid[:])

	// Format the string to include hyphens (UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	formattedUUID := fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8],   // 8 characters
		hexStr[8:12],  // 4 characters
		hexStr[12:16], // 4 characters
		hexStr[16:20], // 4 characters
		hexStr[20:])   // 12 characters

	return formattedUUID
}

func (r *TigerBettleService) HexStringToUint(hexStr string) uint64 {
	// Parse the hexadecimal string to a uint64
	uint64Value, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0
	}

	// Convert uint64 to uint (may be uint32 or uint64 depending on architecture)
	return uint64(uint64Value)
}

func (r *TigerBettleService) CreateAccounts(payloadData []tbTypes.Account, client tb.Client) ([]tbTypes.AccountEventResult, error) {

	res, err := client.CreateAccounts(payloadData)

	if err != nil {
		return nil, err
	}

	for _, err := range res {
		return nil, errors.New(err.Result.String())
	}

	return res, nil
}

func (r *TigerBettleService) LookupAccounts(payloadData []tbTypes.Uint128, client tb.Client) ([]tbTypes.Account, error) {

	//lookup accounts by ID
	accounts, err := client.LookupAccounts(payloadData)

	if err != nil {
		return []tbTypes.Account{}, err
	}

	return accounts, nil
}

func (r *TigerBettleService) QueryAccounts(payloadData tbTypes.QueryFilter, client tb.Client) ([]tbTypes.Account, error) {
	accounts, err := client.QueryAccounts(payloadData)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *TigerBettleService) AccountBalances(payloadData tbTypes.AccountFilter, client tb.Client) ([]tbTypes.AccountBalance, error) {
	balance, err := client.GetAccountBalances(payloadData)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
