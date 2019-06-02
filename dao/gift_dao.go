/**
 * 抽奖系统的数据库操作
 */
package dao

import (
	"github.com/go-xorm/xorm"

	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/models"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao {
	return &GiftDao{
		engine: engine,
	}
}

func (d *GiftDao) Get(id int) *models.LtGift {
	data := &models.LtGift{Id: id}
	ok, err := d.engine.Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *GiftDao) GetAll() []models.LtGift {
	datalist := make([]models.LtGift, 0)
	err := d.engine.
		Asc("sys_status").
		Asc("displayorder").
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *GiftDao) CountAll() int64 {
	num, err := d.engine.
		Count(&models.LtGift{})
	if err != nil {
		return 0
	} else {
		return num
	}
}

//func (d *GiftDao) Search(country string) []models.LtGift {
//	datalist := make([]models.LtGift, 0)
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

func (d *GiftDao) Delete(id int) error {
	data := &models.LtGift{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *GiftDao) Update(data *models.LtGift, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *GiftDao) Create(data *models.LtGift) error {
	_, err := d.engine.Insert(data)
	return err
}

// 获取到当前可以获取的奖品列表
// 有奖品限定，状态正常，时间期间内
// gtype倒序， displayorder正序
func (d *GiftDao) GetAllUse() []models.LtGift {
	now := comm.NowUnix()
	datalist := make([]models.LtGift, 0)
	err := d.engine.
		Cols("id", "title", "prize_num", "left_num", "prize_code",
			"prize_time", "img", "displayorder", "gtype", "gdata").
		Desc("gtype").
		Asc("displayorder").
		Where("prize_num>=?", 0).    // 有限定的奖品
		Where("sys_status=?", 0).    // 有效的奖品
		Where("time_begin<=?", now). // 时间期内
		Where("time_end>=?", now).   // 时间期内
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

//递增num个数量
func (d *GiftDao) IncrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).
		Incr("left_num", num).
		//Where("left_num=?", num).
		Update(&models.LtGift{Id: id})
	return r, err
}

//递减num个数量
func (d *GiftDao) DecrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).
		Decr("left_num", num).
		Where("left_num>=?", num).
		Update(&models.LtGift{Id: id})
	return r, err
}
