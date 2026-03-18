package kanbanStatus

const (
	StatusPlanning  = 0
	StatusPending   = 1
	StatusOngoing   = 2
	StatusDone      = 3
	StatusFinished  = 16
	StatusRework    = 17
	StatusReworking = 18
	StatusReworked  = 19
)

var (
	// 正向映射：code -> name
	statusNameMap = map[int]string{
		StatusPlanning:  "planning",
		StatusPending:   "pending",
		StatusOngoing:   "ongoing",
		StatusDone:      "done",
		StatusFinished:  "finished",
		StatusRework:    "rework",
		StatusReworking: "reworking",
		StatusReworked:  "reworked",
	}
	// 反向映射：name -> code
	statusCodeMap map[string]int
)

func init() {
	// 构建反向映射
	statusCodeMap = make(map[string]int, len(statusNameMap))
	for code, name := range statusNameMap {
		statusCodeMap[name] = code
	}
}

// GetStatusName 根据状态码返回名称，未知码返回空字符串和 false
func GetStatusName(code int) (string, bool) {
	name, ok := statusNameMap[code]
	return name, ok
}

// GetStatusCode 根据名称返回状态码，未知名称返回0和false
func GetStatusCode(name string) (int, bool) {
	code, ok := statusCodeMap[name]
	return code, ok
}
