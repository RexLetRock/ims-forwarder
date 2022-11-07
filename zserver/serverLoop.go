package zserver

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const cDefaultIMCTopicPort = 2082

func GetIMCTopic(conn net.Conn) string {
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	ip := binary.BigEndian.Uint32(remoteAddr.IP)
	port := cDefaultIMCTopicPort
	result := "imc." + strconv.FormatUint(uint64(ip), 10) + "." + strconv.FormatUint(uint64(port), 10)
	logrus.Errorf("IMC Topic %v %v \n", binary.BigEndian.Uint32(remoteAddr.IP), result)
	return result
}

func (s *ConnHandle) LoopToFlush() {
	for {
		time.Sleep(cTimeToFlush)
		s.flush <- []byte{}
	}
}

func (s *ConnHandle) LoopToWrite() {
	for {
		select {
		case msg := <-s.chans:
			s.slice = append(s.slice, msg...)
			s.cntSent++
			if s.cntSent >= cSendSize {
				if len(s.slice) > 0 {
					go s.conn.Write(s.slice)
					s.slice = []byte{}
					s.cntSent = 0
				}
			}
		case <-s.flush:
			if len(s.slice) > 0 {
				go s.conn.Write(s.slice)
				s.slice = []byte{}
				s.cntSent = 0
			}
		}
	}
}

func (s *ConnHandle) LoopToRead() {
	reader := bufio.NewReader(s.readerReq)
	for {
		msg, err := ReadWithEnd(reader)
		if err != nil {
			logrus.Warnf("Msg %v", err)
			continue
		}

		aMsg := strings.Split(string(msg), "|")
		if len(aMsg) < cMsgPartsNum {
			logrus.Warnf("Msg wrong format %v", aMsg)
			continue
		}

		// Happy handle msg
		switch aMsg[1] {
		case MessageBroadcast.Toa():
			connector, _ := gIPData.Get(aMsg[2])
			if cConn, ok := connector.(*ConnHandle); ok {
				cConn.chans <- msg
			} else {
				logrus.Errorf("IMC Topic not found %v \n", aMsg[2])
			}
		case MessageEdit.Toa():
			logrus.Warn("New topic %+v %v", aMsg, aMsg[2])
			gIPData.Set(aMsg[2], s)
		default:
			// s.chans <- msg
		}
	}
}

func ReadWithEnd(reader *bufio.Reader) ([]byte, error) {
	message, err := reader.ReadBytes('#')
	if err != nil {
		return nil, err
	}

	a1, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a1)
	if a1 != '\t' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}

	a2, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	message = append(message, a2)
	if a2 != '#' {
		message2, err := ReadWithEnd(reader)
		if err != nil {
			return nil, err
		}
		ret := append(message, message2...)
		return ret, nil
	}
	return message, nil
}
