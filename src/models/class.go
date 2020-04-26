package models

import (
	"fmt"
	"tcloud-api/src/util"

	"github.com/jinzhu/gorm"
)

type ClassInfo struct {
	gorm.Model
	Name  string
	Code  string
	Total uint
}

type ClassRecord struct {
	gorm.Model
	UID       uint
	privilage uint
}

const ClassCodeLength = 6

func CreateClassInfo(user *User, name string) (*ClassInfo, error) {
	if user == nil || user.Type == 0 {
		return nil, fmt.Errorf("%s", "Invalid operation")
	}
	db := GetOpenConnection()
	defer db.Close()
	info := ClassInfo{
		Name:  name,
		Code:  util.GenerateCaptcha(ClassCodeLength),
		Total: 0,
	}
	err := db.Create(&info).Error
	return &info, err
}

func GetClassInfoList(user *User) ([]ClassInfo, error) {
	db := GetOpenConnection()
	defer db.Close()
	var info []ClassInfo
	err := db.Find(&info).Error
	if err == nil && user != nil && user.Type != 1 {
		for _, f := range info {
			f.Code = ""
		}
	}
	return info, err
}

func GetClassInfoByID(user *User, id uint) (*ClassInfo, error) {
	if user == nil || user.Type == 0 {
		return nil, fmt.Errorf("%s", "invalid operation.")
	}
	db := GetOpenConnection()
	defer db.Close()
	var info ClassInfo
	err := db.Find(&info, "id = ?", id).Error
	return &info, err
}

func UpdateClassInfo(user *User, info *ClassInfo) (*ClassInfo, error) {
	if user == nil || user.Type == 0 {
		return nil, fmt.Errorf("%s", "invalid operation.")
	}
	db := GetOpenConnection()
	defer db.Close()
	err := db.Save(info).First(info).Error
	return info, err
}

func DeleteClassInfo(user *User, id uint) error {
	if user == nil || user.Type == 0 {
		return fmt.Errorf("%s", "invalid operation.")
	}
	db := GetOpenConnection()
	defer db.Close()
	var info ClassInfo
	err := db.Find(&info, "id = ?", id).Delete(&info).Error
	return err
}

func GetClassRecordByUID(user *User, uid uint) (*ClassRecord, error) {
	if user.Type != 1 && uid != user.ID {
		return nil, fmt.Errorf("%s", "Invalid operation.")
	}
	db := GetOpenConnection()
	defer db.Close()
	var record ClassRecord
	err := db.Find(&record, "uid = ?", uid).Error
	return &record, err
}
