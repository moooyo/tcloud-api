package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type TagParams struct {
	ID   int
	Name string
}
type Course struct {
	gorm.Model
	UID         uint
	Name        string
	Description string
	StartTime   uint
	EndTime     uint
}

type CourseClassRecord struct {
	gorm.Model
	Course uint
	Class  uint
}

type CourseTagsRecord struct {
	gorm.Model
	Course uint
	Tag    uint
}

type CourseDirectoryRecord struct {
	gorm.Model
	Course    uint
	Directory uint
	Pre_Index uint
}

type CourseResponse struct {
	ID          uint
	Name        string
	Description string
	StartTime   uint
	EndTime     uint
	Class       []ClassInfo
	FileList    []Directory
	Tags        []Tag
}

func GetCourseResponseByID(ID uint) (*CourseResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	var course Course
	err := db.Find(&course, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var tagsRecord []CourseTagsRecord
	err = db.Find(&tagsRecord, "course = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var classRecord []CourseClassRecord
	err = db.Find(&classRecord, "course = ?", ID).Error
	if err != nil {
		return nil, err
	}
	var directoryRecord []CourseDirectoryRecord
	err = db.Find(&directoryRecord, "course = ? AND pre_index = ?", ID, 0).Error
	if err != nil {
		return nil, err
	}
	var tags []Tag
	var class []ClassInfo
	var dirs []Directory
	for _, v := range tagsRecord {
		var r Tag
		err := db.Find(&r, "id = ?", v.Tag).Error
		if err != nil {
			return nil, err
		}
		tags = append(tags, r)
	}
	for _, v := range directoryRecord {
		var r Directory
		err := db.Find(&r, "id = ?", v.Directory).Error
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, r)
	}
	for _, v := range classRecord {
		var r ClassInfo
		err := db.Find(&r, "id = ?", v.Class).Error
		if err != nil {
			return nil, err
		}
		class = append(class, r)
	}
	ret := CourseResponse{
		ID:          course.ID,
		StartTime:   course.StartTime,
		EndTime:     course.EndTime,
		FileList:    dirs,
		Tags:        tags,
		Class:       class,
		Name:        course.Name,
		Description: course.Description,
	}
	return &ret, nil
}

func GetCourseList(user *User, offset uint, limit uint) ([]CourseResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []CourseResponse
	var course []Course
	if user.Type == 1 {
		err := db.Offset(offset).Limit(limit).Find(&course, "uid = ?", user.ID).Error
		if err != nil {
			return ret, err
		}
		for _, x := range course {
			resp, err := GetCourseResponseByID(x.ID)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *resp)
		}
	} else {
		var record []CourseClassRecord
		err := db.Offset(offset).Limit(limit).Find(&record, "class = ?", user.Class).Error
		if err != nil {
			return ret, err
		}
		for _, x := range record {
			resp, err := GetCourseResponseByID(x.Course)
			if err != nil {
				return ret, err
			}
			ret = append(ret, *resp)
		}
	}
	return ret, nil
}
func CourseWalk(path uint, tx *gorm.DB, id uint) error {
	var dir []Directory
	err := tx.Find(&dir, "pre_index = ?", path).Error
	if err != nil {
		return err
	} else {
		for _, d := range dir {
			var record CourseDirectoryRecord
			record.Course = id
			record.Directory = d.ID
			record.Pre_Index = path
			err = tx.Save(&record).First(&record).Error
			if err != nil {
				return err
			}
			if d.IsDirectory {
				err = CourseWalk(d.ID, tx, id)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CreateCourse(course *Course, tags []TagParams, files []uint, class []uint) (*Course, error) {
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	var ret Course
	err := tx.Create(course).First(&ret).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for i := range tags {
		if tags[i].ID <= 0 {
			tag, _, err := GetTagsByNameWithCreate(tags[i].Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			var record CourseTagsRecord
			record.Course = ret.ID
			record.Tag = tag.ID
			err = tx.Create(&record).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	for _, i := range files {
		var dir Directory
		err := tx.Find(&dir, "id = ?", i).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		var record CourseDirectoryRecord
		record.Course = ret.ID
		record.Directory = dir.ID
		record.Pre_Index = 0
		err = tx.Save(&record).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if dir.IsDirectory {
			err = CourseWalk(dir.ID, tx, ret.ID)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	for _, i := range class {
		var info ClassInfo
		err := tx.Find(&info, "id = ?", i).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		var record CourseClassRecord
		record.Course = ret.ID
		record.Class = info.ID
		err = tx.Save(&record).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return &ret, nil
}

func PatchCourseClass(user *User, id uint, class []uint) (*CourseResponse, error) {
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	var course Course
	err := tx.Find(&course, "id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if user == nil || user.ID != course.UID {
		return nil, fmt.Errorf("Unauthorized operation.")
	}
	err = tx.Unscoped().Delete(CourseClassRecord{}, "course = ?", course.ID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, i := range class {
		var record CourseClassRecord
		record.Class = i
		record.Course = id
		err := tx.Create(&record).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return GetCourseResponseByID(id)
}

func GetCourseDirectoryByCourseAndPath(user *User, course uint, path uint) ([]Directory, error) {
	var ret []Directory
	var record []CourseDirectoryRecord
	var classRecord []CourseClassRecord
	db := GetOpenConnection()
	defer db.Close()
	err := db.Find(&classRecord, "course = ?", course).Error
	if err != nil {
		return ret, err
	}
	f := false
	for _, x := range classRecord {
		if x.Class == user.Class {
			f = true
			break
		}
	}
	if !f {
		return ret, fmt.Errorf("%s", "No authorized.")
	}
	err = db.Find(&record, "course = ? AND pre_index = ?", course, path).Error
	if err != nil {
		return ret, err
	}
	for _, v := range record {
		var tmp Directory
		err = db.Find(&tmp, "id = ?", v.Directory).Error
		if err != nil {
			return ret, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func GetShareCourseFileMeta(user *User, cid uint, dir uint) (*FileMeta, string, error) {
	if user == nil {
		return nil, "", fmt.Errorf("%s", "Please login first.")
	}
	db := GetOpenConnection()
	defer db.Close()
	var course Course
	err := db.Find(&course, "id = ?", cid).Error
	if err != nil {
		return nil, "", err
	}
	var classRecord []CourseClassRecord
	err = db.Find(&classRecord, "course = ?", cid).Error
	if err != nil {
		return nil, "", err
	}
	find := false
	for _, x := range classRecord {
		if x.Class == user.Class {
			find = true
		}
	}
	if !find {
		return nil, "", fmt.Errorf("%s", "Unauthorized.")
	}
	var courseDirectoryRecord CourseDirectoryRecord
	err = db.Find(&courseDirectoryRecord, "course = ? AND directory = ?", cid, dir).Error
	if err != nil {
		return nil, "", err
	}
	var directory Directory
	err = db.Find(&directory, "id = ?", courseDirectoryRecord.Directory).Error
	if err != nil {
		return nil, "", err
	}
	var meta FileMeta
	err = db.Find(&meta, "id = ?", directory.MetaID).Error
	if err != nil {
		return nil, "", err
	}
	return &meta, directory.Name, nil
}
