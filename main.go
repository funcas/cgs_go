package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/ini.v1"

	"github.com/funcas/cgs/manager"

	"github.com/funcas/cgs/gen-go/process"
	"github.com/funcas/cgs/handler"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/funcas/cgs/container"
)

func main() {
	cfg, _ := ini.Load("./conf/cgs.ini")
	section := cfg.Section("cgs")
	addr := section.Key("server").String() + ":" + section.Key("listen").String()

	secure, _ := section.Key("secure").Bool()
	container.Build()
	manager.LoadTransformer()
	go func() {
		conf := &thrift.TConfiguration{
			ConnectTimeout: time.Duration(5000) * time.Millisecond,
			SocketTimeout:  time.Duration(10000) * time.Millisecond,
		}
		runServer(thrift.NewTTransportFactory(),
			thrift.NewTBinaryProtocolFactoryConf(conf),
			addr, secure)
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

	manager.Destroy()
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
