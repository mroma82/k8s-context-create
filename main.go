package main

import (
	"context"
	"fmt"

	v1machinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// main
func main() {

	// info
	fmt.Println("Kubernetes Context Create Utility")

	// get the connection info
	connection, err := getConnection()
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
	token, err := parseToken(connection.Token)
	if err != nil {
		fmt.Println("Error parsing token:")
		fmt.Println(err)
		return
	}

	// get the context defaults
	contextCreate, err := getContextDefaults()
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

	// set host
	contextCreate.Host = connection.Host

	// build connection
	var config rest.Config
	config.BearerToken = connection.Token
	config.Host = connection.Host
	config.Insecure = connection.Insecure

	// create the clientset
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		fmt.Println("Error connecting to cluster:")
		fmt.Println(err)
		return
	}

	// get the service accounts
	fmt.Print("Reading cluster details... ")
	serviceAccount, err := clientset.CoreV1().ServiceAccounts(token.Namespace).Get(context.TODO(), token.ServiceAccount, v1machinery.GetOptions{})
	if err != nil {
		fmt.Println("Error querying service account")
		fmt.Println(err)
		return
	} else {
		fmt.Println("Success")
	}

	// get the secret
	secret, err := clientset.CoreV1().Secrets(token.Namespace).Get(context.TODO(), token.Secret, v1machinery.GetOptions{})
	if err != nil {
		fmt.Println("Error querying secret")
		fmt.Println(err)
		return
	}

	// run
	if err = createContext(contextCreate, serviceAccount, secret); err != nil {
		fmt.Println("Error creating context:")
		fmt.Print(err)
		return
	} else {
		fmt.Println("Context successfully created")
	}
}
