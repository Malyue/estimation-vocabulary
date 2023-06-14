package algorithm

import (
	"encoding/csv"
	_model "estimation-vocabulary/internal/model"
	"log"
	"os"
)

/*
	StartTime: 2023/6/12
	Author:
*/

// 将数据集导入到数据库
const path = "resource/Result.csv"

// 读取resource中的文件插入数据库
func ImportToDb() {
	importVocabulary()
}

// 转换vocabulary表
func importVocabulary() {
	// 读取文件
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("无法打开csv", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		// 具体怎么插入需要根据文件的列数来看
		v := _model.Vocabulary{
			Word:  record[0],
			Level: record[1],
			// 词频排序
			//FrequenceLevel: record[2],
		}
		// FrequenceLevel看需要添加，然后在结构体中加上即可
		v.InsertVocabulary()
		//_model.InsertVocabulary(record[0], record[1], 0)
	}

}
