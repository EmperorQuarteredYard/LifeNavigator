package models

type Kanban struct {
	Projects []Project
}

var kanbanStatus = map[int]string{
	1:  "planning",
	2:  "ongoing",
	3:  "done",
	16: "finished",  //在一段时间后软删除
	17: "rework",    //计划中的变种
	18: "reworking", //进行中的变种
	19: "reworked",  //已完成的变种
}
