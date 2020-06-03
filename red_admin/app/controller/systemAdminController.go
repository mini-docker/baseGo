package controller

import (
	"baseGo/model"
	"baseGo/model/code"
	"baseGo/red_admin/app/controller/common"
	"baseGo/red_admin/app/middleware"
	"baseGo/red_admin/app/middleware/validate"
	"baseGo/red_admin/app/server"
	"baseGo/red_admin/app/services"
)

// 系统管理员控制器
type SystemAdminController struct {
	QueryAdminReq struct {
		StartTime int    `json:"startTime"`
		EndTime   int    `json:"endTime"`
		RoleId    int    `json:"roleId"`
		IsOnline  int    `json:"isOnline"`
		Account   string `json:"account"`
		PageIndex int    `json:"pageIndex"` // 页码
		PageSize  int    `json:"pageSize"`  // 每页条数
	}

	AddAdminReq struct {
		Account         string `json:"account" valid:"Must;ErrorCode(3026)`  // 账号
		Password        string `json:"password" valid:"Must;ErrorCode(3027)` // 密码
		ConfirmPassword string `json:"confirmPassword"`                      // 确认密码
		RoleId          int    `json:"roleId" valid:"Must;ErrorCode(3029)`   // 角色id
	}

	QueryAdminOneReq struct {
		Id int `json:"id"` // 管理员id
	}

	ResetPasswordReq struct {
		Id int `json:"id"` // 管理员id
	}

	EditAdminReq struct {
		Id              int    `json:"id" valid:"Must;ErrorCode(3031)`     // 管理员id
		Password        string `json:"password"`                           // 密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`          // 确认密码
		RoleId          int    `json:"roleId" valid:"Must;ErrorCode(3029)` // 角色id
	}

	DelAdminReq struct {
		Id int `json:"id" valid:"Must;ErrorCode(3031)` // 系统管理员id
	}

	// 登陆请求
	LoginReq struct {
		Account  string `json:"account" valid:"Must;ErrorCode(3026)`  //账号
		Password string `json:"password" valid:"Must;ErrorCode(3027)` //密码
	}

	// 修改密码请求
	EditPasswordReq struct {
		OldPassword     string `json:"oldPassword" valid:"Must;ErrorCode(3027)"` // 原密码
		Password        string `json:"password" valid:"Must;ErrorCode(3027)"`    // 新密码
		ConfirmPassword string `json:"confirmPassword,omitempty"`                // 确认密码
	}
}

var (
	SystemAdminService = new(services.SystemAdminService)
	SessionService     = new(middleware.SessionService)
)

// 查询全部系统管理员
func (m SystemAdminController) QueryAdminList(ctx server.Context) error {
	req := &m.QueryAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	result, err := SystemAdminService.QuerySystemAdminList(req.StartTime, req.EndTime, req.RoleId, req.IsOnline, req.Account, req.PageIndex, req.PageSize)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, result)
}

// 添加系统管理员
func (m SystemAdminController) AddAdmin(ctx server.Context) error {
	req := &m.AddAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAdminService.AddSystemAdmin(req.Account, req.Password, req.ConfirmPassword, req.RoleId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 根据id查询系统管理员信息
func (m SystemAdminController) QueryAdminOne(ctx server.Context) error {
	req := &m.QueryAdminOneReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	Admin, err := SystemAdminService.QueryAdminOne(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, Admin)
}

// 修改系统管理员
func (m SystemAdminController) EidtAdmin(ctx server.Context) error {
	req := &m.EditAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAdminService.EditSystemAdmin(req.Id, req.Password, req.ConfirmPassword, req.RoleId)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 删除系统管理员
func (m SystemAdminController) DelAdmin(ctx server.Context) error {
	req := &m.DelAdminReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAdminService.DelSystemAdmin(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 初始化密码
func (m SystemAdminController) ResetPassword(ctx server.Context) error {
	req := &m.ResetPasswordReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	err := SystemAdminService.ResetPassword(req.Id)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 管理员登陆
func (m SystemAdminController) Login(ctx server.Context) error {
	req := &m.LoginReq
	if err := ctx.Validate(req); err != nil {
		return err
	}
	// 获取登陆ip
	ip := ctx.RealIP()
	device := ctx.Get(model.DEVICE).(int)
	info, err := SystemAdminService.Login(req.Account, req.Password, ip, device)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, info)
}

// 管理员登出
func (m SystemAdminController) Logout(ctx server.Context) error {
	// 获取登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	err = SystemAdminService.Logout(user)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}

// 修改密码
func (m SystemAdminController) EditPassword(ctx server.Context) error {
	req := &m.EditPasswordReq
	if err := ctx.Validate(req); err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	// 获取用户登陆信息
	user, err := SessionService.GetSession(ctx.Get(model.SessionKey).(string))
	if err != nil {
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.LOGIN_INFO_GET_FAIL})
	}
	err = SystemAdminService.EditPassword(user.Account(), req.OldPassword, req.Password, req.ConfirmPassword)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}
	return common.HttpResultJson(ctx, nil)
}
