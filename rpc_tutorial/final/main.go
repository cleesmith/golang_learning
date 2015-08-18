package main

func main() {
	proxy := &LcaProxy{
		Addr:     ":8080",
		RpcAddr:  ":8079",
		Requests: make(map[string]*RequestStats),
		Balancer: MakeBalancer([]string{"127.0.0.1:8081"}),
	}
	proxy.GoToWork()
}
