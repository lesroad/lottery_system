/**
 * 抽奖系统的数据库操作
 */
package dao

import (
	"github.com/go-xorm/xorm"

	"iris项目/my_lottery/models"
)

type UserDao struct {
	engine *xorm.Engine
}

func NewUserDao(engine *xorm.Engine) *UserDao {
	return &UserDao{
		engine: engine,
	}
}

func (d *UserDao) Get(id int) *models.LtUser {
	data := &models.LtUser{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserDao) GetAll(page, size int) []models.LtUser {
	offset := (page - 1) * size
	datalist := make([]models.LtUser, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *UserDao) CountAll() int {
	num, err := d.engine.
		Count(&models.LtUser{})
	if err != nil {
		return 0
	} else {
		return int(num)
	}
}

//func (d *UserDao) Search(country string) []models.LtUser {
//	datalist := make([]models.LtUser, 0)
//	err := d.engine.
//		Where("country=?", country).
//		Desc("id").
//		Find(&datalist)
//	if err != nil {
//		return datalist
//	} else {
//		return datalist
//	}
//}

//func (d *UserDao) Delete(id int) error {
//	data := &models.LtUser{Id: id, SysStatus: 1}
//	_, err := d.engine.Id(data.Id).Update(data)
//	return err
//}

func (d *UserDao) Update(data *models.LtUser, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserDao) Create(data *models.LtUser) error {
	_, err := d.engine.Insert(data)
	return err
}
