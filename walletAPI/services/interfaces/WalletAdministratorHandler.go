package interfaces

import (
	"github.com/labora-wallet/walletAPI/model/dtos"
)

type WalletAdministratorHandler interface {
	AttemptWalletCreation(wallet dtos.WalletDTO) (string, int64, error)
	AttemptWalletRemoval(walletId int64) (int64, error)
	ValidateScore(nationalIdentityNumber, countryId string) string
}
