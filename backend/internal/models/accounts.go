package models

type Account struct {
	ID         uint64  `gorm:"primary_key;auto_increment"`
	Users      []User  `gorm:"many2many:account_users;"`
	Name       string  `gorm:"type:varchar(50);not null"`
	Type       string  `gorm:"type:varchar(50);not null"`
	Unit       string  `gorm:"type:varchar(50);not null"`
	Balance    float64 `gorm:"type:decimal(10,2);default:0" json:"balance"`
	NetBalance float64 `gorm:"type:decimal(10,2);default:0" json:"net_balance"`
	Version    uint64  `gorm:"not null;default:0"`
}

func (m *Account) SetID(id uint64) {
	m.ID = id
}
func (m *Account) GetID() uint64 {
	return m.ID
}
