package models

import (
	"log"
	"tcloud-api/src/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDataBase() {
	db, err := gorm.Open("mysql", util.GetDefaultDSN())
	if err != nil {
		log.Fatal("Error init database.\n", err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Directory{})
	db.AutoMigrate(&FileMeta{})
	db.AutoMigrate(&ShareDirectory{})
	db.AutoMigrate(&ShareRecord{})
	db.AutoMigrate(&ClassInfo{})
	db.AutoMigrate(&ClassRecord{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Course{})
	db.AutoMigrate(&CourseClassRecord{})
	db.AutoMigrate(&CourseDirectoryRecord{})
	db.AutoMigrate(&CourseTagsRecord{})
	db.AutoMigrate(&Practice{})
	db.AutoMigrate(&PracticeClassRecord{})
	db.AutoMigrate(&PracticeTagRecord{})
	db.AutoMigrate(&Notice{})
	db.AutoMigrate(&NoticeClassRecord{})

}

func GetOpenConnection() *gorm.DB {
	db, err := gorm.Open("mysql", util.GetDefaultDSN())
	if err != nil {
		log.Fatal("Error init database.\n", err)
	}
	return db
}
