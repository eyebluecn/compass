package rest

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nu7hatch/gouuid"
	"time"
)

type SiteDao struct {
	BaseDao
}

//创建站点
func (this *SiteDao) Create(site *Site) *Site {

	if site == nil {
		panic("参数不能为nil")
	}

	timeUUID, _ := uuid.NewV4()
	site.Uuid = string(timeUUID.String())
	site.CreateTime = time.Now()
	site.ModifyTime = time.Now()
	site.Sort = time.Now().UnixNano() / 1e6

	db := this.context.DB.Create(site)
	this.PanicError(db.Error)

	return site
}

//按照Id查询站点，找不到返回nil
func (this *SiteDao) FindByUuid(uuid string) *Site {

	// Read
	var site *Site = &Site{}
	db := this.context.DB.Where(&Site{Base: Base{Uuid: uuid}}).First(site)
	if db.Error != nil {
		return nil
	}
	return site
}

//按照Id查询站点,找不到抛panic
func (this *SiteDao) CheckByUuid(uuid string) *Site {

	// Read
	var site *Site = &Site{}
	db := this.context.DB.Where(&Site{Base: Base{Uuid: uuid}}).First(site)
	this.PanicError(db.Error)
	return site
}

//显示站点列表。
func (this *SiteDao) Page(page int, pageSize int, userUuid string, name string, url string, visible string, sortArray []OrderPair) *Pager {

	var wp = &WherePair{}

	if userUuid != "" {
		wp = wp.And(&WherePair{Query: "user_uuid = ?", Args: []interface{}{userUuid}})
	}

	if name != "" {
		wp = wp.And(&WherePair{Query: "name LIKE ?", Args: []interface{}{"%" + name + "%"}})
	}

	if url != "" {
		wp = wp.And(&WherePair{Query: "url LIKE ?", Args: []interface{}{"%" + url + "%"}})
	}

	if visible != "" {
		tmp := 0
		if (visible == "true") {
			tmp = 1
		} else if (visible == "false") {
			tmp = 0
		} else {
			panic("visible为bool类型，格式错误。")
		}

		wp = wp.And(&WherePair{Query: "visible = ?", Args: []interface{}{tmp}})
	}

	count := 0
	db := this.context.DB.Model(&Site{}).Where(wp.Query, wp.Args...).Count(&count)
	this.PanicError(db.Error)

	var sites []*Site
	orderStr := this.GetSortString(sortArray)
	if orderStr == "" {
		db = this.context.DB.Where(wp.Query, wp.Args...).Offset(page * pageSize).Limit(pageSize).Find(&sites)
	} else {
		db = this.context.DB.Where(wp.Query, wp.Args...).Order(orderStr).Offset(page * pageSize).Limit(pageSize).Find(&sites)
	}

	this.PanicError(db.Error)

	pager := NewPager(page, pageSize, count, sites)

	return pager
}

//保存站点
func (this *SiteDao) Save(site *Site) *Site {

	site.ModifyTime = time.Now()
	db := this.context.DB.Save(site)
	this.PanicError(db.Error)
	return site
}

//删除一个文件，数据库中删除，物理磁盘上删除。
func (this *SiteDao) Delete(site *Site) {

	//删除文件夹本身
	db := this.context.DB.Delete(&site)
	this.PanicError(db.Error)

}
