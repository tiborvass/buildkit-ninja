package main

import (
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	bkninja "github.com/tiborvass/buildkit-ninja"
)

func main() {
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), bkninja.Build); err != nil {
		panic(err)
	}
}
