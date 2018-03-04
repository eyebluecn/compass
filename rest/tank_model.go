package rest

type Tank struct {
	Base
	UserUuid   string `json:"userUuid"`
	Name       string `json:"name"`
	MatterUuid string `json:"matterUuid"`
	Size       int64  `json:"size"`
	Privacy    bool   `json:"privacy"`
	Url        string `json:"url"`
	Remark     string `json:"remark"`
	Confirmed  bool   `json:"confirmed"`
	//用于上传的uploadToken.
	UploadTokenUuid int `gorm:"-" json:"uploadTokenUuid"`
	//客户端需要将文件上传到何处去。
	UploadUrl int `gorm:"-" json:"uploadUrl"`
}

// set File's table name to be `profiles`
func (Tank) TableName() string {
	return TABLE_PREFIX + "tank"
}
