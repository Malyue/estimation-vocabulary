package internal

import (
	"errors"
	_pkg "estimation-vocabulary/algorithm"
	_model "estimation-vocabulary/internal/model"
	"fmt"
	"sync"
	"time"
)

// 储存 test-id 以及 对应的一些信息，比如
/*
	"A1":{
		Known:10,
		Un:20,
		WordList:{1223,7321}
	}
*/

//var UserMap map[string]UserTestStruct

// map并发读写不安全,改用sync.Map
// 存储 - m.Store(key,value)
// 获取 value,ok := m.Load(key)
// 删除 m.Delete(key)

// 由于是并发安全的，并发情况下会影响性能，具体影响可能得测试看看（我们没什么并发应该问题不大）
var UserMap sync.Map

// 存储用户一次测试的相关内容(定期清理)
// 用redis也可以，看哪个省事
type UserTestStruct struct {
	// 词汇量
	Score int
	// 当前等级 A1 - C2
	Level string
	// 当前答题数
	TotalNum int64
	//VocabularyInfo *_pkg.VocabularyInfo
	//TODO 确定Vocabulary结构体的类别
	VocabularyInfo *_model.Vocabulary
	// 所有难度的对应数据，key为level
	LadderInfo map[string]*_pkg.LadderInfo
	// 详细信息，包括每个阶段的认识数目以及对应的wordlist,level作为key,wordId为value
	WordInfo  map[string][]int64
	StartTime time.Time
}

// 判断wordList
func JudgeIfRepeated(testId string, level string, wordId int64) (bool, error) {
	// 先判断存不存在testid对应的map

	userInfo, ok := UserMap.Load(testId)
	if !ok {
		return false, errors.New("the testId is not exists")
	}

	var user *UserTestStruct
	// 取值
	// 判断是否该类型,同类型赋值，不同类型报错
	switch userInfo.(type) {
	case *UserTestStruct:
		user = userInfo.(*UserTestStruct)
	default:
		return false, errors.New("userInfo is not *UserTestStruct type")
	}
	fmt.Println(user)

	// 这里需要WordInfo
	wordIdList, ok := user.WordInfo[level]
	if !ok {
		// 不存在这个level，可能是处于刚升级的情况
		// TODO 这一步需要根据业务逻辑看看是否需要更改
		//暂时不需要改
		return true, nil
	}
	// 查看wordId是否存在，这个看一下有没有优化的方法
	for _, value := range wordIdList {
		if value == wordId {
			return true, nil
		}
	}

	return false, nil
}
