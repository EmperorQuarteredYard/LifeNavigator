package models

type Model interface {
	SetID(id uint64)
	GetID() uint64
}
