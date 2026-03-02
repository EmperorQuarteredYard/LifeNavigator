package models

type TaskBudget struct {
	ID     uint64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TaskID uint64  `gorm:"not null;index" json:"task_id"`              // 外键
	Type   string  `gorm:"type:varchar(50);not null" json:"type"`      // 预算类型（如 "资金", "工时"）
	Budget float64 `gorm:"type:decimal(10,2);default:0" json:"budget"` // 预算金额
	Used   float64 `gorm:"type:decimal(10,2);default:0" json:"used"`   // 已用金额
}
type ProjectBudget struct {
	ID        uint64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProjectID uint64  `gorm:"not null;index" json:"project_id"` // 外键
	Type      string  `gorm:"type:varchar(50);not null" json:"type"`
	Budget    float64 `gorm:"type:decimal(10,2);default:0" json:"budget"`
	Used      float64 `gorm:"type:decimal(10,2);default:0" json:"used"`
}

func MergeBudgetItems(items []TaskBudget) []TaskBudget {
	// 使用 map 按类型聚合金额和已用金额
	merged := make(map[string]*TaskBudget)

	for _, item := range items {
		// 如果类型已存在，累加金额
		if p, ok := merged[item.Type]; ok {
			p.Budget += item.Budget
			p.Used += item.Used
		} else {
			// 否则创建新项（必须复制，不能直接取 item 的地址）
			newItem := TaskBudget{
				Type:   item.Type,
				Budget: item.Budget,
				Used:   item.Used,
			}
			merged[item.Type] = &newItem
		}
	}

	// 将 map 中的值转换为切片
	result := make([]TaskBudget, 0, len(merged))
	for _, itemPtr := range merged {
		result = append(result, *itemPtr)
	}
	return result
}
