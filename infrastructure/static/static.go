package static

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

var (
	//go:embed files/*
	embedDirStatic embed.FS
)

func NewStatic(f *fiber.App) {
	// Serve static assets
	fsys, _ := fs.Sub(embedDirStatic, "files/assets")
	f.Use("/assets", filesystem.New(filesystem.Config{
		Root: http.FS(fsys),
	}))

	// Serve index.html for SPA routes
	f.Use(func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip API routes
		if strings.HasPrefix(path, "/api") {
			return c.Next()
		}

		// Skip static assets routes
		if strings.HasPrefix(path, "/assets") {
			return c.Next()
		}

		// Serve index.html from the embedded FS for all other routes (SPA routes)
		indexFile, err := embedDirStatic.ReadFile("files/index.html")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Could not find index.html")
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Status(fiber.StatusOK).Send(indexFile)
	})
}
