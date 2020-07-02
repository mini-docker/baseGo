package grpc

import (
	"errors"

	"baseGo/src/fecho/golog"
	"baseGo/src/imserver/pkg/bufio"
	"baseGo/src/imserver/pkg/bytes"
	"baseGo/src/imserver/pkg/encoding/binary"
	"baseGo/src/imserver/pkg/websocket"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Proto{Op: OpProtoReady}
	// ProtoFinish proto finish
	ProtoFinish = &Proto{Op: OpProtoFinish}
)

// WriteTo write a proto to bytes writer.
func (p *Proto) WriteTo(b *bytes.Writer) {
	var (
		packLen = _rawHeaderSize + int32(len(p.Body))
		buf     = b.Peek(_rawHeaderSize)
	)
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	if p.Body != nil {
		b.Write(p.Body)
	}
}

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		golog.Error("Proto", "ReadTCP", "error(%v)", nil, err)
		return
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	//fmt.Println("----------------------------------------")
	return
}

// WriteTCP write a proto to TCP writer.
func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	var (
		buf     []byte
		packLen int32
	)
	if p.Op == OpRaw {
		// write without buffer, job concact proto into raw buffer
		_, err = wr.WriteRaw(p.Body)
		if err != nil {
			golog.Error("Proto", "WriteTCP", "error(%v)", nil, err)
		}
		return
	}
	packLen = _rawHeaderSize + int32(len(p.Body))
	if buf, err = wr.Peek(_rawHeaderSize); err != nil {
		golog.Error("Proto", "WriteTCP", "error(%v)", nil, err)
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	if p.Body != nil {
		_, err = wr.Write(p.Body)
		if err != nil {
			golog.Error("Proto", "WriteTCP", "error(%v)", nil, err)
		}
	}
	return
}

// WriteTCPHeart write TCP heartbeat with room online.
func (p *Proto) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + _heartSize
	if buf, err = wr.Peek(packLen); err != nil {
		golog.Error("Proto", "WriteTCPHeart", "error(%v)", nil, err)
		return
	}
	// header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	// body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	return
}

// ReadWebsocket read a proto from websocket connection.
func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		golog.Error("Proto", "ReadWebsocket", "error(%v)", nil, err)
		return
	}
	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}
	return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Proto) WriteWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + len(p.Body)
	if err = ws.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		golog.Error("Proto", "WriteWebsocket", "ws.WriteHeader(%v, %v) error(%v)", nil, websocket.BinaryMessage, packLen, err)
		return
	}
	if buf, err = ws.Peek(_rawHeaderSize); err != nil {
		golog.Error("Proto", "WriteWebsocket", "ws.Peek(%v) error(%v)", nil, _rawHeaderSize, err)
		return
	}
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	if p.Body != nil {
		err = ws.WriteBody(p.Body)
	}
	return
}

// WriteWebsocketHeart write websocket heartbeat with room online.
func (p *Proto) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
	var (
		buf     []byte
		packLen int
	)
	packLen = _rawHeaderSize + _heartSize
	// websocket header
	if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
		golog.Error("Proto", "WriteWebsocketHeart", "wr.WriteHeader(%v, %v) error(%v)", nil, websocket.BinaryMessage, packLen, err)
		return
	}
	if buf, err = wr.Peek(packLen); err != nil {
		golog.Error("Proto", "WriteWebsocketHeart", "wr.Peek(%v) error(%v)", nil, packLen, err)
		return
	}
	// proto header
	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(buf[_verOffset:], int16(p.Ver))
	binary.BigEndian.PutInt32(buf[_opOffset:], p.Op)
	binary.BigEndian.PutInt32(buf[_seqOffset:], p.Seq)
	// proto body
	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
	return
}
