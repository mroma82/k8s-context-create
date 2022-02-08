package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"romacode.com/k8s-context/models"
)

// function that creates the context
func createContext(request *models.ContextRequest, sa *v1.ServiceAccount, secret *v1.Secret) error {

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
