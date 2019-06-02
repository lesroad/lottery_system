/*
// 10 不同编码的优惠券发放
*/
package utils

import (
	"fmt"
	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/datasource"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"
	"log"
)

// 不同编码优惠券类的发放 redis+数据库
func PrizeCodeDiff(id int, codeService services.CodeService) string {
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SPOP", key) //随机移除key中元素，返回被移除的元素（唯一编码的优惠券）
	if err != nil {
		log.Println("prizedata.prizeServCodeDiff error=", err)
		return ""
	}
	code := comm.GetString(rs, "")
	if code == "" {
		log.Printf("prizedata.prizeServCodeDiff rs=%s", rs)
		return ""
	}
	// 更新数据库中的发放状态
	codeService.UpdateByCode(&models.LtCode{
		Code:       code,
		SysStatus:  2,
		SysUpdated: comm.NowUnix(),
	}, nil)
	return code
}

// 导入新的优惠券编码
func ImportCacheCodes(id int, code string) bool {
	// 集群版本需要放入到redis中
	// [暂时]本机版本的就直接从数据库中处理吧
	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.RecacheCodes SADD error=", err)
		return false
	} else {
		return true
	}
}

// 数据库的优惠券的编码整理到缓存中
func RecacheCodes(id int, codeService services.CodeService) (sucNum, errNum int) {
	// 集群版本需要放入到redis中
	// [暂时]本机版本的就直接从数据库中处理吧
	list := codeService.Search(id) //数据库中得到所有的券
	if list == nil || len(list) <= 0 {
		return 0, 0
	}
	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	tmpKey := "tmp_" + key //使用临时名然后覆盖
	for _, data := range list {
		if data.SysStatus == 0 { //状态正常
			code := data.Code
			_, err := cacheObj.Do("SADD", tmpKey, code)
			if err != nil {
				log.Println("prizedata.RecacheCodes SADD error=", err)
				errNum++
			} else {
				sucNum++
			}
		}
	}
	_, err := cacheObj.Do("RENAME", tmpKey, key)
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	}
	return sucNum, errNum
}

// 获取当前的缓存中编码数量
// 返回，剩余数据库编码数量和缓存中编码数量
func GetCacheCodeNum(id int, codeService services.CodeService) (int, int) {
	num := 0
	cacheNum := 0
	// 统计数据库中有效编码数量
	list := codeService.Search(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				num++
			}
		}
	}

	// redis中缓存的key值
	key := fmt.Sprintf("gift_code_%d", id)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("SCARD", key) //返回set数量
	if err != nil {
		log.Println("prizedata.RecacheCodes RENAME error=", err)
	} else {
		cacheNum = int(comm.GetInt64(rs, 0))
	}

	return num, cacheNum
}

/*
// 从数据库中发放券(已被缓存替代)
func PrizeLocalCodeDiff(id int, codeService services.CodeService) string {
	lockUid := 0 - id - 100000000
	LockLucky(lockUid)
	defer UnlockLucky(lockUid)

	codeId := 0
	codeInfo := codeService.NextUsingCode(id, codeId)
	if codeInfo != nil && codeInfo.Id > 0 {
		codeInfo.SysStatus = 2
		codeInfo.SysUpdated = comm.NowUnix()
		codeService.Update(codeInfo, nil)
	} else {
		log.Println("prizedata.prizeCodeDiff num codeInfo, gift_id=", id)
		return ""
	}
	return codeInfo.Code
}
*/
