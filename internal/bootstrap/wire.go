// internal/bootstrap/wire.go
package bootstrap

import (
	"helloworld-go/internal/data/db"
	"helloworld-go/internal/data/persistence"
	grpcTransport "helloworld-go/internal/server/grpc"
	httpServerHandler "helloworld-go/internal/server/http/handler"
	httpServerRouter "helloworld-go/internal/server/http/router"
	userService "helloworld-go/internal/service/user"

	"github.com/gin-gonic/gin"
)

func Build() (*gin.Engine, *grpcTransport.UserServer, error) {
	gormDB, err := db.NewDB()
	if err != nil {
		return nil, nil, err
	}

	userRepo := persistence.NewUserRepo(gormDB)
	userSvc := userService.NewUserService(userRepo)

	// HTTP
	ginEngine := gin.Default()
	userHandler := httpServerHandler.NewUserHandler(userSvc)
	httpServerRouter.RegisterUserRoutes(ginEngine, userHandler)

	// gRPC
	grpcServer := grpcTransport.NewUserServer(userSvc)

	return ginEngine, grpcServer, nil
}
