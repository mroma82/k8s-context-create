package main

import (
	"flag"

	"romacode.com/k8s-context/models"
)

// get the connection info from args/env
func getConnection() (*models.Connection, error) {

	// init model
	connection := models.Connection{}

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
func getContextDefaults() (*models.ContextRequest, error) {

	// init model
	contextCreate := models.ContextRequest{}

	// context name
	flag.StringVar(&contextCreate.Name, "context", "", "")

	// cluster name
	flag.StringVar(&contextCreate.ClusterName, "cluster", "", "")

	// parse
	flag.Parse()

	// return
	return &contextCreate, nil
}
