package rest

import "time"

type TankDownloadToken struct {
	TankBase
	UserUuid   string    `json:"userUuid"`
	MatterUuid string    `json:"matterUuid"`
	ExpireTime time.Time `json:"expireTime"`
	Ip         string    `json:"ip"`
}
