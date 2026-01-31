package genoapi

import (
	"context"
	"os"

	"github.com/labbs/nexo/infrastructure"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/labbs/nexo/infrastructure/http"
	"github.com/labbs/nexo/infrastructure/logger"
	routes "github.com/labbs/nexo/interfaces/http"
	"github.com/urfave/cli/v3"
)

func NewInstance(version string) *cli.Command {
	cfg := &config.Config{}
	cfg.Version = version
	serverFlags := getFlags(cfg)

	return &cli.Command{
		Name:  "genoapi",
		Usage: "Generate OpenAPI specification file",
		Flags: serverFlags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runCommand(*cfg)
		},
	}
}

// getFlags returns the list of CLI flags required for the server command.
func getFlags(cfg *config.Config) (list []cli.Flag) {
	list = append(list, config.GenericFlags(cfg)...)
	list = append(list, config.ServerFlags(cfg)...)
	list = append(list, config.LoggerFlags(cfg)...)
	list = append(list, config.ExportOapiFlags(cfg)...)
	return
}

func runCommand(cfg config.Config) error {
	var err error

	// Initialize dependencies
	deps := infrastructure.Deps{
		Config: cfg,
	}

	// Initialize logger
	deps.Logger = logger.NewLogger(cfg.Logger.Level, cfg.Logger.Pretty, cfg.Version)
	logger := deps.Logger.With().Str("component", "interfaces.cli.genoapi.runcommand").Logger()

	deps.Http, err = http.Configure(deps.Config, deps.Logger, deps.SessionApp, true)
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.runserver.http.configure").Msg("Failed to configure HTTP server")
		return err
	}

	routes.SetupRoutes(deps)

	spec, err := deps.Http.FiberOapi.GenerateOpenAPISpecYAML()
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.genoapi.generate_openapi_spec_yaml").Msg("Failed to generate OpenAPI spec in YAML format.")
		return err
	}

	err = os.WriteFile(cfg.ExportOapi.FileName, []byte(spec), 0644)
	if err != nil {
		logger.Fatal().Err(err).Str("event", "http.genoapi.write_openapi_yaml_file").Msg("Failed to write OpenAPI YAML file.")
		return err
	}
	logger.Info().Str("event", "http.genoapi.openapi_yaml_exported").Str("file", cfg.ExportOapi.FileName).Msg("OpenAPI spec exported in YAML format.")

	return nil
}
