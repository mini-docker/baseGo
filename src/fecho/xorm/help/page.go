package help

import (
	"fecho/utility"
	"fecho/xorm"
)

//分页参数
type PageParams struct {
	Page     int    `json:"pageIndex"`     //页码数
	PageSize int    `json:"pageSize"` //页面数量
	Limit    []int  `json:"-"`        //分页
	OrderBy  string `json:"orderBy"`  //排序字段
	Asc      bool   `json:"asc"`      //排序顺序 默认为false,也就是默认倒序
}

// Make 给xorm 补上分页信息
func (l *PageParams) Make(session *xorm.Session, columns ...string) {

	if l.Page < 1 {
		l.Page = 1
	}
	if l.PageSize < 1 {
		l.PageSize = 50
	} else {
		if l.PageSize >= 1000 {
			l.PageSize = 1000
		}
	}
	offset := (l.Page - 1) * l.PageSize
	var limit = []int{l.PageSize, offset}
	l.Limit = limit

	session.Limit(l.Limit[0], l.Limit[1])
	// 排序
	sort := func(columns ...string) {
		if l.Asc == true {
			session.Asc(columns...)
		} else {
			session.Desc(columns...)
		}
	}
	switch {
	case l.OrderBy != "":
		// 优先级最高:前端传了Order字段,按前端的排序
		l.OrderBy = utility.Camel2Underline(l.OrderBy)
		sort(l.OrderBy)
	case len(columns) > 0:
		// 优先级其次,表中存在并明确表示需要排序的字段,例如update_time
		sort(columns...)
	default:
		// 优先级最次,id,有分页的数据都会有id
		sort("id")
	}
}
