package rest

type Site struct {
	Base
	UserUuid        string `json:"userUuid"`
	Name            string `json:"name"`
	FaviconTankUuid string `json:"faviconTankUuid"`
	FaviconUrl      string `json:"faviconUrl"`
	Url             string `json:"url"`
	Hit             int64  `json:"hit"`
	FaviconTank     *Tank  `gorm:"-" json:"faviconTank"`
}

// set Site's table name
func (Site) TableName() string {
	return TABLE_PREFIX + "site"
}
