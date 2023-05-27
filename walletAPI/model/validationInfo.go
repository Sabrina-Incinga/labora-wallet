package model

type ValidationInfo struct {
	Check struct {
		CheckID      string `json:"check_id"`
		CreationDate string `json:"creation_date"`
		Score        int    `json:"score"`
	} `json:"check"`
}

const (
	StatusRejected  = "RECHAZADO"
	StatusCompleted = "COMPLETADO"
)