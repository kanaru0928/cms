package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/auth"
	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func DefineListener(port int) {
	e := echo.New()

	appConfig, err := LoadConfig()
	if err != nil {
		e.Logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	auth, err := auth.NewCognitoAuthenticator(context.Background(), appConfig.JWKSURL)
	if err != nil {
		e.Logger.Error("failed to create authenticator", "error", err)
		os.Exit(1)
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		e.Logger.Error("failed to load AWS config", "error", err)
		os.Exit(1)
	}

	contentRepository := content.NewS3Repository(&awsConfig, appConfig.S3BucketName, appConfig.S3Prefix)
	metadataRepository := metadata.NewDynamoDBRepository(&awsConfig, appConfig.DynamoDBTableName)

	server := NewServer(&ServerDependencies{
		Authenticator:      auth,
		ContentRepository:  contentRepository,
		MetadataRepository: metadataRepository,
	})

	mw, err := server.CreateMiddleware()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Use(mw...)

	api.RegisterHandlers(e, server)

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Error("shutting down the server", "error", err)
	}
}
