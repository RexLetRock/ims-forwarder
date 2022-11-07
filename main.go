package main

import "ims-forwarder/zserver"

func main() {
	addr := "127.0.0.1:9000"
	zserver.ServerStart(addr)
}
