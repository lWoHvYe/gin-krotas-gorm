// internal/bootstrap/wire.go
package bootstrap

import (
	"helloworld-go/internal/data/db"
	"helloworld-go/internal/data/persistence"
	grpcTransport "helloworld-go/internal/server/grpc"
	httpTransport "helloworld-go/internal/server/http"
	"helloworld-go/internal/service/user"

	"github.com/gin-gonic/gin"
)

func Build() (*gin.Engine, *grpcTransport.Server, error) {
	gormDB, err := db.NewDB()
	if err != nil {
		return nil, nil, err
	}

	userRepo := persistence.NewUserRepo(gormDB)
	userSvc := user.NewService(userRepo)

	// HTTP
	ginEngine := gin.Default()
	httpHandler := httpTransport.NewHandler(userSvc)
	httpHandler.RegisterRoutes(ginEngine)

	// gRPC
	grpcServer := grpcTransport.NewServer(userSvc)

	return ginEngine, grpcServer, nil
}
