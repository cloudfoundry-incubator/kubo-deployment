package proxy

import "net"

func SetNetListen(f func(net, laddr string) (net.Listener, error)) {
	netListen = f
}

func ResetNetListen() {
	netListen = net.Listen
}
