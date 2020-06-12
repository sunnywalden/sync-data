package models

type loginId string
type department string
type title string
type name string
type nickname string

type users struct {
	LoginId loginId
	Depart department
	Title title
	Name name
	NickName nickname
}
