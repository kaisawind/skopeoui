package service

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"strconv"
)

type Status string

type IService interface {
	// Start 启动服务
	Start() error
	// Stop 停止服务
	Close() error
	// Serve 启动并保持服务运行
	Serve() error
	// Address 返回监听地址
	Address() string
}

func IsListenedAddress(addr string) bool {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}

func GetListeningAddress(addr string) (naddr string, err error) {
	if !IsListenedAddress(addr) {
		return addr, nil
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return
	}
	for range 100 {
		iport := 30000 + rand.IntN(30000) // [30000 , 60000)
		naddr = net.JoinHostPort(host, strconv.Itoa(iport))
		if !IsListenedAddress(naddr) {
			return
		}
	}
	return
}

func GetOutboundIP() (net.IP, error) {
	// 使用一个不一定会连接的外部地址
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func GetLocalIPs() ([]string, error) {
	var ips []string

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 跳过回环接口和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口地址
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 跳过IPv6地址和本地回环地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			ips = append(ips, ip.String())
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no available network interfaces")
	}

	return ips, nil
}

func GetHostIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found for hostname")
}
