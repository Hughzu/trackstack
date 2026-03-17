package main

import (
	"os"

	bootstrap "github.com/Hughzu/trackstack/apps/server/internal/app/bootstrap"
	functionurltransport "github.com/Hughzu/trackstack/apps/server/internal/transport/functionurl"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	runtime, err := bootstrap.NewRuntime()
	if err != nil {
		_, _ = os.Stderr.WriteString("startup error: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		_ = runtime.Close()
	}()

	runtime.Logger.Info("lambda runtime started", "env", runtime.Config.Env)
	functionurltransport.StartBuffered(runtime.Handler)
}
