package main

import (
	"context"
	"flag"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func deplymentList() {
	kubeconfig := flag.String("kubeconfig", "/Users/kritiksachdeva/.kube/config", "Location to k8s configuration")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		// Handle error
		fmt.Printf("error %s building config from kubeconfig\n", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// Handle error
		fmt.Printf("erorr %s building clientset from config\n", err.Error())
	}

	deplymentList, err := clientset.AppsV1().Deployments("kube-system").List(context.Background(), v1.ListOptions{})
	if err != nil {
		// Handle error
		fmt.Printf("error %s getting deployment list from kube-system namespace\n", err.Error())
	}

	for _, deploy := range deplymentList.Items { // This is similar to for i := 0; i < 10; i++ ... Hence, the code is the first item i.e declaration of vars
		fmt.Println(deploy.Name)
	}
}
