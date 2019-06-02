/**
 * 抽奖系统的数据库操作
 */
package dao

import (
	"github.com/go-xorm/xorm"

	"iris项目/my_lottery/models"
)

type UserdayDao struct {
	engine *xorm.Engine
}

func NewUserdayDao(engine *xorm.Engine) *UserdayDao {
	return &UserdayDao{
		engine: engine,
	}
}

func (d *UserdayDao) Get(id int) *models.LtUserday {
	data := &models.LtUserday{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserdayDao) GetAll(page, size int) []models.LtUserday {
	offset := (page - 1) * size
	datalist := make([]models.LtUserday, 0)
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

func (d *UserdayDao) CountAll() int64 {
	num, err := d.engine.
		Count(&models.LtUserday{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

func (d *UserdayDao) Search(uid, day int) []models.LtUserday {
	datalist := make([]models.LtUserday, 0)
	err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Desc("id").
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *UserdayDao) Count(uid, day int) int {
	info := &models.LtUserday{}
	ok, err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Get(info)
	if !ok || err != nil {
		return 0
	} else {
		return info.Num
	}
}

//func (d *UserdayDao) Delete(id int) error {
//	data := &models.LtUserday{Id: id, SysStatus: 1}
//	_, err := d.engine.Id(data.Id).Update(data)
//	return err
//}

func (d *UserdayDao) Update(data *models.LtUserday, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *UserdayDao) Create(data *models.LtUserday) error {
	_, err := d.engine.Insert(data)
	return err
}
