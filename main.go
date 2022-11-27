package main

import "ims-forwarder/zserver"

func main() {
	addr := "0.0.0.0:19000"
	zserver.ServerStart(addr)
}
