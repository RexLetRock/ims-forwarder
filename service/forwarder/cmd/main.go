package main

import "ims-forwarder/service/forwarder"

func main() {
	addr := "0.0.0.0:19000"
	forwarder.ServerStart(addr)
}
