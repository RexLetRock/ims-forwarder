package forwarder

import (
	"io"
	"log"
	"net"
	"time"

	"ims-forwarder/util"
)

const ThreadPerConn = 5

const cSendSize = 999
const cChanSize = 10000
const cBuffSize = 1024 * 1024
const cTimeToFlush = 5 * time.Millisecond
const cMsgPartsNum = 4

var gIPData util.ConcurrentMap

func ServerStart(host string) {
	gIPData = util.NewConcurrentMap()

	log.Printf("Start server at %v\n", host)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Start server err %v \n", err)
	}

	defer listener.Close()
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Printf("Conn accept failed %v \n", err)
		} else {
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	handle := ConnHandleCreate(conn)
	handle.Handle()
}

type ConnHandle struct {
	readerReq *io.PipeReader
	writerReq *io.PipeWriter

	chans chan []byte
	flush chan []byte

	buffer []byte
	slice  []byte
	conn   net.Conn

	cntSent int
}

func ConnHandleCreate(conn net.Conn) *ConnHandle {
	s := &ConnHandle{
		chans:   make(chan []byte, cChanSize),
		flush:   make(chan []byte),
		buffer:  make([]byte, cBuffSize),
		conn:    conn,
		cntSent: 0,
	}
	s.readerReq, s.writerReq = io.Pipe()

	go s.LoopToFlush()
	go s.LoopToWrite()
	go s.LoopToRead()

	return s
}

func (s *ConnHandle) Handle() error {
	defer s.conn.Close()
	for {
		n, err := s.conn.Read(s.buffer)
		if err != nil {
			return err
		}
		s.writerReq.Write(s.buffer[:n])
	}
}
