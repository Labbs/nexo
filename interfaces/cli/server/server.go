package server

import (
	"context"
	"strconv"

	"github.com/labbs/nexo/application/action"
	"github.com/labbs/nexo/application/apikey"
	"github.com/labbs/nexo/application/auth"
	databaseApp "github.com/labbs/nexo/application/database"
	"github.com/labbs/nexo/application/document"
	"github.com/labbs/nexo/application/drawing"
	"github.com/labbs/nexo/application/group"
	"github.com/labbs/nexo/application/session"
	"github.com/labbs/nexo/application/space"
	"github.com/labbs/nexo/application/user"
	"github.com/labbs/nexo/application/webhook"
	"github.com/labbs/nexo/infrastructure"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/labbs/nexo/infrastructure/cronscheduler"
	"github.com/labbs/nexo/infrastructure/database"
	"github.com/labbs/nexo/infrastructure/http"
	"github.com/labbs/nexo/infrastructure/jobs"
	"github.com/labbs/nexo/infrastructure/logger"
	"github.com/labbs/nexo/infrastructure/persistence"
	routes "github.com/labbs/nexo/interfaces/http"

	"github.com/urfave/cli/v3"
)

// NewInstance creates a new CLI command for starting the server.
// It's called by the main application to add the "server" command to the CLI.
func NewInstance(version string) *cli.Command {
	cfg := &config.Config{}
	cfg.Version = version
	serverFlags := getFlags(cfg)

	return &cli.Command{
		Name:  "server",
		Usage: "Start the Nexo HTTP server",
		Flags: serverFlags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runServer(*cfg)
		},
	}
}

// getFlags returns the list of CLI flags required for the server command.
func getFlags(cfg *config.Config) (list []cli.Flag) {
	list = append(list, config.GenericFlags(cfg)...)
	list = append(list, config.ServerFlags(cfg)...)
	list = append(list, config.LoggerFlags(cfg)...)
	list = append(list, config.DatabaseFlags(cfg)...)
	list = append(list, config.SessionFlags(cfg)...)
	list = append(list, config.RegistrationFlags(cfg)...)
	return
}

// runServer initializes the necessary dependencies and starts the HTTP server.
func runServer(cfg config.Config) error {
	var err error

	// Initialize dependencies
	deps := infrastructure.Deps{
		Config: cfg,
	}

	// Initialize logger
	deps.Logger = logger.NewLogger(cfg.Logger.Level, cfg.Logger.Pretty, cfg.Version)
	logger := deps.Logger.With().Str("component", "interfaces.cli.http.runserver").Logger()

	// Initialize other cron scheduler (go-cron)
	deps.CronScheduler, err = cronscheduler.Configure(deps.Logger)
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.cronscheduler.configure").Msg("Failed to configure cron scheduler")
		return err
	}

	// Initialize database connection (gorm)
	deps.Database, err = database.Configure(deps.Config, deps.Logger)
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.database.configure").Msg("Failed to configure database connection")
		return err
	}

	// Initialize application services
	userPers := persistence.NewUserPers(deps.Database.Db)
	groupPers := persistence.NewGroupPers(deps.Database.Db)
	sessionPers := persistence.NewSessionPers(deps.Database.Db)
	spacePers := persistence.NewSpacePers(deps.Database.Db)
	documentPers := persistence.NewDocumentPers(deps.Database.Db)
	permissionPers := persistence.NewPermissionPers(deps.Database.Db)
	favoritePers := persistence.NewFavoritePers(deps.Database.Db)
	commentPers := persistence.NewCommentPers(deps.Database.Db)
	documentVersionPers := persistence.NewDocumentVersionPers(deps.Database.Db)

	apiKeyPers := persistence.NewApiKeyPers(deps.Database.Db)
	webhookPers := persistence.NewWebhookPers(deps.Database.Db)
	webhookDeliveryPers := persistence.NewWebhookDeliveryPers(deps.Database.Db)
	databasePers := persistence.NewDatabasePers(deps.Database.Db)
	databaseRowPers := persistence.NewDatabaseRowPers(deps.Database.Db)
	drawingPers := persistence.NewDrawingPers(deps.Database.Db)
	actionPers := persistence.NewActionPers(deps.Database.Db)
	actionRunPers := persistence.NewActionRunPers(deps.Database.Db)

	deps.UserApp = user.NewUserApp(deps.Config, deps.Logger, userPers, groupPers, favoritePers)
	deps.SessionApp = session.NewSessionApp(deps.Config, deps.Logger, sessionPers, deps.UserApp)
	deps.SpaceApp = space.NewSpaceApp(deps.Config, deps.Logger, spacePers, documentPers, permissionPers)
	deps.DocumentApp = document.NewDocumentApp(deps.Config, deps.Logger, documentPers, spacePers, permissionPers, commentPers, documentVersionPers)
	deps.AuthApp = auth.NewAuthApp(deps.Config, deps.Logger, deps.UserApp, deps.SessionApp, deps.SpaceApp, deps.DocumentApp)
	deps.ApiKeyApp = apikey.NewApiKeyApp(deps.Config, deps.Logger, apiKeyPers)
	deps.WebhookApp = webhook.NewWebhookApp(deps.Config, deps.Logger, webhookPers, webhookDeliveryPers)
	deps.DatabaseApp = databaseApp.NewDatabaseApp(deps.Config, deps.Logger, databasePers, databaseRowPers, spacePers, permissionPers)
	deps.DrawingApp = drawing.NewDrawingApp(deps.Config, deps.Logger, drawingPers, permissionPers, spacePers)
	deps.ActionApp = action.NewActionApp(deps.Config, deps.Logger, actionPers, actionRunPers)
	deps.GroupApp = group.NewGroupApp(deps.Config, deps.Logger, groupPers, userPers)

	// Initialize HTTP server (fiber + fiberoapi)
	deps.Http, err = http.Configure(deps.Config, deps.Logger, deps.SessionApp, true)
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.http.configure").Msg("Failed to configure HTTP server")
		return err
	}

	// Setup cron jobs
	configJobs := jobs.Config{
		Logger:        deps.Logger,
		CronScheduler: deps.CronScheduler,
		SessionApp:    *deps.SessionApp,
	}

	err = configJobs.SetupJobs()
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.jobs.setup").Msg("Failed to setup cron jobs")
		return err
	}

	// Setup routes
	routes.SetupRoutes(deps)

	// Start HTTP server
	logger.Info().Str("event", "http.runserver.http.listen").Msgf("Starting HTTP server on port %d", cfg.Server.Port)
	err = deps.Http.Fiber.Listen(":" + strconv.Itoa(cfg.Server.Port))
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.http.listen").Msg("Failed to start HTTP server")
		return err
	}

	return nil
}
