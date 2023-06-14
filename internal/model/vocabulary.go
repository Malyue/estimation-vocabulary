package model

import (
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// vocabularies
type Vocabulary struct {
	Id             int64         `json:"id" gorm:"column:id"`
	Word           string        `json:"word" gorm:"column:word"`
	Level          string        `json:"level" gorm:"column:level"`
	FrequenceLevel int           `json:"frequence_level" gorm:"column:frequence_level"`
	CreateAt       time.Duration `json:"createAt" `
	UpdateAt       time.Duration `json:"updateAt"`
	DeleteFlag     int           `json:"delete_flag" gorm:"column:delete_flag"`
}

func (v *Vocabulary) InsertVocabulary() (err error) {
	v.DeleteFlag = 0
	result := db.Model(&Vocabulary{}).Create(&v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 根据等级随机取查找单词
func (v *Vocabulary) SelectVocabularyByLevelRandom() error {
	// 在业务逻辑层保证抽取的单词不重复，这里只负责随机抽取
	//利用Gorm设置随机数种子进行随机抽取

	//设置随机数种子
	rand.Seed(time.Now().UnixNano())

	err := db.Model(&Vocabulary{}).Where("level =?", v.Level).
		Order(gorm.Expr("RAND()")).
		Limit(1).
		First(v).
		Error
	if err != nil {
		return err
	}
	return nil
}

//func (v *Vocabulary) SelectByID() error {
//	result := db.Model(&Vocabulary{}).Where("id=?", id).Select(&v)
//
//	if result.Error != nil {
//		return result.Error
//	}
//}
