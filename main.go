package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1a "k8s.io/api/core/v1"
	v1b "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"romacode.com/k8s-context/models"
)

func main() {

	fmt.Println("Hello World!")

	k8sTest()
}

func k8sTest() {

	// get the args
	//args := os.Args[1:]

	var err error

	// get the connection info
	var connection *models.Connection
	if connection, err = getConnection(); err != nil {
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

	// get the namespace
	if len(connection.Namespace) == 0 {
		fmt.Print("Namespace: ")
		fmt.Scanf("%s", &connection.Namespace)
	}

	// get the context defaults
	contextCreate, err := getContextDetaults()
	if err != nil {
		fmt.Println("Error getting context details:")
		fmt.Println(err)
	}

	// get the cluster name
	if len(contextCreate.ClusterName) == 0 {
		fmt.Print("Cluster: ")
		fmt.Scanf("%s", &contextCreate.ClusterName)
	}

	// validate
	if len(connection.Host) == 0 || len(connection.Token) == 0 || len(connection.Namespace) == 0 || len(contextCreate.ClusterName) == 0 {
		fmt.Println("Error, --host, --token, --namespace, --cluster required")
		return
	}

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
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts(connection.Namespace).List(context.TODO(), v1b.ListOptions{})
	if err != nil {
		fmt.Println("Error querying service accounts")
		fmt.Println(err)
		return
	} else {
		fmt.Println("Success")
	}

	// go through each service account
	for _, sa := range serviceAccounts.Items {

		// get the secret
		if secret, err := clientset.CoreV1().Secrets(sa.Namespace).Get(context.TODO(), sa.Secrets[0].Name, v1b.GetOptions{}); err == nil {

			// check if a match on the token
			if string(secret.Data["token"]) == connection.Token {
				fmt.Println("*** FOUND ****")
				fmt.Println(sa.Name)

				// update connect request
				contextCreate.Host = connection.Host

				// run
				if err = execute(contextCreate, &sa, secret); err != nil {
					fmt.Println("Error creating context:")
					fmt.Print(err)
					return
				} else {
					fmt.Println("Context successfully created")
				}
			}

		}
	}
}

// go
func execute(request *models.ContextRequest, sa *v1a.ServiceAccount, secret *v1a.Secret) error {

	// get the details from the secreate
	token := string(secret.Data["token"])
	crt := string(secret.Data["ca.crt"])

	// set context name
	contextName := request.Name
	if len(contextName) == 0 {
		contextName = fmt.Sprintf("%s-%s", request.ClusterName, sa.Name)
	}

	// get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// define cert paths
	certDir := filepath.Join(homeDir, ".kube", request.ClusterName)
	certPath := filepath.Join(certDir, fmt.Sprintf("%s.crt", contextName))

	// show details
	fmt.Printf("Home directory: %s\n", homeDir)
	fmt.Printf("Cluster directory: %s\n", certDir)
	fmt.Printf("Certificate path: %s\n", certPath)

	// make the folder
	err = os.MkdirAll(certDir, os.ModePerm)
	if err != nil {
		return err
	}

	// save cert file
	fil, err := os.Create(certPath)
	if err != nil {
		return err
	}

	// write the file
	if _, err := fmt.Fprint(fil, crt); err != nil {
		fil.Close()
		return err
	} else {
		fil.Close()
	}

	// set credentials
	c := exec.Command("kubectl", "config", "set-credentials", fmt.Sprintf("%s-user", contextName), "--token", token)
	_, err = c.Output()
	if err != nil {
		return err
	}

	// create cluster
	c = exec.Command("kubectl", "config", "set-cluster", request.ClusterName, fmt.Sprintf("--server=%s", request.Host), fmt.Sprintf("--certificate-authority=%s", certPath))
	_, err = c.Output()
	if err != nil {
		return err
	}

	// set context
	c = exec.Command("kubectl", "config", "set-context", contextName, "--user", fmt.Sprintf("%s-user", contextName), "--cluster", request.ClusterName, "--namespace", "default")
	_, err = c.Output()
	if err != nil {
		return err
	}

	// if here, ok
	return nil
}
