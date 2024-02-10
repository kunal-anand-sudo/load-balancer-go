package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}
type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func (s *SimpleServer) Address() string                                 { return s.addr }
func (s *SimpleServer) IsAlive() bool                                   { return true }
func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) { s.proxy.ServeHTTP(rw, req) }

func newSimpleServer(addr string) *SimpleServer {
	serveUrl, err := url.Parse(addr)
	handleErr(err)
	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serveUrl),
	}
}
func newLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}
func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to next avalaiable server %s \n", targetServer.Address())
	targetServer.Serve(rw, req)
}
func main() {
	servers := []Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.duckduckgo.com"),
	}
	lb := newLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		rw.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		rw.Header().Set("Expires", "0")                                         // Proxies.
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)

}
func handleErr(err error) {
	if err != nil {
		fmt.Println("error: %v\n", err)
		os.Exit(1)
	}
}
