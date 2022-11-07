package util

import (
	"encoding/binary"
	"net"
	"strconv"
	"time"
)

const (
// defaultPageSize             = 20
// maxSOSMsgNum                = 5
// reachMaxLiveLocationError   = "reach_max_live_location"
// noPermissionError           = "no_permission"
// liveLocationInProgressError = "live_location_in_progress"
)

// DirectResponse direct to response
type DirectResponse string

// Direct define
const (
	KafkaDirect      DirectResponse = "kafka"
	KafkaMultiDirect DirectResponse = "kafkamulti"
	TCPDirect        DirectResponse = "tcp"
)

func GetOutboundIP() uint32 {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return 0
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return binary.BigEndian.Uint32(localAddr.IP)
}

func GetIMCTopic(ip, port uint32) string {
	return "imc." + strconv.FormatUint(uint64(ip), 10) + "." + strconv.FormatUint(uint64(port), 10)
}

func Hash33(high uint32, nonce []byte) uint64 {
	hash := uint32(5381)
	for i, b := range nonce {
		if i >= 32 {
			break
		}
		hash = ((hash << 5) + hash) + uint32(b) /* hash * 33 + b */
	}

	return (uint64(high) << 32) | uint64(hash)
}

func GetUTCTime() uint64 {
	return uint64(time.Now().UTC().UnixNano() / 1e3)
}

func GetUTCTimeSecond() uint64 {
	return uint64(time.Now().UTC().UnixNano() / 1e9)
}

func GetTime() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Microsecond))
}

func GetTimeSecond() uint64 {
	return uint64(time.Now().Unix())
}
