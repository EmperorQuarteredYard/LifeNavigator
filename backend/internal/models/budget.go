package models

type TaskPayment struct {
	ID       uint64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TaskID   uint64  `gorm:"not null;index" json:"task_id"` // 外键
	BudgetID uint64  `gorm:"not null;index" json:"budget_id"`
	Amount   float64 `gorm:"type:decimal(10,2);default:0" json:"used"` // 已用金额
}
type ProjectBudget struct {
	ID        uint64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProjectID uint64  `gorm:"not null;index" json:"project_id"`           // 外键
	AccountID uint64  `json:"account_id"`                                 // 扣除哪个账户
	Unit      string  `gorm:"unit:varchar(50);default:'CNY'" json:"unit"` // 预算单位：CNY/USD/token/hour
	Budget    float64 `gorm:"type:decimal(10,2);default:0" json:"budget"` // 预算金额
	Used      float64 `gorm:"type:decimal(10,2);default:0" json:"used"`   // 已用金额
}

func (m *ProjectBudget) SetID(id uint64) {
	m.ID = id
}
func (m *ProjectBudget) GetID() uint64 {
	return m.ID
}
func (m *TaskPayment) SetID(id uint64) {
	m.ID = id
}
func (m *TaskPayment) GetID() uint64 {
	return m.ID
}
