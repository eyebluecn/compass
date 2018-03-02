package rest

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nu7hatch/gouuid"
	"time"
)

type TankDao struct {
	BaseDao
}

//创建
func (this *TankDao) Create(tank *Tank) *Tank {

	timeUUID, _ := uuid.NewV4()
	tank.Uuid = string(timeUUID.String())
	tank.CreateTime = time.Now()
	tank.ModifyTime = time.Now()
	db := this.context.DB.Create(tank)
	this.PanicError(db.Error)

	return tank
}

//修改一个文件
func (this *TankDao) Save(tank *Tank) *Tank {

	tank.ModifyTime = time.Now()
	db := this.context.DB.Save(tank)
	this.PanicError(db.Error)

	return tank
}

//按照Id查询用户，找不到返回nil
func (this *TankDao) FindByUuid(uuid string) *Tank {

	// Read
	var tank *Tank = &Tank{}
	db := this.context.DB.Where(&Tank{Base: Base{Uuid: uuid}}).First(tank)
	if db.Error != nil {
		return nil
	}
	return tank
}

//按照Id查询用户,找不到抛panic
func (this *TankDao) CheckByUuid(uuid string) *Tank {

	// Read
	var tank *Tank = &Tank{}
	db := this.context.DB.Where(&Tank{Base: Base{Uuid: uuid}}).First(tank)
	this.PanicError(db.Error)
	return tank
}
