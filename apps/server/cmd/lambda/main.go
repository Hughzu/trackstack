package main

import (
	"os"

	"github.com/Hughzu/trackstack/apps/server/internal/app/monolithapi"
	functionurltransport "github.com/Hughzu/trackstack/apps/server/internal/platform/aws/functionurl"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	runtime, err := monolithapi.NewRuntime()
	if err != nil {
		_, _ = os.Stderr.WriteString("Startup error : " + err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		_ = runtime.Close()
	}()

	runtime.Logger.Info("Lambda runtime started", "env", runtime.Config.Env)
	functionurltransport.StartBuffered(runtime.Handler)
}
