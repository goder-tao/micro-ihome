package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"ihome/service/reguser/handler"
	pb "ihome/service/reguser/proto"
)

func main() {
	// mysql.Init()
	// Create service
	srv := service.New(
		service.Name("reguser"),
	)

	// Register handler
	pb.RegisterReguserHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
