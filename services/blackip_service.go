/**
 * 抽奖系统数据处理（包括数据库，也包括缓存等其他形式数据）
 */
package services

import (
	"fmt"
	"iris项目/lottery/comm"
	"iris项目/my_lottery/dao"
	"iris项目/my_lottery/datasource"
	"iris项目/my_lottery/models"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

// IP信息，可以缓存(本地或者redis)，有更新的时候，再根据具体情况更新缓存
var cachedBlackipList = make(map[string]*models.LtBlackip)
var cachedBlackipLock = sync.Mutex{}

type BlackipService interface {
	GetAll(page, size int) []models.LtBlackip
	CountAll() int64
	Search(ip string) []models.LtBlackip
	Get(id int) *models.LtBlackip
	//Delete(id int) error
	Update(user *models.LtBlackip, columns []string) error
	Create(user *models.LtBlackip) error
	GetByIp(ip string) *models.LtBlackip
}

type blackipService struct {
	dao *dao.BlackipDao
}

func NewBlackipService() BlackipService {
	return &blackipService{
		dao: dao.NewBlackipDao(datasource.InstanceDbMaster()),
	}
}

func (s *blackipService) GetAll(page, size int) []models.LtBlackip {
	return s.dao.GetAll(page, size)
}

func (s *blackipService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackipService) Search(ip string) []models.LtBlackip {
	return s.dao.Search(ip)
}

func (s *blackipService) Get(id int) *models.LtBlackip {
	return s.dao.Get(id)
}

//func (s *blackipService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

func (s *blackipService) Update(data *models.LtBlackip, columns []string) error {
	// 先更新缓存的数据
	s.updateByCache(data, columns)
	// 再更新数据的数据
	return s.dao.Update(data, columns)
}

func (s *blackipService) Create(data *models.LtBlackip) error {
	return s.dao.Create(data)
}

// 根据IP读取IP的黑名单数据
func (s *blackipService) GetByIp(ip string) *models.LtBlackip {
	// 先从缓存中读取数据
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		// 再从数据库中读取数据
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.LtBlackip{Ip: ip}
		}
		s.setByCache(data)
	}
	return data
}

func (s *blackipService) getByCache(ip string) *models.LtBlackip {
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", ip)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blackip_service.getByCache HGETALL key=", key, ", error=", err)
		return nil
	}
	dataIp := comm.GetStringFromStringMap(dataMap, "Ip", "")
	if dataIp == "" {
		return nil
	}
	data := &models.LtBlackip{
		Id:         int(comm.GetInt64FromStringMap(dataMap, "Id", 0)),
		Ip:         dataIp,
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
	}
	return data
}

func (s *blackipService) setByCache(data *models.LtBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()
	// 数据更新到redis缓存
	params := []interface{}{key}
	params = append(params, "Ip", data.Ip)
	if data.Id > 0 {
		params = append(params, "Blacktime", data.Blacktime)
		params = append(params, "SysCreated", data.SysCreated)
		params = append(params, "SysUpdated", data.SysUpdated)
	}
	_, err := rds.Do("HMSET", params...)
	if err != nil {
		log.Println("blackip_service.setByCache HMSET params=", params, ", error=", err)
	}
}

// 数据更新了，直接清空缓存数据
func (s *blackipService) updateByCache(data *models.LtBlackip, columns []string) {
	if data == nil || data.Ip == "" {
		return
	}
	// 集群模式，redis缓存
	key := fmt.Sprintf("info_blackip_%s", data.Ip)
	rds := datasource.InstanceCache()
	// 删除redis中的缓存
	rds.Do("DEL", key)
}
