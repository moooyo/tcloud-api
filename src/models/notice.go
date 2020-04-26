package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Notice struct {
	gorm.Model
	Title       string
	Description string
	Level       uint
	Time        int64
	UID         uint
}
type NoticeClassRecord struct {
	gorm.Model
	Notice uint
	Class  uint
}

type NoticeResponse struct {
	ID          uint
	Title       string
	Description string
	Level       uint
	Time        int64
	UID         uint
	Class       []ClassInfo
}

/*
interface courseNoticeRecord {
  ID: number;
  Title: string;
  Description: string;
  Level: number;
  Time: number;
  Class: number[];
}
*/

func GetNoticeResponseByID(ID uint) (*NoticeResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	var notice Notice
	err := db.Find(&notice, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var classRecord []NoticeClassRecord
	err = db.Find(&classRecord, "notice = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var class []ClassInfo
	for _, v := range classRecord {
		var r ClassInfo
		err := db.Find(&r, "id = ?", v.Class).Error
		if err != nil {
			return nil, err
		}
		class = append(class, r)
	}
	ret := NoticeResponse{
		ID:          ID,
		Title:       notice.Title,
		Description: notice.Description,
		Level:       notice.Level,
		Time:        notice.Time,
		UID:         notice.UID,
		Class:       class,
	}
	return &ret, nil
}
func GetNoticeList(user *User, offset uint, limit uint) ([]NoticeResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	var notice []Notice
	var ret []NoticeResponse
	if user.Type == 1 {
		err := db.Offset(offset).Limit(limit).Find(&notice, "uid = ?", user.ID).Error
		if err != nil {
			return ret, err
		}
		for _, x := range notice {
			v, err := GetNoticeResponseByID(x.ID)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *v)
		}
	} else {
		var record []NoticeClassRecord
		err := db.Offset(offset).Limit(limit).Find(&record, "class = ?", user.Class).Error
		if err != nil {
			return ret, err
		}
		for _, x := range record {
			v, err := GetNoticeResponseByID(x.Notice)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *v)
		}
	}
	return ret, nil
}

func CreateNotice(user *User, title string, desc string, level uint) (*Notice, error) {
	db := GetOpenConnection()
	defer db.Close()
	var notice Notice
	notice.Level = level
	notice.Description = desc
	now := time.Now()
	time := now.Unix()
	notice.Time = time
	notice.Title = title
	notice.UID = user.ID
	err := db.Save(&notice).Error
	return &notice, err
}

func PatchNoticeClass(user *User, id uint, class []uint) (*NoticeResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	var notice Notice
	err := tx.Find(&notice, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if user == nil || user.ID != notice.UID {
		return nil, fmt.Errorf("Unauthorized operation.")
	}
	err = tx.Unscoped().Delete(NoticeClassRecord{}, "notice = ?", notice.ID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, i := range class {
		var record NoticeClassRecord
		record.Class = i
		record.Notice = id
		err := tx.Create(&record).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return GetNoticeResponseByID(id)
}
