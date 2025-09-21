package rest

import (
	"club/internal/api/rest/handler/v1/internalroute/mergereq"

	"git.oceantim.com/backend/packages/golang/essential/healthcheck"
	essmiddleware "git.oceantim.com/backend/packages/golang/essential/middlewares"
)

// SetupAPIRoutes
// @title           			Loyalty Club Service
// @version         			1.0.0
// @description     			This is the API for the Loyalty Club service
// @Host 						api.example.com
// @BasePath  					/
// @Schemes 					https
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func (s *Server) SetupAPIRoutes(
	errRespM essmiddleware.HttpErrorResponse,
	captureSentry essmiddleware.CaptureSentry,
	health healthcheck.HealthHandler,
	mergeRequestHandler *mergereq.Handler,
) {
	r := s.engine
	r.Use(errRespM.Handle(), captureSentry.CaptureSentry(), health.HealthCheck)
	{
		intV1 := r.Group("internal/v1")
		intV1.POST("merge-requests", mergeRequestHandler.MergeRequestWebhook)
	}
}
