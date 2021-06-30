package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github/rossi1/go-api-microservice-example/cmd/internal"
	"github/rossi1/go-api-microservice-example/handler"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"gorm.io/gorm"
)

// @title API-example
// @version 1.0
// @description This is a go-api-microservice-example.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic BasicAuth

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
func main() {
	ctx := context.Background()

	cfg, err := internal.GetConfig(ctx)
	fmt.Println(cfg.ENVIRONMENT)

	if err != nil {
		log.Fatalf("config error: %s", err)
		os.Exit(1)
	}

	es, err := internal.NewElasticSearch(cfg)

	if err != nil {
		log.Fatalf("elasticsearch error: %s", err)
		os.Exit(1)
	}
	info, code, err := es.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	//err = internal.NewSentry(cfg)

	//if err != nil {
	//	log.Fatalf("sentry error: %s", err)
	//	os.Exit(1)
	//}

	db, err := internal.NewPostgreSQL(cfg)

	if err != nil {
		log.Fatalf("db error: %s", err)
		os.Exit(1)
	}

	//kaftaproducer, err := internal.NewKafkaProducer(cfg)

	//if err != nil {
	//	log.Fatalf("kafta producer error: %s", err)
	//	os.Exit(1)
	//}

	errC, err := run(ctx, cfg, db, es, nil, "")

	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
		os.Exit(1)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
		os.Exit(1)
	}

	//if err := cmd.Execute(ctx, es); err != nil {
	//	log.Fatalf("Error running command: %s", err)
	//	os.Exit(1)

	//}

}

func run(ctx context.Context, cfg internal.Config, db *gorm.DB, es interface{}, producer interface{}, topic string) (<-chan error, error) {

	srv := newServer(internal.GetDeps(cfg, db, es, nil, topic))

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		log.Println("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {

			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Println("Server shutdown completed")
	}()

	go func() {

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil

}

func newServer(conf internal.Deps) *http.Server {
	address := fmt.Sprintf("%s:%s", conf.ServerHost, conf.ServerPort)
	db := conf.DB
	//prod := conf.Producer
	//prodTopic := conf.ProducerTopic
	search := conf.Search
	handler := handler.NewHandler(db, search, nil, "")

	return &http.Server{
		Handler:           handler.GetRouter(),
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}
