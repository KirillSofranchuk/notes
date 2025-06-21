package model

type BusinessEntity interface {
	GetInfo() string
	SetId(id int)
	GetId() int
	SetTimestamp()
}
