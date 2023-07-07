package global

import "github.com/ncuhome/GeniusAuthoritarian/pkg/departments"

var Departments = []string{
	departments.UDev,
	departments.UPm,
	departments.UGame,
	departments.UOp,
	departments.UAdm,
	departments.UDes,
	departments.UCe,
	departments.USenior,
}

var DepartmentRelation map[string]uint
