package algorithm

// TODO 应该遍历所有阶梯 在用户点击结束的时候才调用
// @Title forecastVocabulary
// @Description
// @Param curLadder
// @Return int64
func forecastVocabulary(userinfo *UserInfo) {
	result := int64(0)
	// TODO 计算分数
	// 反应不同阶级回答情况的表
	var diffTagIdentifyMap map[string]*LadderInfo
	diffTagIdentifyMap = userinfo.LadderInfo
	// 1 遍历用户再不同阶层的作答情况
	for curLadderID, curLadderInfo := range diffTagIdentifyMap {
		// 当前阶层词汇量的上界
		upperBound := levelVocabulary[curLadderID][1]
		// 当前阶层词汇量的下界
		lowerBound := levelVocabulary[curLadderID][0]
		// 当前阶层认识的词
		realizeNum := curLadderInfo.knownNun
		// 当前阶层回答的词的总数
		curSum := curLadderInfo.curNum
		// 当前阶层回答的词的正确率
		rate := 1.0 * realizeNum / curSum
		// 预测当前阶层知道的词的总数
		estVcbSize := lowerBound + ((upperBound - lowerBound) * rate)
		// 当前阶层的权重
		weight := 1.0 * curSum / userinfo.TotalNum

		// 当前阶层在加权后知道的词的总数
		result += estVcbSize * weight

	}
	userinfo.Score = result
}

// @Title ladderHandler
// @Description   每次提交都要调用，根据单词的认识情况判断是否切换阶层
// @Param curLadder 当前阶层信息
// @Param vocabulary 当前单词信息
// @Return LadderInfo 修改过后的阶层信息
func ladderHandler(userinfo *UserInfo) {
	// TODO 有可能需要更改
	if userinfo.EndFlag {
		forecastVocabulary(userinfo)
	}
	level := userinfo.Level
	// 触发阶梯变化最少的单词个数
	baseChangeNum := int64(5)

	// 触发阶梯变化的临界识别率
	baseRealizeRate := 0.2
	if userinfo.VocabularyInfo.known {
		// 认识单词
		baseRealizeRate = 0.8
		userinfo.LadderInfo[level].knownNun++
	}
	userinfo.TotalNum++
	userinfo.LadderInfo[level].curNum++
	// TODO 记录当前词汇
	// 计算用户在当前阶段的认识率
	var realizeRate = float64(userinfo.LadderInfo[level].knownNun) / float64(userinfo.LadderInfo[level].curNum)
	if userinfo.VocabularyInfo.known && realizeRate >= baseRealizeRate && userinfo.LadderInfo[level].curNum >= baseChangeNum {
		upgradeLadder(userinfo)
	}
	if !userinfo.VocabularyInfo.known && realizeRate <= baseRealizeRate && userinfo.LadderInfo[level].curNum >= baseChangeNum {
		downgradeLadder(userinfo)
	}
}

// @Title upgradeLadder
// @Description  降低难度
// @Param curLadder
// @Return LadderInfo
func upgradeLadder(userinfo *UserInfo) {
	level := level2Num[userinfo.Level]
	if level == 5 {
		return
	}
	level++
	userinfo.Level = num2Level[level]
}

// @Title downgradeLadder
// @Description 提高难度
// @Param curLadder
// @Return LadderInfo
func downgradeLadder(userinfo *UserInfo) {
	level := level2Num[userinfo.Level]
	if level == 0 {
		return
	}
	level--
	userinfo.Level = num2Level[level]
}
