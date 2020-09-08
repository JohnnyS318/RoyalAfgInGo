package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/JohnnyS318/RoyalAfgInGo/pkg/auth/pkg/auth/config"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/auth/pkg/auth/handlers"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/protos"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/shared/pkg/log"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/shared/pkg/mw"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/shared/pkg/utils"
	"google.golang.org/grpc"

	"github.com/Kamva/mgm/v3"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Start starts the account service
func Start() {
	logger := log.NewLogger()

	logger.Warn("Application started. Router will be configured next")
	defer logger.Warn("Application shut down")
	defer logger.Desugar().Sync()

	config.ConfigureDefaults()

	r := mux.NewRouter()

	cfg := &mgm.Config{CtxTimeout: viper.GetDuration(config.DatabaseTimeout)}
	err := mgm.SetDefaultConfig(cfg, viper.GetString(config.DatabaseName), options.Client().ApplyURI(viper.GetString(config.DatabaseUrl)))
	if err != nil {
		logger.Fatalw("Connection to mongo failed", "error", err)
	}
	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		logger.Fatalw("Could not get the mongo client", "error", err)
	}
	defer utils.DisconnectClient(logger, client)

	// Grpc Setup

	conn, err := grpc.DialContext(context.Background(), viper.GetString("UserService.Url"), grpc.WithBlock(), grpc.WithInsecure())

	if err != nil {
		logger.Fatalw("Connection could not be established", "error", err, "target", viper.GetString("UserService.Url"))
	}

	defer conn.Close()

	userService := protos.NewUserServiceClient(conn)

	//userDb := database.NewUserDatabase(logger)

	// Register Middleware
	loggerHandler := mw.NewLoggerHandler(logger)

	stdChain := alice.New(loggerHandler.LogRouteWithIP)

	// Handlers
	userHandler := handlers.NewUserHandler(logger, userService)
	authmwHandler := mw.NewAuthMWHandler(logger, viper.GetString(config.JwtSigningKey))

	// Get Subrouters
	postRouter := r.Methods(http.MethodPost).Subrouter()
	getRouter := r.Methods(http.MethodGet).Subrouter()

	prChain := stdChain.Append(loggerHandler.ContentTypeJSON)

	postRouter.Handle("/account/register", prChain.ThenFunc(userHandler.Register))
	postRouter.Handle("/account/login", prChain.ThenFunc(userHandler.Login))
	postRouter.Handle("/account/logout", prChain.Append(authmwHandler.AuthMW).ThenFunc(userHandler.Logout))

	getRouter.Handle("/account/verify", stdChain.Append(authmwHandler.AuthMW).ThenFunc(userHandler.VerifyLoggedIn))

	logger.Debug("Setup Routes")

	// SERVER SETUP
	port := viper.GetString(config.Port)

	srv := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	utils.StartGracefully(logger, srv, viper.GetDuration(config.GracefulTimeout))
}

func disconnectClient(logger *zap.SugaredLogger, client *mongo.Client) {
	err := client.Disconnect(mgm.Ctx())
	if err != nil {
		logger.Error("MongoDB disconnect", "error", err)
	}
	logger.Warn("Mongodb disconnected")
}