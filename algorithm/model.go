package algorithm

// 传入算法结构，有什么需要调整的就是在这里加结构
type UserInfo struct {
	// 词汇量
	Score int64
	// 所有难度总答题数
	TotalNum int64
	// 这个是答题到目前所有难度的相关信息
	LadderInfo     map[string]*LadderInfo
	VocabularyInfo *VocabularyInfo
}

// 阶梯相关数据
type LadderInfo struct {
	// 该难度答题数
	curNum int64
	// 该难度认识数
	knownNun int64
}

// 词汇相关数据
type VocabularyInfo struct {
	wordId int64
	word   string
	level  int64
	known  bool
}

type Ladder struct {
	// 在当前阶梯（难度）所答的词数
	curNum int64
	// 在当前阶梯（难度）所认识的词数
	realizeNum int64
	// 总共回答的词数
	totalNum int64
	// 当前难度
	level int64
	// 是否结束
	endFlag bool
}

func LaderHandler(curLader *Ladder, vocabulary *VocabularyInfo) *Ladder {
	return new(Ladder)
}
