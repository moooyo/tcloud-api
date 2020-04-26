package models

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model
	Name string
}

func GetTagByName(Name string) (*Tag, error) {
	db := GetOpenConnection()
	defer db.Close()
	var tag Tag
	err := db.Find(&tag, "name = ?", Name).Error
	return &tag, err
}

func GetTagsByID(ID []uint) ([]Tag, error) {
	db := GetOpenConnection()
	defer db.Close()
	var tags []Tag
	err := db.Find(&tags, "id in (?)", ID).Error
	return tags, err
}

func GetTagsByNameWithCreate(Name string) (*Tag, bool, error) {
	db := GetOpenConnection()
	defer db.Close()
	var tag Tag
	create := false
	err := db.Find(&tag, "Name = ?", Name).Error
	if err != nil || tag.ID == 0 {
		tag.Name = Name
		err = db.Create(&tag).Error
		create = true
	}
	return &tag, create, err
}

func GetTagsList() ([]Tag, error) {
	db := GetOpenConnection()
	defer db.Close()
	var tags []Tag
	err := db.Find(&tags).Error
	return tags, err
}
