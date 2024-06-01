package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go-tech/internal/configs"
	"go-tech/internal/dao/db"
	"go-tech/internal/logging"
	"go-tech/internal/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var apiServerCmd = &cobra.Command{
	Use:  "api",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		apiServer()
	},
}

func AddApi(root *cobra.Command) {
	root.AddCommand(apiServerCmd)
}

func apiServer() {
	ctx := context.Background()
	err := configs.InitConfig(ctx)
	err = db.InitDb(ctx)
	if err != nil {
		panic(err)
	}
	httpSrv := startHttpServer(ctx)
	// signal
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if httpSrv != nil {
		stopHttpServer(ctx, httpSrv, 5*time.Second)
		httpSrv = nil
	}
}
func startHttpServer(ctx context.Context) *http.Server {
	//@mark: initialize http web server and start
	router := routers.InitRouter()
	addr := configs.Conf().HTTPServer.Addr
	logging.Info(ctx).Msgf("Application started, listening and serving HTTP on: %s", addr)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Fatal(ctx).Msgf("listen: %v", err.Error())
		}
	}()
	return srv
}

func stopHttpServer(ctx context.Context, server interface{}, ts time.Duration) {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	srv := server.(*http.Server)
	ctx, cancel := context.WithTimeout(ctx, ts)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Fatal(ctx).Msgf("Server forced to shutdown: %v", err.Error())
	}
}
