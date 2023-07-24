package main

import (
	"main/ping"
)

func main() {
	tgt := "192.168.31.%s"
	ping.Scan(tgt, 1, 254)
}
