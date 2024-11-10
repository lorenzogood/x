package main

import (
	"github.com/lorenzogood/x/internal/startup"
)

func main() {
	ctx, cancel := startup.Run("TEST")
	defer cancel()

	startup.Metrics(ctx)

	<-ctx.Done()
}
