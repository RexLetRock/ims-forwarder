package main

import "ims-forwarder/zserver"

func main() {
	addr := "0.0.0.0:9000"
	zserver.ServerStart(addr)
}
