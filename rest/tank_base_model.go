package rest

import "time"

type TankBase struct {
	Uuid       string    `json:"uuid"`
	Sort       int64     `json:"sort"`
	ModifyTime time.Time `json:"modifyTime"`
	CreateTime time.Time `json:"createTime"`
}
