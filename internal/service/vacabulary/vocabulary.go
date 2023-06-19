package vacabulary

import (
	"bytes"
	"encoding/json"
	_alg "estimation-vocabulary/algorithm"
	_internal "estimation-vocabulary/internal"
	_model "estimation-vocabulary/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
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
	{Label: "", Value: ""},
	{Label: "", Value: ""},
}

// 用于存储每个等级的上下限
var levelVocabulary = map[string][2]int64{
	"A1": {0, 1000},
	"A2": {1000, 2000},
	"B1": {2000, 3000},
	"B2": {3000, 5000},
	"C1": {5000, 8000},
	"C2": {8000, 10000},
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
	BroadVocabTest  []VocabularyBatch `json:"broad_vocab_test"`
	NarrowVocabTest []VocabularyBatch `json:"narrow_vocab_test"`
	Result          int               `json:"result"`
	//WordList []VocabularyBatch `json:"wordList"`
}

type VocabularyBatch struct {
	//Word  string `json:"word"`
	Value string `json:"value"`
	Known bool   `json:"known"`
}

// 选中单词认识与否的请求结构
type WordKnownReq struct {
	TestId string `json:"test_id"`
	WordId string `json:"word_id"`
	Word   string `json:"word"`
	Known  bool   `json:"known"`
}

// 获取单词的请求结构 GetWord（）
type WordGetReq struct {
	TestId string `json:"test_id"`
	//TotalNum int64  `json:"total_num"`
}

// 获取单词的响应结构
type WordGetRes struct {
	TestId   string `json:"test_id"`
	WordId   string `json:"word_id"`
	Word     string `json:"word"`
	TotalNum int64  `json:"total_num"`
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
	level := c.Query("level")
	if level == "" {
		_internal.ResponseError(c, _internal.CodeLevelRequire)
		return
	}
	if _, ok := levelVocabulary[level]; !ok {
		_internal.ResponseError(c, _internal.CodeWrongLevel)
		return
	}
	score := int64(0)

	//TODO 根据level设置初始的Score（Score初始值为当前等级的区间下限）

	// 2. 创建一个map
	testId := uuid.New().String()
	user := _internal.UserTestStruct{
		Level:    level,
		Score:    score,
		TotalNum: 0,
		//VocabularyInfo: &_alg.VocabularyInfo{},
		VocabularyInfo: &_alg.VocabularyInfo{},
		LadderInfo:     make(map[string]*_alg.LadderInfo),
		WordInfo:       make(map[string][]int64),
		EndFlag:        false,
		StartTime:      time.Now(),
	}

	_internal.UserMap.Store(testId, &user)
	// 3.返回一个testId

	c.Set("test_id", testId)

	if batch, _ := c.Get("batch"); batch == true {
		return
	}
	_internal.ResponseSuccess(c, testId)

}

//	@Method Get
//  @Param test_id
//  @Return word,total_num,

func GetWord(c *gin.Context) {
	// 1. TODO 接收testId,只有一个方法就直接从路由里读取算了

	testId := uuid.New().String()
	//var totalNum int64 = 0
	testId = c.Query("test_id")
	//totalNum = c.Query("total_num")

	// 2. 判断是否存在对应的map，没有直接给他停了
	userInfo, exist := _internal.UserMap.Load(testId)
	if !exist {
		_internal.ResponseError(c, _internal.CodeInvalidTestId)
		return
	}
	user := userInfo.(*_internal.UserTestStruct)

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
			log.Println(err)
			_internal.ResponseError(c, _internal.CodeWordRepeat)
			return
		}
		if !ok {
			break
		}
	}

	//TODO WordInfo记录已返回的单词,这里是不是要放在单词被算法处理之后
	user.WordInfo[user.Level] = append(user.WordInfo[user.Level], v.Id)
	//修改user信息
	user.VocabularyInfo = &_alg.VocabularyInfo{
		WordId: v.Id,
		Word:   v.Word,
		Known:  false, //先初始化为false
	}
	// 4.TODO 返回需要的数据
	res := WordGetRes{
		TestId:   testId,
		WordId:   strconv.FormatInt(v.Id, 10),
		Word:     v.Word,
		TotalNum: user.TotalNum,
	}
	if batch, _ := c.Get("batch"); batch == true {
		return
	}
	//返回
	_internal.ResponseSuccess(c, res)
}

