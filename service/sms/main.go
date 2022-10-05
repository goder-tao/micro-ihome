package main

import (
	"ihome/service/sms/handler"
	pb "ihome/service/sms/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("sms"),
	)
	//client.NewRegistry()
	//server.DefaultServer = grpc.NewServer(server.Registry())

	// Register handler
	pb.RegisterSmsHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
