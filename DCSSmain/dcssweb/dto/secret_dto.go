package dto

import "github.com/Alan-Lxc/crypto_contest/dcssweb/model"

type SecretDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

type RetrieveSecretByIdDto struct {
	Degree      int64  `json:"degree"`
	Counter     int64  `json:"counter"`
	UserId      uint   `json:"user_id"`
	Description string `json:"description"`
}

func ToRetrieveSecretByIdDto(secret model.Secret) RetrieveSecretByIdDto {
	return RetrieveSecretByIdDto{
		Degree:      secret.Degree,
		Counter:     secret.Counter,
		UserId:      secret.UserId,
		Description: secret.Description,
	}
}

type RetrieveSecretByUseridDto struct {
	secrets []model.Secret `json:"secrets"`
}

func ToRetrieveSecretByUseridDto(secrets []model.Secret) RetrieveSecretByUseridDto {
	return RetrieveSecretByUseridDto{secrets: secrets}
}
