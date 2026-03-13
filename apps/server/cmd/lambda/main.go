package main

import (
	"os"

	appcore "github.com/Hughzu/trackstack/apps/server/internal/core/app"
	"github.com/aws/aws-lambda-go/lambdaurl"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	runtime, err := appcore.NewRuntime()
	if err != nil {
		_, _ = os.Stderr.WriteString("startup error: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		_ = runtime.Close()
	}()

	runtime.Logger.Info("lambda runtime started", "env", runtime.Config.Env)
	lambdaurl.Start(runtime.Handler)
}
