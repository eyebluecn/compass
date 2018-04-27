package rest

import "time"

type TankUploadToken struct {
	TankBase
	UserUuid   string    `json:"userUuid"`
	FolderUuid string    `json:"folderUuid"`
	MatterUuid string    `json:"matterUuid"`
	ExpireTime time.Time `json:"expireTime"`
	Filename   string    `json:"filename"`
	Privacy    bool      `json:"privacy"`
	Size       int64     `json:"size"`
	Ip         string    `json:"ip"`
}
