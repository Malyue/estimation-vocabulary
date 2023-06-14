package vacabulary

import (
	_pkg "estimation-vocabulary/algorithm"
	_internal "estimation-vocabulary/internal"
	_model "estimation-vocabulary/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func StartTest(c *gin.Context) {
	// TODO 从前端拿到用户所选level
	level := "A1"

	// 2. 创建一个map
	testId := uuid.New()
	_internal.UserMap.Store(testId, &_internal.UserTestStruct{
		Level:          level,
		Score:          0,
		TotalNum:       0,
		VocabularyInfo: &_pkg.VocabularyInfo{},
		LadderInfo:     make(map[string]*_pkg.LadderInfo),
		WordInfo:       make(map[string][]int64),
		StartTime:      time.Now(),
	})

	// 3.返回一个testId
	_internal.ResponseSuccess(c, testId)

}

func GetWord(c *gin.Context) {
	// 1.TODO 接收testId
	testId := uuid.New()

	// 2.TODO 可以判断是否存在对应的map，没有直接给他停了
	userInfo, exist := _internal.UserMap.Load(testId)
	if !exist {
		_internal.ResponseError(c, _internal.CodeInvalidTestId)
		return
	}
	user := userInfo.(*_internal.UserTestStruct)

	// 3.TODO model那边根据level随机拿一个
	v := _model.Vocabulary{
		Level: user.Level,
	}
	// TODO 函数未实现
	err := v.SelectVocabularyByLevelRandom()
	if err != nil {

	}

	// 4.TODO 返回需要的数据
	_internal.ResponseSuccess(c, _internal.UserTestStruct{})
	// testid
	// 随机获得，
	/*
		{
			"testid":xx,
			"wordId":1,
			"word":a
		}
	*/
}

func GetScore(c *gin.Context) {
	// TODO 改一下函数名字，没想好，作用是获取用户每个单词的认识与否,并告知算法
	/*
		testid:xx,
		wordId:1,
		Known:true/false
	*/
	// TODO 1.获得请求参数
	testId := uuid.New()
	/*
		调用算法
		wordId:1,
		curNum:2,
		curKnown:0,
		Known:true/false,
		Score:3000,
	*/
	// TODO 2.从全局map获取当前testId的一些数据，然后构造算法需要的结构
	userTestMap, exist := _internal.UserMap.Load(testId)
	if !exist {
		_internal.ResponseError(c, _internal.CodeInvalidTestId)
		return
	}

	_ = userTestMap.(*_internal.UserTestStruct)

	// TODO 构造算法需要的结构，具体根据算法需求改.
	//userInfo := &_pkg.UserInfo{
	//	Score:          user.Score,
	//	TotalNum:       user.TotalNum,
	//	LadderInfo:     user.LadderInfo,
	//	VocabularyInfo: user.VocabularyInfo,
	//}

	// TODO 3.调用算法层，参数统一为UserInfo结构,具体怎么调用看算法层的方法,然后根据返回结构去修改全局map的信息
	// ladderInfo,exist := _internal.UserMap[testId]

	// TODO 4.返回前端，告知请求成功，正常的话不需要数据返回
	/*
		return
		score:
		level:
	*/
}

// 接口
func GetResult(c *gin.Context) {
	// 返回结果
	// 1. TODO 获得testid

	// 2. TODO 两种可能吧.大概率是直接从map里读取score直接返回即可，或者看是否要从算法层再计算一遍
	//_internal.ResponseSuccess()
}

func Exit(c *gin.Context) {
	// 1.TODO 获取testid
	testId := uuid.New()
	// 2.TODO 删掉map
	_internal.UserMap.Delete(testId)
	// 3. TODO 返回前端是否成功
}
