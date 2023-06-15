package vacabulary

import (
	"encoding/json"
	_alg "estimation-vocabulary/algorithm"
	_internal "estimation-vocabulary/internal"
	_model "estimation-vocabulary/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"path/filepath"
	"time"
)

// 可选等级显示结构
type LevelStruct struct {
	Label string `json:"lable"`
	Value string `json:"value"`
}

// 可选等级列表
var levels = []LevelStruct{
	// TODO 创建一个级别的结构返回,大概三个吧，根据业务看看怎么设置
	// 例子 {Label:"小白",Value:"A1"}
	{Label: "小白", Value: "A1"},
}

// TODO 批处理接收文件格式 - 暂定这种,具体还得取决于批处理的具体含义
/*
	{
		"wordList":{
			{"word":a,"known":false},
			{"word":b,"known":true},
		}
	}
*/

// 批处理接收结构
type Batch struct {
	WordList []VocabularyBatch `json:"wordList"`
}

type VocabularyBatch struct {
	Word  string `json:"word"`
	Known bool   `json:"known"`
}

// 选中单词认识与否的请求结构
type WordKnownReq struct {
	TestId string `json:"testId"`
	WordId string `json:"wordId"`
	Word   string `json:"word"`
	Known  bool   `json:"known"`
}

// @Method GET
// Describe 显示可选的初始级别
func ShowLevelList(c *gin.Context) {
	_internal.ResponseSuccess(c, levels)
}

// @Method Get
// @Param level string
func StartTest(c *gin.Context) {
	// TODO 从前端拿到用户初始所选level
	level := "A1"

	//
	//level := c.Query("level")

	//fmt.Println(level)

	// 2. 创建一个map
	testId := uuid.New().String()
	_internal.UserMap.Store(testId, &_internal.UserTestStruct{
		Level:    level,
		Score:    0,
		TotalNum: 0,
		//VocabularyInfo: &_alg.VocabularyInfo{},
		VocabularyInfo: &_alg.VocabularyInfo{},
		LadderInfo:     make(map[string]*_alg.LadderInfo),
		WordInfo:       make(map[string][]int64),
		StartTime:      time.Now(),
	})

	// 3.返回一个testId
	_internal.ResponseSuccess(c, testId)

}

//	@Method Get
//  @Param test_id

func GetWord(c *gin.Context) {
	// 1. 接收testId,只有一个方法就直接从路由里读取算了

	testId := uuid.New().String()
	var totalNum int64 = 0
	//testId := c.Query("test_id")
	//totalNum := c.Query("total_num")

	// 2. 判断是否存在对应的map，没有直接给他停了
	userInfo, exist := _internal.UserMap.Load(testId)
	if !exist {
		_internal.ResponseError(c, _internal.CodeInvalidTestId)
		return
	}
	user := userInfo.(*_internal.UserTestStruct)

	//TODO 调用算法，获取当前的level，和算法小组沟通需要的参数

	// 3.TODO model那边根据level随机拿一个，且不重复
	v := _model.Vocabulary{
		Level: user.Level,
	}

	// 抽取单词且保证单词不重复
	for {
		err := v.SelectVocabularyByLevelRandom()
		if err != nil {
			_internal.ResponseError(c, _internal.CodeWordSelectErr)
		}

		//判断随机抽出来的单词是否重复
		//TODO 解决uuid存储int问题（uuid是int128）
		ok, err := _internal.JudgeIfRepeated(testId, v.Level, v.Id)
		if err != nil {
			_internal.ResponseError(c, _internal.CodeWordRepeat)
		}
		if ok {
			break
		}
	}

	// 4.TODO 返回需要的数据
	vInfo := _alg.VocabularyInfo{
		WordId: v.Id,
		Word:   v.Word,
	}
	//返回
	_internal.ResponseSuccess(c, _internal.UserTestStruct{
		Level:          user.Level,
		TotalNum:       totalNum,
		VocabularyInfo: &vInfo,
	})
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

	// body 里的 json 解析方法
	wordReq := WordKnownReq{}
	if err := c.ShouldBindJSON(&wordReq); err != nil {
		log.Println("解析body错误", err)
		_internal.ResponseError(c, _internal.CodeErrParseBody)
		return
	}

	// TODO 1.获得请求参数
	testId := uuid.New().String()
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

	user := userTestMap.(*_internal.UserTestStruct)

	// TODO 构造算法需要的结构，具体根据算法需求改.

	userInfo := &_alg.UserInfo{
		Score:      user.Score,
		TotalNum:   user.TotalNum,
		LadderInfo: user.LadderInfo,
		//VocabularyInfo: user.VocabularyInfo,
	}
	fmt.Println(userInfo)

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
	testId := uuid.New().String()
	//testId := c.Query("test_id")
	// 2.TODO 删掉map
	_internal.UserMap.Delete(testId)
	// 3. TODO 返回前端是否成功  是否需要返回分数（待定）
	_internal.ResponseSuccess(c, nil)
}

func Test(c *gin.Context) {
	fmt.Println("test service success")
	_internal.ResponseSuccess(c, nil)
}

// @Method POST
// @Parm form-data
// @Describe 批处理
func GetScoreBatch(c *gin.Context) {
	file, _ := c.FormFile("file")
	// 识别后缀，这里直接限制json
	ext := filepath.Ext(file.Filename)
	if ext != ".json" {
		log.Println("文件类型出错,批处理需要json文件")
		_internal.ResponseError(c, _internal.CodeErrFileFormat)
		return
	}

	src, _ := file.Open()
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		log.Println("读取文件数据错误", err)
		_internal.ResponseError(c, _internal.CodeServerBusy)
		return
	}

	var batchData Batch
	err = json.Unmarshal(data, &batchData)
	if err != nil {
		log.Println("解析json错误", err)
		_internal.ResponseError(c, _internal.CodeErrJsonFormat)
		return
	}

	// TODO 如何根据解析出的json去调用我们自己的方法

	// TODO 计算出最后成绩然后返回

	//_internal.ResponseSuccess()
}

//TODO 逻辑
//TODO 1、建立测试链接：创建testId，加入map，返回testId
//TODO 2、前端接收到链接建立成功的状态码，则发起请求获取单词（需要携带当前的总题目数）
