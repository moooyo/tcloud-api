package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	FileTypeUnset = iota
	FileTypeOther
	FileTypePdf
	FileTypeCsv
	FileTypeXlsx
	FileTypeDocx
	FileTypeMp4
	FileTypeWebm
	FileTypeMp3
	FileTypePng
	FileTypeJpeg
	FileTypeGif
	FileTypeBmp
	FileTypeJpg

	// End
	FileTypeEnd
)

const (
	FileListTypeUnset = iota
	FileListTypeOther
	FileListTypeMusic
	FileListTypeVideo
	FileListTypeDoc
	FileListTypeImage
)

func FileListType2Array(Type uint) []uint {
	switch Type {
	case FileListTypeDoc:
		return []uint{FileTypePdf, FileTypeDocx, FileTypeXlsx, FileTypeCsv}
	case FileListTypeImage:
		return []uint{FileTypePng, FileTypeJpeg, FileTypeGif, FileTypeBmp, FileTypeJpg}
	case FileListTypeMusic:
		return []uint{FileTypeMp3}
	case FileListTypeVideo:
		return []uint{FileTypeMp4, FileTypeWebm}
	case FileListTypeOther:
		return []uint{FileTypeEnd}
	default:
		return []uint{FileTypeEnd}
	}
}

func FileName2FileType(name string) uint {
	strs := strings.Split(name, ".")
	if len(strs) <= 1 {
		return FileTypeOther
	}
	switch strs[len(strs)-1] {
	case "pdf":
		return FileTypePdf
	case "csv":
		return FileTypeCsv
	case "xlsx":
		return FileTypeXlsx
	case "docx":
		return FileTypeDocx
	case "mp4":
		return FileTypeMp4
	case "webm":
		return FileTypeWebm
	case "mp3":
		return FileTypeMp3
	case "png":
		return FileTypePng
	case "jpeg":
		return FileTypeJpeg
	case "jpg":
		return FileTypeJpg
	case "bmp":
		return FileTypeBmp
	case "gif":
		return FileTypeGif
	default:
		return FileTypeOther
	}
}

func FileType2Str(fileType uint) string {
	switch fileType {
	case FileTypePdf:
		return "pdf"
	case FileTypeCsv:
		return "csv"
	case FileTypeDocx:
		return "docx"
	case FileTypeMp3:
		return "mp3"
	case FileTypeMp4:
		return "mp4"
	case FileTypeXlsx:
		return "xlsx"
	case FileTypeWebm:
		return "webm"
	case FileTypePng:
		return "png"
	case FileTypeJpeg:
		return "jpeg"
	case FileTypeGif:
		return "gif"
	case FileTypeBmp:
		return "bmp"
	case FileTypeJpg:
		return "jpg"
	default:
		return "other"
	}
}

type Directory struct {
	gorm.Model
	UID         uint
	Name        string `gorm:"type:varchar(128)"`
	IsDirectory bool   `gorm:"not null"`
	PreIndex    uint
	MetaID      uint `gorm:"default:'0'" json:"-"`
	Size        uint `gorm:"default:'0'"`
	Type        uint `gorm:"default:'0'"`
}

func SearchDirectoryByID(ID uint) (*Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var d Directory
	err := db.First(&d, "ID=?", ID).Error
	return &d, err
}

func InsertDirectory(d *Directory) error {
	db := GetOpenConnection()
	defer db.Close()
	err := db.Create(d).First(d).Error
	return err
}

func UpdateDirectory(d *Directory) error {
	db := GetOpenConnection()
	defer db.Close()
	err := db.Update(d).First(d).Error
	return err
}

func SearchFileListByPreIndex(ID uint, offset, limit uint) ([]Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []Directory
	err := db.Offset(offset).Limit(limit).Find(&ret, "pre_index=?", ID).Error
	return ret, err
}
func SearchFileListByPreIndexWithoutLimit(ID uint) ([]Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []Directory
	err := db.Find(&ret, "pre_index=?", ID).Error
	return ret, err
}

func CreateDirectoryByPathID(ID uint, UID uint, Name string) (*Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	directory := Directory{
		Name:        Name,
		PreIndex:    ID,
		IsDirectory: true,
		UID:         UID,
	}
	err := db.Create(&directory).First(&directory).Error
	return &directory, err
}

func ChangeDirectoryName(user *User, ID uint, Name string) error {
	var directory Directory
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	err := tx.Where("id=?", ID).First(&directory).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if directory.UID != user.ID {
		tx.Rollback()
		return fmt.Errorf("%s", "unauthorized")
	}
	directory.Name = Name
	err = tx.Save(&directory).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func DeleteDirectories(id []uint, uid uint) error {
	db := GetOpenConnection()
	tx := db.Begin()
	defer db.Close()
	for _, i := range id {
		var dir Directory
		dir.ID = i
		err := tx.Find(&dir).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if dir.UID != uid {
			return fmt.Errorf("%s", "unauthorize")
		}
		err = tx.Delete(&dir).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func SearchFileListByType(user *User, Type, offset, limit uint) ([]Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var data []Directory
	typeArr := FileListType2Array(Type)
	var err error
	if Type != FileListTypeOther {
		err = db.Offset(offset).Limit(limit).Find(&data, "uid = ? AND type IN (?) AND is_directory = ?", user.ID, typeArr, false).Error
	} else {
		err = db.Offset(offset).Limit(limit).Find(&data, "uid = ? AND is_directory = ? AND (type >= ? OR type = ?)", user.ID, false, typeArr[0], FileListTypeOther).Error
	}
	return data, err
}

func DeleteFileCompletely(user *User, id []uint) ([]uint, error) {
	db := GetOpenConnection()
	defer db.Close()
	var data []Directory
	err := db.Unscoped().Find(&data, "id in (?) AND uid = ?", id, user.ID).Error
	var ret []uint
	if err != nil {
		return ret, err
	}
	for _, x := range data {
		var tmp Directory
		err := db.Unscoped().Find(&tmp, "id = ?", x.ID).Delete(&tmp).Error
		if err != nil {
			return ret, err
		}
		ret = append(ret, x.ID)
	}

	return ret, err
}

func GetDeletedFileList(user *User, offset uint, limit uint) ([]Directory, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []Directory
	err := db.Unscoped().Offset(offset).Limit(limit).Find(&ret, "uid = ? AND deleted_at <> ?", user.ID, "").Error
	return ret, err
}

func RestoreTrash(user *User, ids []uint) ([]uint, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []uint
	for _, id := range ids {
		var dir Directory
		err := db.Unscoped().Find(&dir, "id = ?", id).Error
		if err != nil {
			return ret, err
		}
		dir.DeletedAt = nil
		err = db.Unscoped().Model(&dir).Update("deleted_at", nil).Error
		if err != nil {
			return ret, err
		}
		ret = append(ret, id)
	}
	return ret, nil
}
