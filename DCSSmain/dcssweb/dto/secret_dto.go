package dto

import "github.com/Alan-Lxc/crypto_contest/dcssweb/model"

type SecretDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

type GetSecretDto struct {
	Name            string `json:"name"`
	Degree          int64  `json:"degree"`
	Counter         int64  `json:"counter"`
	UserId          uint   `json:"user_id"`
	CreateTime      string `json:"create_time"`
	LastUpdateTime  string `json:"last_update_time"`
	LastHandoffTime string `json:"last_handoff_time"`
	Description     string `json:"description"`
}

func ToGetSecretDto(secret model.Secret) GetSecretDto {
	var timeLayoutStr = "2006-01-02 15:04:05" //go中的时间格式化必须是这个时间
	return GetSecretDto{
		Name:            secret.Secretname,
		Degree:          secret.Degree,
		Counter:         secret.Counter,
		UserId:          secret.UserId,
		CreateTime:      secret.CreatedAt.Format(timeLayoutStr),
		LastUpdateTime:  secret.UpdatedAt.Format(timeLayoutStr),
		LastHandoffTime: secret.LastHandoffAt.Format(timeLayoutStr),
		Description:     secret.Description,
	}
}

type GetSecretListDto struct {
	secrets []model.Secret `json:"secrets"`
}

func ToRetrieveSecretByUseridDto(secrets []model.Secret) GetSecretListDto {
	return GetSecretListDto{secrets: secrets}
}
