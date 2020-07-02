package core

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/utility"
	"baseGo/src/red_robot/conf"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/gansidui/gotcp"
)

var (
	n            []string
	WhiteListKey = ""
)

const (
	MaxBodyLen = 5 << 10 //5120KB
)

// Packet
type Packet struct {
	pData []byte `json:"-"`
	Data  string `json:"data"`
	Key   string `json:"key"`
}

func (this *Packet) Serialize() []byte {
	return this.pData
}

func (this *Packet) GetLength() uint32 {
	return binary.BigEndian.Uint32(this.pData[0:4])
}

func (this *Packet) GetBody() []byte {
	return this.pData[4:]
}

func (this *Packet) DecBody() error {
	return json.Unmarshal([]byte(this.pData[4:]), this)
}

func NewPacket(buff []byte, hasLengthField bool) *Packet {
	p := &Packet{}

	if hasLengthField {
		p.pData = buff

	} else {
		p.pData = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.pData[0:4], uint32(len(buff)))
		copy(p.pData[4:], buff)

		//log.Println(p.pData[0:4])
	}

	return p
}

func NewPacketFromStr(buff string, hasLengthField bool) *Packet {
	p := &Packet{}

	bBuff := []byte(buff)

	if hasLengthField {
		p.pData = bBuff

	} else {
		p.pData = make([]byte, 4+len(bBuff))
		binary.BigEndian.PutUint32(p.pData[0:4], uint32(len(bBuff)))
		copy(p.pData[4:], bBuff)
	}

	return p
}

type PacketProtocol struct {
}

func (this *PacketProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {

	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)
	// read length 数据太大的时候，防火墙会屏蔽
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	if length = binary.BigEndian.Uint32(lengthBytes); length > MaxBodyLen {
		return nil, errors.New("the size of packet is larger than the limit")
	}
	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}
	return NewPacket(buff, true), nil
}

type PacketCallback struct {
}

func (this *PacketCallback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	golog.Info("PacketCallback", "OnConnect", "OnConnect:%s", addr)
	//global.Logger.Debug("OnConnect:%s", addr)
	return true
}

func (this *PacketCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	addr := c.GetRawConn().RemoteAddr().String()
	c.PutExtraData(addr)
	addrs := strings.Split(addr, ":")[0]
	num := CheckWhiteList(addrs)
	if num == 0 {
		golog.Error("PacketCallback", "OnMessage", "TCP wallet ip is not in whitelist: %s", nil, addrs)
		//global.Logger.Error("TCP wallet ip is not in whitelist: %s", addrs)
		Send(c, `{"code":1006,"msg":"TCP wallet ip is not in whitelist: `+addrs+`"}`)
		return true
	}
	packet := p.(*Packet)
	golog.Debug("PacketCallback", "OnMessage", "OnMessage:[%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
	err := packet.DecBody()
	if err != nil {
		golog.Error("PacketCallback", "OnMessage", "DecBody err:", err)
		Send(c, `{"code":1001,"msg":"DecBody err"}`)
		return true
	}
	key_real := utility.Md5(packet.Data + conf.GetConfig().Listening.Md5key)
	if key_real != packet.Key {
		golog.Error("PacketCallback", "OnMessage", "MD5Key or DESKey err key_real:%s packet_key:%s", nil, key_real, packet.Key)
		Send(c, `{"code":1002,"msg":"MD5Key or DESKey err"}`)
		return true
	}
	_, errr := DecParams(packet.Data, conf.GetConfig().Listening.Deskey)
	if errr != nil {
		golog.Error("PacketCallback", "OnMessage", "DecParams error Data:%s DESKey:%s  %s", nil, packet.Data, conf.GetConfig().Listening.Deskey, errr.Error())
		Send(c, `{"code":1003,"msg":"MD5Key or DESKey err"}`)
		return true
	}
	body, err := SendDistance(string(packet.GetBody()))
	if err != nil {
		golog.Error("PacketCallback", "OnMessage", "SendDistance err:", err)
		Send(c, `{"code":1005,"msg":"SendDistance`+err.Error()+`"}`)
	} else {
		Send(c, body)
	}
	return true
	//解析并返回数据
}

func SendDistance(msg string) (body string, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", conf.GetConfig().Listening.SendAdd)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		golog.Error("", "SendDistance", "error:", err)
		return body, err
	}
	echoProtocol := PacketProtocol{}
	int, err := conn.Write(NewPacket([]byte(msg), false).Serialize())
	if err != nil {
		golog.Error("", "SendDistance", "error: %v", err, int)
		return body, err
	}

	p, err := echoProtocol.ReadPacket(conn)
	if err == nil {
		q := p.(*Packet)
		body = string(q.GetBody())
		golog.Debug("", "SendDistance", "ReadPacket:%s", body)
		return body, nil
	} else {
		golog.Error("", "SendDistance", "error:", err)
		return body, err
	}
	return body, nil
}

func (this *PacketCallback) OnClose(c *gotcp.Conn) {
	golog.Debug("PacketCallback", "OnClose", "OnClose:s", c.GetExtraData())
}

func Send(c *gotcp.Conn, data string) error {
	return c.AsyncWritePacket(NewPacketFromStr(data, false), 3*time.Second)
}

func DecParams(params, DESKey string) (*conf.RequestBody, error) {
	params = strings.Replace(params, " ", "+", -1) //base64bug
	if len(params) == 0 {
		return nil, nil
	}
	param, err := utility.DesDecrypt([]byte(params), []byte(DESKey))
	if err != nil {
		return nil, err
	}
	golog.Debug("", "DecParams", "params:"+string(param))
	result := new(conf.RequestBody)
	if err := json.Unmarshal(param, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func CheckWhiteList(ip string) (isExist int) {
	WhiteListKey = "whitelist_wallet"
	redisClient := conf.GetRedis().Get()
	datas, err := redisClient.Do("Get", WhiteListKey)
	if err != nil {
		golog.Debug("", "CheckWhiteList", "GetWhiteList error: ", err)
		return 1
	}
	fResults := string(datas.([]byte))
	if len(fResults) > 0 {
		err := json.Unmarshal([]byte(fResults), &n)
		if err != nil {
			golog.Error("", "CheckWhiteList", "json Unmarshal error: ", err)
			return 0
		}
		for _, v := range n {
			if v == ip {
				isExist = 1
				return
			}
		}
		return isExist

	}
	return 1
}
