package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

/*
 * status
 * 0: create by login (missing nickname)
 * 1: create by register
 * 2: should confirm email
 * 3: status ok
 * 4: blocked
 * 5: delete
 */
const (
	UserStatusUnregister = iota
	UserStatusCreateByLogin
	UserStatusCreateByRegister
	UserStatusOK
	UserStatusBlocked
	UserStatusDelete
)

type User struct {
	gorm.Model
	Nickname string `gorm:"type:varchar(32)"`
	Email    string `gorm:"type:varchar(100);UNIQUE;unique_index"`
	Password string `gorm:"type:char(32)" json:"-"`
	Class    uint
	Type     int
	Status   int
	DiskRoot uint
}

func SearchUserByLoginForm(username, password string, Type int) (*User, bool) {
	db := GetOpenConnection()
	defer db.Close()
	user := User{
		Email: username,
	}
	db.Where(&user).First(&user)
	if user.Password != password || user.Type != Type {
		return &user, false
	}
	return &user, true
}

func SearchUserByEmail(email string) (*User, error) {
	db := GetOpenConnection()
	defer db.Close()
	var user User
	err := db.Where("email=?", email).First(&user).Error
	return &user, err
}

func InsertUser(user *User) error {
	db := GetOpenConnection()
	defer db.Close()
	tx := db.Begin()
	err := tx.Create(user).First(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	dir := Directory{
		UID:         user.ID,
		Name:        "/",
		IsDirectory: true,
	}
	err = tx.Create(&dir).First(&dir).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	user.DiskRoot = dir.ID
	err = tx.Save(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func UpdateUser(user *User) error {
	db := GetOpenConnection()
	defer db.Close()
	err := db.Save(user).Error
	return err
}
func GetUserByID(ID uint) (*User, error) {
	db := GetOpenConnection()
	defer db.Close()
	var u User
	err := db.Find(&u, "id = ?", ID).Error
	return &u, err
}

func GetUsersByID(ID []uint) ([]User, error) {
	db := GetOpenConnection()
	defer db.Close()
	var r []User
	err := db.Find(&r, "id IN (?)", ID).Error
	return r, err
}

func UpdateUsers(users []User) ([]User, error) {
	db := GetOpenConnection()
	defer db.Close()
	var ret []User
	for _, u := range users {
		var user User
		err := db.Save(&u).First(&user).Error
		if err != nil {
			return ret, err
		}
		ret = append(ret, user)
	}
	return ret, nil
}

func (user *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(user)
}
func (user *User) UnMarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, user)
}

func GetUserList(user *User, offset uint, limit uint) ([]User, error) {
	var ret []User
	/*
		if user.Type != 1 {
			return ret, fmt.Errorf("%s", "invalid operation")
		}*/
	db := GetOpenConnection()
	defer db.Close()
	err := db.Offset(offset).Limit(limit).Find(&ret, "type = ?", 0).Error
	return ret, err
}
