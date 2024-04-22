package main

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"
)

type handler struct {
	gameCode string
	logger   zerolog.Logger
}
type Response struct {
	Balance   string `json:"balance"`
	UserName  string `json:"user_name"`
	ErrorCode int    `json:"error_code"`
}

func NewHandler(code string, logger zerolog.Logger) handler {
	return handler{
		gameCode: code,
		logger:   logger,
	}
}

func (h handler) Balance(ctx *fiber.Ctx) error {

	gameId := utils.CopyString(ctx.Get("Game-Id", ""))
	sessionToken := utils.CopyString(ctx.Get("Session-Token", ""))

	if gameId == "" || gameId != h.gameCode {
		h.logger.Error().Msg("missing header")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if sessionToken == "" {
		h.logger.Error().Msg("missing header")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	a := fiber.Get("https://api.restful-api.dev/objects/7")
	defer fiber.ReleaseAgent(a)
	statusCode, _, errs := a.Bytes()
	if len(errs) > 0 || statusCode != 200 {
		var err error
		for _, e := range errs {
			err = errors.Join(err, e)
		}
		h.logger.Error().Err(err).
			Str("method", "get").
			Int("status_code", statusCode).
			Msg("failed sending request to game")
	}

	username := ctx.Query("username")
	res := Response{
		Balance:   "100",
		UserName:  username,
		ErrorCode: 0,
	}
	b, err := sonic.Marshal(res)
	if err != nil {
		h.logger.Error().Msg("failed to marshal response")
	}

	h.logger.Info().Msg(fmt.Sprintf("Sending response %+v", res))
	return ctx.Send(b)
}
