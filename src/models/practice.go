package models

import (
	"fmt"
	"tcloud-api/src/util"

	"github.com/jinzhu/gorm"
)

type Practice struct {
	gorm.Model
	OJ        uint
	UID       uint
	Title     string
	Url       string
	Total     uint
	ProblemID uint
	Accept    uint
}

type PracticeClassRecord struct {
	gorm.Model
	Practice uint
	Class    uint
}

type PracticeTagRecord struct {
	gorm.Model
	Practice uint
	Tag      uint
}

/*
  ID: number;
  OJ: number;
  Title: string;
  URL: string;
  Total: number;
  AcRate: number;
  Tags: tag[];
  Class: number[];
*/
type PracticeResponse struct {
	ID     uint
	OJ     uint
	Title  string
	Url    string
	Total  uint
	AcRate float32
	Tags   []Tag
	Class  []ClassInfo
}

func CreatePractice(user *User, OJ uint, ID uint) (*Practice, error) {
	resp, err := util.GetProblem(OJ, ID)
	if err != nil {
		return nil, err
	}
	var practice Practice
	practice.Title = resp.Title
	practice.OJ = resp.OJ
	practice.Url = resp.Url
	practice.UID = user.ID
	practice.ProblemID = resp.ProblemID
	db := GetOpenConnection()
	defer db.Close()
	err = db.Save(&practice).First(&practice).Error
	return &practice, err
}

func GetPracticeResponseByID(ID uint) (*PracticeResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	var practice Practice
	err := db.Find(&practice, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var tagsRecord []PracticeTagRecord
	err = db.Find(&tagsRecord, "practice = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var classRecord []PracticeClassRecord
	err = db.Find(&classRecord, "practice = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var tags []Tag
	var class []ClassInfo
	for _, v := range tagsRecord {
		var r Tag
		err := db.Find(&r, "id = ?", v.Tag).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, r)
	}
	for _, v := range classRecord {
		var r ClassInfo
		err := db.Find(&r, "id = ?", v.Class).Error
		if err != nil {
			return nil, err
		}
		class = append(class, r)
	}
	ret := PracticeResponse{
		ID:    practice.ID,
		OJ:    practice.OJ,
		Title: practice.Title,
		Url:   practice.Url,
		Total: practice.Total,
		Tags:  tags,
		Class: class,
	}
	if practice.Total == 0 {
		ret.AcRate = 0.0
	} else {
		ret.AcRate = (float32(practice.Accept) / float32(practice.Total))
	}
	return &ret, nil
}

func GetPracticeList(user *User, offset uint, limit uint) ([]PracticeResponse, error) {
	var ret []PracticeResponse
	var err error
	db := GetOpenConnection()
	defer db.Close()
	if user.Type == 1 {
		var practice []Practice
		err = db.Limit(limit).Offset(offset).Find(&practice, "uid = ?", user.ID).Error
		for _, i := range practice {
			var tmp *PracticeResponse
			tmp, err = GetPracticeResponseByID(i.ID)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *tmp)
		}
	} else {
		var record []PracticeClassRecord
		err = db.Limit(limit).Offset(offset).Find(&record, "class = ?", user.Class).Error
		for _, i := range record {
			var tmp *PracticeResponse
			tmp, err = GetPracticeResponseByID(i.Practice)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *tmp)
		}
	}
	return ret, err
}

func PatchPracticeClass(user *User, id uint, class []uint) (*PracticeResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	var practice Practice
	err := tx.Find(&practice, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if user == nil || user.ID != practice.UID {
		return nil, fmt.Errorf("Unauthorized operation.")
	}
	err = tx.Unscoped().Delete(PracticeClassRecord{}, "practice = ?", practice.ID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, i := range class {
		var record PracticeClassRecord
		record.Class = i
		record.Practice = id
		err := tx.Create(&record).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return GetPracticeResponseByID(id)
}
