package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1machinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// connection
type Connection struct {
	Host     string
	Token    string
	Insecure bool
}

// request
type ContextRequest struct {
	Name        string
	Host        string
	ClusterName string
}

// function that creates the context
func CreateContext(request *ContextRequest, connection *Connection, token *Token) error {

	// set host
	request.Host = connection.Host

	// build connection
	var config rest.Config
	config.BearerToken = connection.Token
	config.Host = connection.Host
	config.Insecure = connection.Insecure

	// create the clientset
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		fmt.Errorf("Error connecting to cluster\n")
		return err
	}

	// get the service accounts
	fmt.Print("Reading cluster details... ")
	serviceAccount, err := clientset.CoreV1().ServiceAccounts(token.Namespace).Get(context.TODO(), token.ServiceAccount, v1machinery.GetOptions{})
	if err != nil {
		fmt.Errorf("Error querying service account\n")
		return err
	} else {
		fmt.Println("Success")
	}

	// get the secret
	secret, err := clientset.CoreV1().Secrets(token.Namespace).Get(context.TODO(), token.Secret, v1machinery.GetOptions{})
	if err != nil {
		fmt.Errorf("Error querying secret\n")
		return err
	}

	// get the details from the secreate
	crt := string(secret.Data["ca.crt"])

	// set context name
	contextName := request.Name
	if len(contextName) == 0 {
		contextName = fmt.Sprintf("%s-%s", request.ClusterName, serviceAccount.Name)
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
	c := exec.Command(
		"kubectl",
		"config",
		"set-credentials",
		fmt.Sprintf("%s-user", contextName),
		"--token",
		token.Val)
	_, err = c.Output()
	if err != nil {
		return err
	}

	// create cluster
	c = exec.Command(
		"kubectl",
		"config",
		"set-cluster",
		request.ClusterName,
		fmt.Sprintf("--server=%s", request.Host),
		fmt.Sprintf("--certificate-authority=%s", certPath))
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
