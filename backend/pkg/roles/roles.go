package roles

const (
	Administrator = "administrators" //管理员
	User          = "users"          //用户
	Guest         = "guests"         //游客
	Developer     = "developers"     //开发者
	Advocator     = "advocators"     //支持者
	Default       = Guest            //默认为游客
)

var privilegeValue = map[string]int{
	Developer:     10000,
	Administrator: 1000,
	Advocator:     100,
	User:          10,
	Guest:         0,
}

func GetPrivilegeValue(name string) int {
	value, ok := privilegeValue[name]
	if ok {
		return value
	}
	return -1
}
