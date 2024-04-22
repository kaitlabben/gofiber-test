package main

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	h := NewHandler("303", zerolog.New(os.Stdout))

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	app.Use(recover.New())
	app.Use(pprof.New())

	app.Get("/balance", h.Balance)

	err := app.Listen(":9010")
	if err != nil {
		log.Fatal().Msg("failed to listen and serve: err = " + err.Error())
	}
}
