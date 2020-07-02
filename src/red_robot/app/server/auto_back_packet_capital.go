package server

import (
	"baseGo/src/fecho/common"
	"baseGo/src/fecho/golog"
	"baseGo/src/model"
	"baseGo/src/model/bo"
	"baseGo/src/red_robot/conf"
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func AutoBackPacketCapital() {
	s := gocron.NewScheduler()
	s.Every(2).Minutes().Do(backPacketCapital)
	<-s.Start()
}

func backPacketCapital() {
	sess := conf.GetXormSession()
	defer sess.Close()

	// 获取所有结算退还保证金异常的会员
	users, err := new(bo.User).GetAllBackCapital(sess)

	if err != nil {
		golog.Error("AutoBackPacketCapital", "backPacketCapital", "获取所有结算退还保证金异常的会员失败", err)
		return
	}

	// 便利数据退还保证金
	for _, user := range users {
		// 获取线路信息获取当前是钱包模式还是额度转换模式
		lineInfo, err := GetStytemLineInfo(user.LineId)
		if err != nil {
			golog.Error("AutoBackPacketCapital", "backPacketCapital", "获取线路信息失败:", err)
			continue
		}
		// 返还异常红包押金
		var remark string
		if user.Capital > 0 {
			remark = fmt.Sprintf("返还结算异常未退还的保证金%v元", user.Capital)
		} else {
			remark = fmt.Sprintf("扣除结算异常多退还的保证金%v元", common.DecimalSub(0, user.Capital))
		}

		if user.IsRobot != model.USER_IS_ROBOT_YES {
			err = new(UserServer).ChangeAmount(sess, user.LineId, user.AgencyId, user.Account, lineInfo.ApiUrl, lineInfo.TransType, user.Id, user.Capital, remark, common.DecimalSub(0, user.Capital))
			if err != nil {
				sess.Rollback()
				golog.Error("AutoBackPacketCapital", "backPacketCapital", "退还异常保证金失败:", nil, err.Error(), user.LineId, user.AgencyId, user.Account)
				continue
			}
		}
	}
}
