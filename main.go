package main

import (
	"context"
	"database/sql"
	"github.com/aalug/blog-go/api"
	db "github.com/aalug/blog-go/db/sqlc"
	_ "github.com/aalug/blog-go/docs/statik"
	"github.com/aalug/blog-go/gapi"
	"github.com/aalug/blog-go/mail"
	"github.com/aalug/blog-go/pb"
	"github.com/aalug/blog-go/utils"
	"github.com/aalug/blog-go/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rakyll/statik/fs"
	zerolog "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	// === env ===
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env file: ", err)
	}

	// === postgres ===
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	store := db.NewStore(conn)

	// === redis ===
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	//serverType := os.Getenv("SERVER_TYPE")
	//if serverType == "gin" {
	//	runGinServer(config, store)
	//} else {
	//	runGrpcServer(config, store)
	//}
	go runTaskProcessor(redisOpt, store)
	sender := mail.NewHogSender("gulczas977@o2.pl")
	_ = sender.SendEmail(mail.Data{
		To:      []string{"gulczas977@o2.pl"},
		Subject: "Some Test Subject",
		Content: "<h1>Hello!</h1><br><hr><h5>This is test</h5>",
		Files: []mail.AttachFile{
			{
				Name: "test",
				Path: "./README.md",
			},
		},
	})
	go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}

func runGrpcServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBlogGoServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create a listener:", err)
	}

	log.Printf("gRPC server listening at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start the gRPC server:", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBlogGoHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server: ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// swagger docs server
	statikFileSystem, err := fs.New()
	if err != nil {
		log.Fatal("cannot create a statik file system: ", err)
	}

	docsHandler := http.StripPrefix("/docs/", http.FileServer(statikFileSystem))
	mux.Handle("/docs/", docsHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create a listener:", err)
	}

	log.Printf("HTTP gateway server starting at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start the HTTP gateway server:", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	zerolog.Info().Msg("task processor started")
	err := taskProcessor.Start()
	if err != nil {
		zerolog.Fatal().Err(err).Msg("failed to start task processor")
	}
}
