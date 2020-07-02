package server

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/red_robot/conf"
	"baseGo/src/red_robot/core"
	"encoding/json"
	"net"
	"strconv"
	"strings"
)

func Wallet(adds, cmd, username, md5key, deskey string, reqMember *conf.ReqMember, remark string, extensionField int) (*conf.ResponseBody, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", adds)
	if err != nil {
		golog.Error("server", "Wallet", "error:", err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		golog.Error("server", "Wallet", "error:", err)
		return nil, err
	}
	defer conn.Close()

	data := new(conf.RequestBody)
	data.Cmd = cmd
	data.SubVt = "pkplus"
	data.Member = reqMember
	data.Data = remark
	if extensionField == 1 {
		data.ExtensionField = strconv.Itoa(extensionField)
	} else {
		data.ExtensionField = ""
	}
	if strings.ToUpper(data.Cmd) == "TRANSFER" {
		order := OderNo(0)
		data.RequestId = order
	}

	b, _ := json.Marshal(data)
	qParams, err := utility.DesEncrypt(b, []byte(deskey))
	key := utility.Md5(qParams + md5key)
	msgData := new(conf.Packet)
	msgData.Platform = "pkplus"
	msgData.Key = key
	msgData.Data = qParams
	msgByte, _ := json.Marshal(msgData)

	echoProtocol := &core.PacketProtocol{}
	// ping <--> pong
	// write
	//超过2M的防火墙会屏蔽
	_, err = conn.Write(core.NewPacket(msgByte, false).Serialize())
	if err != nil {
		golog.Error("server", "Wallet", "error:", err)
		return nil, err
	}

	// read
	p, err := echoProtocol.ReadPacket(conn)
	if err != nil {
		golog.Error("server", "Wallet", "error:", err)
		return nil, err
	}
	packet := p.(*core.Packet)
	respData := new(conf.ResponseBody)
	json.Unmarshal(packet.GetBody(), respData)
	return respData, nil
}
