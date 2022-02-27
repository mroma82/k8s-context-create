package main

import (
	"fmt"

	"romacode.com/k8s-context/pkg"
)

// main
func main() {

	// info
	fmt.Println("Kubernetes Context Create Utility")

	// get the connection info
	connection, err := getConnectionArgs()
	if err != nil {
		fmt.Println("Error getting connection details:")
		fmt.Println(err)
		return
	}

	// get the host
	if len(connection.Host) == 0 {
		fmt.Print("Host: ")
		fmt.Scanf("%s", &connection.Host)
	}

	// get the token
	if len(connection.Token) == 0 {
		fmt.Print("Token: ")
		fmt.Scanf("%s", &connection.Token)
	}

	// show the token
	token, err := pkg.ParseToken(connection.Token)
	if err != nil {
		fmt.Println("Error parsing token:")
		fmt.Println(err)
		return
	}

	// get the context defaults
	contextCreate, err := getContextDefaultsArgs()
	if err != nil {
		fmt.Println("Error getting context details:")
		fmt.Println(err)
		return
	}

	// get the cluster name
	if len(contextCreate.ClusterName) == 0 {
		fmt.Print("Cluster: ")
		fmt.Scanf("%s", &contextCreate.ClusterName)
	}

	// validate
	if len(connection.Host) == 0 || len(connection.Token) == 0 || len(contextCreate.ClusterName) == 0 {
		fmt.Println("Error, --host, --token, --cluster required")
		return
	}

	// create the contenst
	if err = pkg.CreateContext(contextCreate, connection, token); err != nil {
		fmt.Errorf("Error creating context:%s\n", "")
		fmt.Println(err)
	} else {
		fmt.Println("Context successfully created")
	}
}
