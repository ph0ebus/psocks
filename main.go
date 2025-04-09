package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"psocks/server"
)

var Socks5Proxy string

func main() {

	var (
		host = flag.String("host", "0.0.0.0", "监听地址, 默认为0.0.0.0")
		port = flag.Int("port", 1080, "监听端口, 默认为1080")
	)

	flag.StringVar(&Socks5Proxy, "socks", "", "socks5代理地址, 默认为空")
	flag.Parse()

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		slog.Error(fmt.Sprintf("listen failed: %v\n", err))
		return
	}

	slog.Info(fmt.Sprintf("listening at %s:%d", *host, *port))

	for {
		client, err := server.Accept()
		if err != nil {
			slog.Error(fmt.Sprintf("accept failed: %v", err))
			continue
		}
		go process(client)
	}
}

func process(client net.Conn) {
	slog.Info("accept from: " + client.RemoteAddr().String())

	if Socks5Proxy != "" {
		proxy, err := net.Dial("tcp", Socks5Proxy)
		if err != nil {
			fmt.Println("connect to proxy failed:", err)
			client.Close()
			return
		}
		server.Socks5Forward(client, proxy)
		return
	}

	if err := server.Socks5Auth(client); err != nil {
		fmt.Println("auth error:", err)
		client.Close()
		return
	}

	target, err := server.Socks5Connect(client)
	if err != nil {
		fmt.Println("connect error:", err)
		client.Close()
		return
	}

	server.Socks5Forward(client, target)
}
