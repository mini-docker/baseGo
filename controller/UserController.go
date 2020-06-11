package controller

import (
	"baseGo/common"
	"baseGo/dto"
	"baseGo/model"
	"baseGo/response"
	"baseGo/util"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// body传入
// {
//   "name": "test",
//   "account": "12345678900",
//   "password": "123456"
// }

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser)
	// 获取参数
	name := requestUser.Name
	account := requestUser.Account
	password := requestUser.Password
	fmt.Println(password, "password", account)
	// 数据验证
	// if len(account) != 11 {
	// 	response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "账号必须为11位")
	// 	return
	// }
	if len(account) < 1 || len(account) > 13 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "账号必须在2-12位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 如果名称为空给一个随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	if isAccountExist(DB, account) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Account:  account,
		Password: string(hasePassword),
	}
	// todo
	DB.Create(&newUser)

	// 发送token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Print("token genrate error:%v", err)
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")

}

func Login(c *gin.Context) {
	db := common.GetDB()
	// post body 提交
	var requestUser = model.User{}
	c.Bind(&requestUser)
	// name := requestUser.Name
	account := requestUser.Account
	password := requestUser.Password

	// 获取参数 // form 表单查询
	// account := c.PostForm("account")
	// password := c.PostForm("password")

	// fmt.Println(len(account), "lenphone", account, c)
	// 数据验证
	if len(account) < 1 || len(account) > 13 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "账号必须在2-12位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 判断账号是否存在
	var user model.User
	db.Where("account=?", account).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	// 发送token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
	}

	// 返回结果
	// response.Success(c, gin.H{"code": 200, "data": gin.H{"token": token}, "success": true}, "登陆成功")
	// response.Success(c, gin.H{"token": token}, "登陆成功")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"token": token}, "success": true, "msg": "登陆成功"})

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func isAccountExist(db *gorm.DB, account string) bool {
	var user model.User
	db.Where("account=?", account).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
