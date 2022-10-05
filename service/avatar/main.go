package main

import (
	"ihome/service/avatar/handler"
	pb "ihome/service/avatar/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("avatar"),
	)

	// Register handler
	pb.RegisterAvatarHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
