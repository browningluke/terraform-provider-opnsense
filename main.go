package main

import (
	"context"
	"flag"
	"log"

	"github.com/browningluke/terraform-provider-opnsense/internal/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	opnsenseProviderName        = "registry.terraform.io/browningluke/opnsense"
	version              string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serverFactory, _, err := provider.ProtoV6ProviderServerFactory(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		opnsenseProviderName,
		serverFactory,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err.Error())
	}
}
