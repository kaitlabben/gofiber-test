package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type player struct {
	Index    int
	UserName string
	JwtToken string
	Token    string
}
type Response struct {
	Balance   string `json:"balance"`
	UserName  string `json:"user_name"`
	ErrorCode int    `json:"error_code"`
}

var c = fiber.AcquireClient()
var logger = zerolog.New(os.Stdout)

func main() {
	players := []player{
		{
			Index:    0,
			UserName: "Hajen",
			JwtToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6IkZha2VpdF9GYWtlT3BlcmF0b3JfSGFqZW4iLCJzZXNzaW9uX3Rva2VuIjoiMjE1NDNjYWItMmE1Yy00MDgzLTkzYTUtYWFkN2NlZTk0NWQzIiwibmJmIjoxNzEzNzc0MTIxLCJleHAiOjE3MTM3Nzc3MjEsImlhdCI6MTcxMzc3NDEyMSwiaXNzIjoiZ2VwLWJhY2tlbmQiLCJhdWQiOiJnYW1lLWVuZ2luZSJ9.OAz3x3sYlZRb0ryPrvrIyWuxlTbcP1geA_yloEdXiAhkX56Xm-mjJPXBQlu4BsvN1HB7S_j5npZZCKS2p6kUzg",
			Token:    "21543cab-2a5c-4083-93a5-aad7cee945d3",
		},
		{
			Index:    1,
			UserName: "Valen",
			JwtToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1bmlxdWVfbmFtZSI6IkZha2VpdF9GYWtlT3BlcmF0b3JfVmFsZW4iLCJzZXNzaW9uX3Rva2VuIjoiYjUzYzllZDctY2IzYS00MTUyLWE5OWYtZTY2M2MxNDY2MDZiIiwibmJmIjoxNzEzNzc0MjEwLCJleHAiOjE3MTM3Nzc4MTAsImlhdCI6MTcxMzc3NDIxMCwiaXNzIjoiZ2VwLWJhY2tlbmQiLCJhdWQiOiJnYW1lLWVuZ2luZSJ9.4rdTXcxVHnFmCM7qzmshmqP47c26LyCU7abzalyh-i2jNgkcnROWJjSoJktfrYkMgANQu-khGQ3Zepgg2UDQHw",
			Token:    "b53c9ed7-cb3a-4152-a99f-e663c146606b",
		},
	}

	iterations := 500
	t := time.NewTicker(1000 * time.Millisecond)
	defer t.Stop()

	for i := 0; i < iterations; i++ {
		<-t.C

		for _, p := range players {
			go sendReq(p)
		}

	}

}

func sendReq(p player) {
	//a := c.Get("http://localhost:9016/balance")
	a := c.Get("http://localhost:9010/balance")
	defer fiber.ReleaseAgent(a)
	//a.Add("Jwt-Session-Token", p.JwtToken)
	a.Add("Session-Token", p.Token)
	a.Add("Game-Id", "303")

	a.QueryString("username=" + p.UserName)

	statusCode, body, errs := a.Bytes()
	if len(errs) > 0 || statusCode != 200 {
		var err error
		for _, e := range errs {
			err = errors.Join(err, e)
		}
		logger.Error().Err(err).
			Str("method", "get").
			Int("status_code", statusCode).
			Msg("failed sending request to game")
	}
	var res Response
	err := sonic.Unmarshal(body, &res)
	if err != nil {
		logger.Error().Err(err).
			Str("user_name", p.UserName).
			Msg("failed unmarshaling balance response")
	}

	if res.UserName != p.UserName {
		logger.Error().Str("response_user", res.UserName).Str("expected_user", p.UserName).Msg("got wrong user in response")
	}

}
