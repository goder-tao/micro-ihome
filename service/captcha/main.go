package main

import (
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"ihome/service/captcha/handler"
	pb "ihome/service/captcha/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("captcha"),
	)

	// Register handler
	pb.RegisterCaptchaHandler(srv.Server(), handler.New())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
