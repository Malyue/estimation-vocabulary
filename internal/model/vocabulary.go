package model

import "time"

// vocabularies
type Vocabulary struct {
	Id             int           `json:"id" gorm:"column:id"`
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
	// 需要保证不重复

}

func (v *Vocabulary) SelectByID() error {
	result := db.Model(&Vocabulary{}).Where("id=?", id).Select(&v)

	if result.Error != nil {
		return result.Error
	}
}
