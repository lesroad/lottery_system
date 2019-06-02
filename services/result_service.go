/**
 * 抽奖系统数据处理（包括数据库，也包括缓存等其他形式数据）
 */
package services

import (
	"iris项目/my_lottery/dao"
	"iris项目/my_lottery/datasource"
	"iris项目/my_lottery/models"
)

type ResultService interface {
	GetAll(page, size int) []models.LtResult
	CountAll() int64
	GetNewPrize(size int, giftIds []int) []models.LtResult
	SearchByGift(giftId, page, size int) []models.LtResult
	SearchByUser(uid, page, size int) []models.LtResult
	CountByGift(giftId int) int64
	CountByUser(uid int) int64
	Get(id int) *models.LtResult
	Delete(id int) error
	Update(user *models.LtResult, columns []string) error
	Create(user *models.LtResult) error
}

type resultService struct {
	dao *dao.ResultDao
}

func NewResultService() ResultService {
	return &resultService{
		dao: dao.NewResultDao(datasource.InstanceDbMaster()),
	}
}

func (s *resultService) GetAll(page, size int) []models.LtResult {
	return s.dao.GetAll(page, size)
}

func (s *resultService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *resultService) GetNewPrize(size int, giftIds []int) []models.LtResult {
	return s.dao.GetNewPrize(size, giftIds)
}

func (s *resultService) SearchByGift(giftId, page, size int) []models.LtResult {
	return s.dao.SearchByGift(giftId, page, size)
}

func (s *resultService) SearchByUser(uid, page, size int) []models.LtResult {
	return s.dao.SearchByUser(uid, page, size)
}

func (s *resultService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *resultService) CountByUser(uid int) int64 {
	return s.dao.CountByUser(uid)
}

func (s *resultService) Get(id int) *models.LtResult {
	return s.dao.Get(id)
}

func (s *resultService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *resultService) Update(data *models.LtResult, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *resultService) Create(data *models.LtResult) error {
	return s.dao.Create(data)
}
