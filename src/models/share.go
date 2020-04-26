package models

import (
	"fmt"
	"tcloud-api/src/util"
	"time"

	"github.com/jinzhu/gorm"
)

type ShareRecord struct {
	gorm.Model
	InternalID           string `gorm:"UNIQUE_INDEX"`
	ShareRootID          uint
	UID                  uint
	NickName             string
	Secret               bool
	Password             string
	ShareName            string
	ShareFileType        uint
	ShareFileIsDirectory bool
	Expired              time.Time
}

type ShareDirectory struct {
	gorm.Model
	ShareID     uint
	UID         uint
	Name        string `gorm:"type:varchar(128)"`
	IsDirectory bool   `gorm:"not null"`
	PreIndex    uint
	MetaID      uint `gorm:"default:'0'" json:"-"`
	Size        uint `gorm:"default:'0'"`
	Type        uint `gorm:"default:'0'"`
	Expires     time.Time
}

func MakeShareDirectoryFromDirectory(dir *Directory, pre uint, expires time.Time) *ShareDirectory {
	return &ShareDirectory{
		Name:        dir.Name,
		IsDirectory: dir.IsDirectory,
		MetaID:      dir.MetaID,
		Size:        dir.Size,
		Type:        dir.Type,
		PreIndex:    pre,
		Expires:     expires,
	}
}

type shareStack struct {
	PreIndex uint
	Data     *Directory
}

/*
 * Create share from some path
 */
const shareCodeLength = 8

func CreateShare(user *User, name string, path []uint, expires time.Time, secret bool) (string, string, error) {
	if len(path) == 0 {
		return "", "", fmt.Errorf("%s", "Path params error.")
	}
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	var shareRoot ShareDirectory
	shareRoot.Expires = expires
	shareRoot.UID = user.ID
	shareRoot.Name = name
	shareRoot.PreIndex = 0
	shareRoot.IsDirectory = true
	shareRoot.ShareID = 0
	err := tx.Create(&shareRoot).First(&shareRoot).Error
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	record := ShareRecord{
		UID:         user.ID,
		ShareRootID: shareRoot.ID,
		Expired:     expires,
		InternalID:  util.GenerateUUID(),
		Secret:      secret,
		NickName:    user.Nickname,
	}

	source := make([]shareStack, 0, 2*len(path))
	var rootFile *Directory = nil

	for _, pid := range path {
		var dir Directory
		err = tx.Where("id = ?", pid).Find(&dir).Error
		if err != nil {
			tx.Rollback()
			return "", "", err
		}
		if dir.UID != user.ID {
			tx.Rollback()
			return "", "", fmt.Errorf("%s", "unauthorize")
		}
		if rootFile == nil {
			record.ShareName = dir.Name
			record.ShareFileType = dir.Type
			record.ShareFileIsDirectory = dir.IsDirectory
			rootFile = &dir
		}
		source = append(source, shareStack{
			PreIndex: shareRoot.ID,
			Data:     &dir,
		})
	}
	if secret {
		record.Password = util.GenerateCaptcha(shareCodeLength)
	}
	err = tx.Create(&record).First(&record).Error
	if err != nil {
		tx.Rollback()
		return "", "", err
	}
	for len(source) > 0 {
		data := source[len(source)-1]
		source = source[:len(source)-1]
		share := MakeShareDirectoryFromDirectory(data.Data, data.PreIndex, expires)
		share.ShareID = record.ID
		err := tx.Create(&share).First(&share).Error
		if err != nil {
			tx.Rollback()
			return "", "", err
		}
		if data.Data.IsDirectory {
			var nextDir []Directory
			err := tx.Where("pre_index = ?", data.Data.ID).Find(&nextDir).Error
			if err != nil {
				tx.Rollback()
				return "", "", err
			}
			for _, dir := range nextDir {
				source = append(source, shareStack{
					PreIndex: share.ID,
					Data:     &dir,
				})
			}
		}
	}
	tx.Commit()
	return record.InternalID, record.Password, nil
}

func GetShareRecord(uuid string) (*ShareRecord, error) {
	db := GetOpenConnection()
	defer db.Close()
	var record ShareRecord
	err := db.Where("internal_id = ?", uuid).First(&record).Error
	return &record, err
}

func GetShareDirectoriesByPreIndex(index uint) ([]ShareDirectory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var dirs []ShareDirectory
	var err error
	err = db.Find(&dirs, "pre_index = ?", index).Error
	return dirs, err
}

func GetShareRecordListByUserID(uid, offset, limit uint) ([]ShareRecord, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []ShareRecord
	err := db.Limit(limit).Offset(offset).Find(&ret, "uid = ?", uid).Error
	return ret, err
}

func GetMetaIDByDirectoryAndShare(id uint, share uint) (*FileMeta, string, error) {
	db := GetOpenConnection().Debug()
	defer db.Close()
	var directory ShareDirectory
	err := db.Find(&directory, "share_id = ? AND id = ?", share, id).Error
	if err != nil {
		return nil, "", err
	}
	var meta FileMeta
	err = db.Find(&meta, "id = ?", directory.MetaID).Error
	return &meta, directory.Name, err
}
