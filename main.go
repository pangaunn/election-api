package main

import (
	"election-api/graph"
	"election-api/graph/generated"
	"election-api/graph/model"
	"election-api/pkg/cache"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/caarlos0/env"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	logger "github.com/sirupsen/logrus"
)

type Config struct {
	RedisHost     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func main() {
	if currentEnvironment, ok := os.LookupEnv("ENV"); ok {
		if currentEnvironment == "dev" {
			err := godotenv.Load("./.env")
			if err != nil {
				logger.Info("Can't load .env", err)
			}
		}
	}

	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	gqlConfig := generated.Config{
		Resolvers: &graph.Resolver{
			Cache:     cache.InitRedis(config.RedisHost, config.RedisPassword),
			Observers: map[string]chan *model.CandidateVoteUpdated{},
		},
	}

	gqlConfig.Directives.ValidIDCard = graph.ValidIDCard
	srv := handler.New(generated.NewExecutableSchema(gqlConfig))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	http.Handle("/pg", playground.Handler("Election API", "/query"))
	http.Handle("/query", graph.InjectIDCardToCtx(c.Handler(srv)))

	log.Fatal(http.ListenAndServe(":3001", nil))
}
