package algorithm

// 用于存储每个等级的上下限
var levelVocabulary = map[string][2]int64{
	"A1": {0, 1000},
	"A2": {1000, 2000},
	"B1": {2000, 3000},
	"B2": {3000, 5000},
	"C1": {5000, 8000},
	"C2": {8000, 10000},
}
var level2Num = map[string]int64{
	"A1": 0,
	"A2": 1,
	"B1": 2,
	"B2": 3,
	"C1": 4,
	"C2": 5,
}

var num2Level = [6]string{"A1", "A2", "B1", "B2", "C1", "C2"}

// 传入算法结构，有什么需要调整的就是在这里加结构
type UserInfo struct {
	// 词汇量
	Score int64
	// 所有难度总答题数
	TotalNum int64
	// 这个是答题到目前所有难度的相关信息
	LadderInfo     map[string]*LadderInfo
	VocabularyInfo *VocabularyInfo
	// 是否结束 初始化应该为false
	EndFlag bool
	// 用户当前的答题难度
	Level string
}

// 阶梯相关数据
type LadderInfo struct {
	// 该难度答题数
	CurNum int64
	// 该难度认识数
	KnownNum int64
}

// 词汇相关数据
type VocabularyInfo struct {
	WordId int64
	Word   string
	Known  bool
}

type Ladder struct {
	// 在当前阶梯（难度）所答的词数
	CurNum int64
	// 在当前阶梯（难度）所认识的词数
	RealizeNum int64
	// 总共回答的词数
	TotalNum int64
	// 是否结束
	EndFlag bool
}

func LaderHandler(curLader *Ladder, vocabulary *VocabularyInfo) *Ladder {
	return new(Ladder)
}
