package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/funcas/cgs/gen-go/process"
	"github.com/funcas/cgs/handler"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/funcas/cgs/container"
)

func main() {

	container.Build()

	go func() {
		runServer(thrift.NewTTransportFactory(), thrift.NewTBinaryProtocolFactoryDefault(), "127.0.0.1:6060", false)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	err := server.Stop()
	if err != nil {
		log.Fatal(err.Error())
	}

	container.Destroy()
	log.Println("Server exited. ")

}

var server *thrift.TSimpleServer

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) {
	var transport thrift.TServerTransport
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		}
		transport, _ = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, _ = thrift.NewTServerSocket(addr)
	}

	log.Printf("%T\n", transport)

	entryServiceHandler := handler.NewEntryServiceHandler()
	processor := process.NewEntryServiceProcessor(entryServiceHandler)
	server = thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	log.Println("Starting the simple server... on ", addr)
	server.Serve()
}
