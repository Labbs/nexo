package dto

import "github.com/gofiber/fiber/v2"

type SSORedirectOutput struct {
	URL   string
	State string
}

type SSOCallbackInput struct {
	Code    string
	State   string
	Context *fiber.Ctx
}

type SSOCallbackOutput struct {
	Token string
}
