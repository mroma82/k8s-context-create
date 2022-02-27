package main

import (
	"flag"

	"romacode.com/k8s-context/pkg"
)

// get the connection info from args/env
func getConnectionArgs() (*pkg.Connection, error) {

	// init model
	connection := pkg.Connection{}

	// get the host
	flag.StringVar(&connection.Host, "host", "", "")

	// get the token
	flag.StringVar(&connection.Token, "token", "", "")

	// insecure
	flag.BoolVar(&connection.Insecure, "insecure", false, "")

	// parse
	flag.Parse()

	// return
	return &connection, nil
}

// get the context defaults
func getContextDefaultsArgs() (*pkg.ContextRequest, error) {

	// init model
	contextCreate := pkg.ContextRequest{}

	// context name
	flag.StringVar(&contextCreate.Name, "context", "", "")

	// cluster name
	flag.StringVar(&contextCreate.ClusterName, "cluster", "", "")

	// parse
	flag.Parse()

	// return
	return &contextCreate, nil
}
