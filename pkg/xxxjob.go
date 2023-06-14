package pkg

import "sync"

// 定时任务，定期删除map
func DeleteMap(userMap sync.Map) {
	// 定期遍历，之后根据开始时间超过一定时间未删除则在这里处理
	userMap.Range(func(key, value interface{}) bool {
		// TODO
		//userInfo := value.(*)
		// 返回true继续遍历，false则停止
		return true
	})
	return
}
