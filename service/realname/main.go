package main

import (
	"ihome/service/realname/handler"
	pb "ihome/service/realname/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("realname"),
	)

	// Register handler
	pb.RegisterRealnameHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
