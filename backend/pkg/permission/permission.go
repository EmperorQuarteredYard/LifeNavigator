package permission

import (
	"errors"
	"fmt"
	"strings"
)

// Role 表示权限角色
type Role uint8

const (
	RoleOwner Role = iota
	RoleWorkmate
	RoleViewer
	RoleGuest
)

// String 返回角色的可读名称
func (r Role) String() string {
	switch r {
	case RoleOwner:
		return "Owner"
	case RoleWorkmate:
		return "Workmate"
	case RoleViewer:
		return "Viewer"
	case RoleGuest:
		return "Guest"
	default:
		return "Unknown"
	}
}

// Operation 表示具体操作
type Operation uint8

const (
	OpCreate Operation = iota
	OpRead
	OpDelete
	OpUpdate
	OpBranch
)

// String 返回操作的可读名称
func (o Operation) String() string {
	switch o {
	case OpCreate:
		return "create"
	case OpRead:
		return "read"
	case OpDelete:
		return "delete"
	case OpUpdate:
		return "update"
	case OpBranch:
		return "branch"
	default:
		return "unknown"
	}
}

// 位偏移计算：每个角色占用7位，偏移 = 角色索引 * 7 + 操作索引
func bitOffset(role Role, op Operation) uint {
	return uint(role)*7 + uint(op)
}

// PermissionSet 封装 uint32 权限值
type PermissionSet uint32

// New 创建一个空的权限集
func New() PermissionSet {
	return 0
}

// FromUint32 从已有的 uint32 值创建权限集
func FromUint32(v uint32) PermissionSet {
	return PermissionSet(v)
}

// ParsePermissionSet 从字符串解析权限集，字符串格式应与 PermissionSet.String() 的输出兼容。
// 支持额外的空格，不区分大小写。若字符串为 "none" 则返回空权限集。
func ParsePermissionSet(s string) (PermissionSet, error) {
	ps := New()
	s = strings.TrimSpace(s)
	if s == "" || s == "none" {
		return ps, nil
	}

	// 按 "; " 分割各个角色部分，但允许分号前后可能有空格
	parts := strings.Split(s, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 角色名与操作列表由 ": " 分隔
		roleAndOps := strings.SplitN(part, ":", 2)
		if len(roleAndOps) != 2 {
			return ps, errors.New("invalid permission format: missing colon separator in '" + part + "'")
		}
		roleName := strings.TrimSpace(roleAndOps[0])
		opsStr := strings.TrimSpace(roleAndOps[1])
		if opsStr == "" {
			continue
		}

		// 解析角色
		role, err := ParseRole(roleName)
		if err != nil {
			return ps, err
		}

		// 按 ", " 分割操作名，允许逗号后可能有空格
		opNames := strings.Split(opsStr, ",")
		for _, opName := range opNames {
			opName = strings.TrimSpace(opName)
			if opName == "" {
				continue
			}
			op, err := ParseOperation(opName)
			if err != nil {
				return ps, err
			}
			ps.Set(role, op)
		}
	}
	return ps, nil
}

// Uint32 返回底层 uint32 值
func (p *PermissionSet) Uint32() uint32 {
	return uint32(*p)
}

// Has 检查指定角色是否拥有指定操作权限
func (p *PermissionSet) Has(role Role, op Operation) bool {
	offset := bitOffset(role, op)
	mask := uint32(1) << offset
	return (uint32(*p) & mask) != 0
}

// Set 授予指定角色指定操作权限
func (p *PermissionSet) Set(role Role, op Operation) {
	offset := bitOffset(role, op)
	mask := uint32(1) << offset
	*p = PermissionSet(uint32(*p) | mask)
}

// Clear 移除指定角色指定操作权限
func (p *PermissionSet) Clear(role Role, op Operation) {
	offset := bitOffset(role, op)
	mask := uint32(1) << offset
	*p = PermissionSet(uint32(*p) & ^mask)
}

// Toggle 翻转指定角色指定操作权限
func (p *PermissionSet) Toggle(role Role, op Operation) {
	if p.Has(role, op) {
		p.Clear(role, op)
	} else {
		p.Set(role, op)
	}
}

// SetRoleAll 授予某个角色的所有操作权限
func (p *PermissionSet) SetRoleAll(role Role) {
	for op := OpCreate; op <= OpBranch; op++ {
		p.Set(role, op)
	}
}

// ClearRoleAll 移除某个角色的所有操作权限
func (p *PermissionSet) ClearRoleAll(role Role) {
	for op := OpCreate; op <= OpBranch; op++ {
		p.Clear(role, op)
	}
}

// RolePermissions 返回指定角色的权限掩码（7位）
func (p *PermissionSet) RolePermissions(role Role) uint8 {
	return uint8((*p) >> (role * 7) & 0x7f)
}

// String 返回可读的权限表示，格式如 "User: create,read; Workmate: delete"
func (p *PermissionSet) String() string {
	var result string
	roles := []Role{RoleOwner, RoleWorkmate, RoleViewer, RoleGuest}
	for _, role := range roles {
		var ops []string
		for op := OpCreate; op <= OpBranch; op++ {
			if p.Has(role, op) {
				ops = append(ops, op.String())
			}
		}
		if len(ops) > 0 {
			if result != "" {
				result += "; "
			}
			result += fmt.Sprintf("%s: %s", role.String(), join(ops, ", "))
		}
	}
	if result == "" {
		return "none"
	}
	return result
}

// join 是简单的字符串连接辅助函数
func join(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	var b = make([]byte, n)
	bp := copy(b, elems[0])
	for _, s := range elems[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}

// ParseRole 将字符串转换为 Role，不区分大小写
func ParseRole(s string) (Role, error) {
	switch s {
	case "Owner", "owner":
		return RoleOwner, nil
	case "Workmate", "workmate":
		return RoleWorkmate, nil
	case "Viewer", "viewer":
		return RoleViewer, nil
	case "Guest", "guest":
		return RoleGuest, nil
	default:
		return 0, errors.New("unknown role: " + s)
	}
}

// ParseOperation 将字符串转换为 Operation，不区分大小写
func ParseOperation(s string) (Operation, error) {
	switch s {
	case "create", "Create":
		return OpCreate, nil
	case "read", "Read":
		return OpRead, nil
	case "delete", "Delete":
		return OpDelete, nil
	case "update", "Update":
		return OpUpdate, nil
	case "branch", "Branch":
		return OpBranch, nil
	default:
		return 0, errors.New("unknown operation: " + s)
	}
}
