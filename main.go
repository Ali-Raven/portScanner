package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	host := flag.String("h", "", "Usage :\n\t./portscan -h <HOST>")
	flag.Parse()
	ch := make(chan string)

	if *host == "" {
		ips, err := getIPs()
		if err != nil {
			panic(err)
		}

		for _, ip := range ips {
			wg.Add(1)
			go portScan(ip, &wg, ch)
		}

	}
	wg.Add(1)
	go portScan(*host, &wg, ch)

	go func() {
		wg.Wait()
		close(ch)
	}()

	var msg string
	if len(os.Args) <= 0 {
		for msg = range ch {
			fmt.Printf("Alert : %s\n", msg)
		}
	} else {
		msg = <-ch
		fmt.Printf("data‌: %s\n", msg)
	}
	fmt.Println("successfully accepted!")
}

func getIPs() ([]string, error) {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}, fmt.Errorf("can't get ips => %v", err)
	}

	var ipv4List []string
	for _, ip := range ips {
		addr, ok := ip.(*net.IPNet)
		if !ok {
			continue
		}

		if addr.IP.IsLoopback() {
			continue
		}

		if addr.IP.To4() != nil {
			ipv4List = append(ipv4List, addr.IP.String())
		}
	}
	return ipv4List, nil
}

func portScan(ip string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	var fullAddr string
	ports := 10000
	for p := range ports {
		if p == 0 {
			continue
		}
		fullAddr = net.JoinHostPort(ip, fmt.Sprintf("%d", p))

		conn, err := net.DialTimeout("tcp", fullAddr, 4*time.Second)
		if err != nil {
			continue
		}

		conn.Close()

		ch <- fmt.Sprintf("IP %s is OPEN", fullAddr)
	}
}