// @Description 接收对于单词的认识与否，调用预测算法，重新计算Level
// @Method POST

func UpdateLevel(c *gin.Context) {
	//获取用户每个单词的认识与否,并告知算法
	/*
		testid:xx,
		wordId:1,
		Known:true/false
	*/

	// TODO 1.获得请求参数
	//testId := uuid.New().String()
	// body 里的 json 解析方法
	wordReq := WordKnownReq{}
	if err := c.ShouldBindJSON(&wordReq); err != nil {
		log.Println("解析body错误", err)
		_internal.ResponseError(c, _internal.CodeErrParseBody)
		return
	}
	testId := wordReq.TestId
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

	//TODO 更新userTestStruct的VocabularyInfo
	//wordId 将string转int64
	wordIdInt, err := strconv.ParseInt(wordReq.WordId, 10, 64)
	if err != nil {
		_internal.ResponseError(c, _internal.CodeErrParseInt)
		return
	}

	user.VocabularyInfo = &_alg.VocabularyInfo{
		WordId: wordIdInt,
		Word:   wordReq.Word,
		Known:  wordReq.Known,
	}

	// TODO 构造算法需要的结构，具体根据算法需求改.

	userInfo := &_alg.UserInfo{
		Score:          user.Score,
		TotalNum:       user.TotalNum,
		LadderInfo:     user.LadderInfo,
		VocabularyInfo: user.VocabularyInfo,
		EndFlag:        user.EndFlag,
		Level:          user.Level,
	}
	//fmt.Println(userInfo)

	// TODO 3.调用算法层，参数统一为UserInfo结构,具体怎么调用看算法层的方法,然后根据返回结构去修改全局map的信息
	// ladderInfo,exist := _internal.UserMap[testId]
	//调用算法层
	ok, err := _alg.LadderHandler(userInfo)
	if !ok {
		_internal.ResponseErrorWithData(c, _internal.CodeLevelInvalid, err.Error())
		return
	}

	//覆盖算法返回结果
	user.Score = userInfo.Score
	user.TotalNum = userInfo.TotalNum
	user.Level = userInfo.Level

	// TODO 4.返回前端，告知请求成功，正常的话不需要数据返回

	if batch, _ := c.Get("batch"); batch == true {
		return
	}
	_internal.ResponseSuccess(c, nil)
}

// 接口
func GetResult(c *gin.Context) {
	// 返回结果
	// 1. TODO 获得testid
	testId := c.Query("test_id")
	if testId == "" {
		a, _ := c.Get("test_id")
		testId = a.(string)
	}

	userTestMap, exist := _internal.UserMap.Load(testId)
	if !exist {
		_internal.ResponseError(c, _internal.CodeInvalidTestId)
		return
	}
	user := userTestMap.(*_internal.UserTestStruct)
	user.EndFlag = true
	// 2. TODO 调用forcastVocabulary（）
	userInfo := &_alg.UserInfo{
		Score:          user.Score,
		TotalNum:       user.TotalNum,
		LadderInfo:     user.LadderInfo,
		VocabularyInfo: user.VocabularyInfo,
		EndFlag:        user.EndFlag,
		Level:          user.Level,
	}

	ok, err := _alg.LadderHandler(userInfo)
	if !ok {
		_internal.ResponseErrorWithData(c, _internal.CodeLevelInvalid, err.Error())
		return
	}
	//这里省略赋值回userTestStruct,直接返回
	user.Score = userInfo.Score
	user.TotalNum = userInfo.TotalNum
	user.LadderInfo = userInfo.LadderInfo
	user.Level = userInfo.Level

	score := user.Score
	//if batch, _ := c.Get("batch"); batch == true {
	//	return
	//}

	_internal.ResponseSuccess(c, score)
	return
}

