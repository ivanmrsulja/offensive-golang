package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <network>\n", os.Args[0])
		fmt.Println("Example: go run main.go 192.168.1.0/24")
		os.Exit(1)
	}

	network := os.Args[1]
	port := "22" // Default SSH port

	ip, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		fmt.Println("Error parsing CIDR:", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	ipCh := make(chan string, 256)

	for i := 0; i < 256; i++ {
		wg.Add(1)
		go worker(ipCh, port, &wg)
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ipCh <- ip.String()
	}
	close(ipCh)

	wg.Wait()
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func worker(ipCh chan string, port string, wg *sync.WaitGroup) {
	defer wg.Done()
	for ip := range ipCh {
		address := fmt.Sprintf("%s:%s", ip, port)
		if isUp(address) {
			fmt.Printf("Host %s is up\n", ip)
		}
	}
}

func isUp(address string) bool {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return false
	}
	conn.Close()
	fmt.Printf("Successfully connected to %s\n", address)
	return true
}
