package dto

type UnitLogDto struct {
	Content string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type UnitDto struct {
	UnitId int `json:"unit_id"`
	UnitIp string `json:"unit_ip"`
}