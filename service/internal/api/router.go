package api

import (
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/api/handler"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/api/middleware"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {

	e.Use(middleware.RequestLoggingMiddleware)
	e.Use(middleware.ErrorHandlingMiddleware)
	e.Use(middleware.CORSConfig())

	// Маршрут для проверки доступности сервера
	e.GET("/api/ping", handler.PingHandler)

	orgRepo := repository.NewOrganizationRepository(db)

	tenderRepo := repository.NewTenderRepository(db)
	tenderHandler := handler.NewTenderHandler(tenderRepo, orgRepo)

	// Группа маршрутов для работы с тендерами
	tenderGroup := e.Group("/api/tenders")

	tenderGroup.Use(middleware.ValidationMiddleware)

	tenderGroup.GET("", tenderHandler.ListTendersHandler)
	tenderGroup.GET("/my", tenderHandler.ListUserTendersHandler)
	tenderGroup.POST("/new", tenderHandler.CreateTenderHandler)

	tenderGroup.GET("/:id/status", tenderHandler.GetStatusHandlerHandler, middleware.AccessControlTenderMiddleware(orgRepo))
	tenderGroup.PUT("/:id/status", tenderHandler.SetStatusTenderHandlerHandler, middleware.AccessControlTenderMiddleware(orgRepo))
	tenderGroup.PATCH("/:id/edit", tenderHandler.EditTenderHandler, middleware.AccessControlTenderMiddleware(orgRepo))
	tenderGroup.PUT("/:id/rollback/:version", tenderHandler.RollbackTenderHandler, middleware.AccessControlTenderMiddleware(orgRepo))

	//Группа маршрутов для работы с предложениями
	bidRepo := repository.NewBidRepository(db)
	bidHandler := handler.NewBidHandler(bidRepo, orgRepo)

	bidGroup := e.Group("/api/bids")

	bidGroup.POST("/new", bidHandler.CreateBidHandler)

	bidGroup.GET("/my", bidHandler.ListUserBidHandler)
	bidGroup.GET("/:id/list", bidHandler.ListOrganizationBids)

	bidGroup.GET("/:id/status", bidHandler.GetStatusBidHandler, middleware.AccessControlBidMiddleware(orgRepo))
	bidGroup.PUT("/:id/status", bidHandler.SetStatusBidHandler, middleware.AccessControlBidMiddleware(orgRepo))

	bidGroup.PATCH("/:id/edit", bidHandler.EditBidHandler, middleware.AccessControlBidMiddleware(orgRepo))
	bidGroup.PUT("/:id/rollback/:version", bidHandler.RollbackBidHandler, middleware.AccessControlBidMiddleware(orgRepo))

	bidGroup.PUT("/:id/submit_decision", bidHandler.SendBidDecisionHandler, middleware.AccessControlBidMiddleware(orgRepo))
	bidGroup.PUT("/:id/feedback", bidHandler.SendBidFeedbackHandler, middleware.AccessControlBidMiddleware(orgRepo))
	bidGroup.GET("/:id/reviews", bidHandler.GetAuthorBidReviewsHandler)
}
