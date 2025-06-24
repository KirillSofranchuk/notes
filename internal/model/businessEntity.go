package model

type BusinessEntity interface {
	SetId(id int)
	GetId() int
	SetTimestamp()
}
