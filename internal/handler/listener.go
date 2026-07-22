package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/auth"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func DefineListener(port int) {
	e := echo.New()

	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		e.Logger.Error("JWKS_URL environment variable is not set")
		os.Exit(1)
	}

	auth, err := auth.NewCognitoAuthenticator(context.Background(), jwksURL)
	if err != nil {
		e.Logger.Error("failed to create authenticator", "error", err)
		os.Exit(1)
	}
	server := NewServer(auth)

	mw, err := server.CreateMiddleware()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Use(mw...)

	api.RegisterHandlers(e, server)

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Error("shutting down the server", "error", err)
	}
}