func Exit(c *gin.Context) {
	// 1.TODO 获取testid
	testId := uuid.New().String()
	testId = c.Query("test_id")
	// 2.TODO 删掉map
	_internal.UserMap.Delete(testId)
	// 3. TODO 返回前端是否成功
	if batch, _ := c.Get("batch"); batch == true {
		return
	}
	_internal.ResponseSuccess(c, nil)
	return
}

func Test(c *gin.Context) {
	fmt.Println("test service success")
	_internal.ResponseSuccess(c, nil)
	return
}

// @Method POST
// @Parm form-data
// @Describe 批处理
func GetScoreBatch(c *gin.Context) {

	//TODO 新增一个全局count计量查不到的单词数量，如果占比较多，说明这组数据没有意义，则直接退出，报错。
	vacCount := float64(0)

	file, _ := c.FormFile("file")
	c.Set("batch", true)
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

	// 创建一个map,这里自己设置从哪个等级开始
	queryParams, _ := url.ParseQuery(c.Request.URL.RawQuery)
	//c.Request.URL.Query().Set("level", "A1")
	queryParams.Set("level", "A1")
	c.Request.URL.RawQuery = queryParams.Encode()
	cCopy := c.Copy()
	c.Request.URL.RawQuery = cCopy.Request.URL.RawQuery
	StartTest(c)

	testid, _ := c.Get("test_id")
	testId := testid.(string)

	// TODO 如何根据解析出的json去调用我们自己的方法
	vocabularyList := batchData.BroadVocabTest

	vocabularyList = append(vocabularyList, batchData.NarrowVocabTest...)
	// 拆分一下单词构造一下然后逐个调用接口
	for _, vocabulary := range vocabularyList {
		v := &_model.Vocabulary{
			Word: vocabulary.Value,
		}
		// 根据名称查wordid
		// TODO 如果查无此单词，则vacCount++ 和 continue
		err = v.LoadByWord()
		if err != nil {
			if err.Error() == "RecordNotFound" {
				vacCount++
				log.Println(err.Error())
				continue
			}
			//TODO 这里可以直接返回报错，或者直接continue继续执行
			log.Println(err)
			_internal.ResponseError(c, _internal.CodeWordSelectErr)
			//continue
		}
		// 去调用提交单词接口
		jsonData := map[string]interface{}{
			"test_id": testId,
			"word_id": fmt.Sprint(v.Id),
			"word":    v.Word,
			"known":   vocabulary.Known,
		}
		requestBody, err := json.Marshal(jsonData)
		if err != nil {
			log.Println(err)
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

		UpdateLevel(c)
	}

	//TODO 判断vacCount占用比例，如果
	wordInvalidRate := vacCount / float64(len(vocabularyList))
	if wordInvalidRate >= 0.3 {
		_internal.ResponseError(c, _internal.CodeOops)
	}
	// TODO 计算出最后成绩然后返回
	GetResult(c)
	Exit(c)
	log.Println(batchData.Result)
	//_internal.ResponseSuccess()

}

//TODO 逻辑
//TODO 1、建立测试链接：创建testId，加入map，返回testId

//TODO 2、前端接收到链接建立成功的状态码，则发起请求获取单词（需要携带当前的总题目数）

//TODO 3、获取单词之后，前端发起单词认识辨别的请求，后端调用算法对该用户的level进行更新
//TODO 4、前端受到level更新成功的返回后，继续调用获取单词的请求
//TODO 5、当前端获取单词到达某一数目的时候(该计数器由前端和后端一同保持),前端可以选择发起结束请求
//TODO 6、后端处理结束请求，调用算法的forecastVocabulary函数，然后将Score返回给前端展示

type MapResp struct {
	TestId   string                    `json:"test_id"`
	UserInfo *_internal.UserTestStruct `json:"user_info"`
}

// 获得当前全局变量
func GetMap(c *gin.Context) {
	resp := []MapResp{}
	_internal.UserMap.Range(func(key, value interface{}) bool {
		if value == nil {
			return false
		}
		res := MapResp{}
		res.TestId = key.(string)
		res.UserInfo = value.(*_internal.UserTestStruct)

		resp = append(resp, res)
		return true
	})
	_internal.ResponseSuccess(c, resp)
}
