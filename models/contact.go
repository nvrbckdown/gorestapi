package models

import (
	u "../utils"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id"` //Пользователь, которому принадлежит этот контакт
}

func (contact *Contact) Validate() (map[string]interface{}, bool) {

	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserId <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	return u.Message(true, "Success."), true

}

func (contact *Contact) Create() map[string]interface{} {

	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "Success.")
	resp["contact"] = contact
	return resp

}

func GetContact(id uint) *Contact {

	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		return nil
	}
	return contact

}

func GetContacts(user uint) []*Contact {

	contatcs := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contatcs).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contatcs

}
