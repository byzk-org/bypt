package socket

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/byzk-org/bypt/consts"
	"github.com/byzk-org/bypt/tools"
	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
)

var endMsg = []byte("end!!&&")

type Conn struct {
	*gmtls.Conn
	msgChannel chan []byte
}

func (c *Conn) WriteDataStr(data string) *Conn {
	c.WriteData([]byte(data))
	return c
}

func (c *Conn) WriteData(data []byte) *Conn {
	_, err := c.Write(bytes.Join([][]byte{
		[]byte(hex.EncodeToString(data)),
		splitMsg,
	}, nil))
	if err != nil {
		tools.ErrOut("发送数据到服务端失败")
	}
	return c
}

func (c *Conn) ReadMsg() []byte {
	msgByte, isOpen := <-c.msgChannel
	if !isOpen || len(msgByte) == 0 {
		tools.ErrOut("未知的处理异常")
	}

	msgStr, err := hex.DecodeString(string(msgByte))
	if err != nil {
		tools.ErrOut("解析消息结构错误")
	}
	if string(msgStr) == "error" {
		msgByte, isOpen = <-c.msgChannel
		if !isOpen || len(msgByte) == 0 {
			tools.ErrOut("未知的处理异常")
		}
		msgStr, err = hex.DecodeString(string(msgByte))
		if err != nil {
			tools.ErrOut("解析消息结构错误")
		}
		tools.ErrOut(string(msgStr))
	}

	msgByte, isOpen = <-c.msgChannel
	if !isOpen || len(msgByte) == 0 {
		tools.ErrOut("未知的处理异常")
	}

	msgStr, err = hex.DecodeString(string(msgByte))
	if err != nil {
		tools.ErrOut("解析消息结构错误")
	}
	return msgStr
}

func (c *Conn) ReadMsgStr() string {
	return string(c.ReadMsg())
}

func (c *Conn) Wait() *Conn {
	_ = c.ReadMsgStr()
	return c
}

func (c *Conn) SendEndMsg() {
	_, _ = c.Write(endMsg)
	_ = c.Close()
}

func GetClientConn() *Conn {

	if !tools.ValidCertExpire(consts.UserCert) {
		tools.ErrOutAndExit("客户端已过期")
	}

	userCert, err := gmtls.GMX509KeyPairsSingle(consts.UserCert, consts.UserKey)
	if err != nil {
		tools.ErrOutAndExit("解析用户证书失败")
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(consts.CaCert)

	conn, err := gmtls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", consts.ServerPort), &gmtls.Config{
		GMSupport:    &gmtls.GMSupport{},
		ServerName:   consts.CertServerName,
		Certificates: []gmtls.Certificate{userCert},
		RootCAs:      certPool,
		ClientAuth:   gmtls.RequireAndVerifyClientCert,
	})
	if err != nil {
		tools.ErrOutAndExit("访问本地应用服务失败!")
	}

	msgChannel := make(chan []byte)
	tmpMsg = nil
	go func() {
		if err = readData(conn, msgChannel); err != nil {
			close(msgChannel)
		}
	}()

	return &Conn{
		Conn:       conn,
		msgChannel: msgChannel,
	}
}
