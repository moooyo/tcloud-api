package models

import "github.com/jinzhu/gorm"

type FileMeta struct {
	gorm.Model
	RealName string
	Size     uint
}

func InsertFileMeta(meta *FileMeta) error {
	db := GetOpenConnection()
	defer db.Close()
	err := db.Create(meta).First(meta).Error

	return err
}

func SearchFileMetaByID(ID uint) (*FileMeta, error) {
	db := GetOpenConnection()
	defer db.Close()
	var f FileMeta
	err := db.Find(&f, "id=?", ID).Error
	return &f, err
}
