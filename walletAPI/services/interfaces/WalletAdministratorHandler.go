package interfaces

import (
	"github.com/labora-wallet/walletAPI/model"
)

type WalletAdministratorHandler interface {
	AttemptWalletCreation(wallet model.WalletDTO) (string, int64, error)
	ValidateScore(nationalIdentityNumber, countryId string) string
}
