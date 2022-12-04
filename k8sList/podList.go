package main

import (
	"context"
	"flag"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func podsList() {
	// It will expects the program to provide flag name "kubeconfig"
	// Here the flag.String cannot parse the "~" operator for home directory

	kubeconfig := flag.String("kubeconfig", "/Users/kritiksachdeva/.kube/config", "Location to k8s configuration")

	// Now we have the k8s config file, next we would need to client-go library to connect to k8s cluster
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		// Handle error for reading the kubeconfig
		fmt.Println(config)
	}
	// Next we need to make a clientset using the above configuration reference,
	// which is used to interact the k8s resources from different API version
	// For ex: Pod resources are available in API version "coreV1" whereas Deployement
	// resources comes from "appsV1" API version

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// Handle error for client set
		fmt.Println(err)
	}

	// Next we need to talk to k8s APi server to collect Pods from a namespace

	podsList, err := clientset.CoreV1().Pods("kube-system").List(context.Background(), v1.ListOptions{})
	if err != nil {
		//Handle error
		fmt.Println(err)
	}

	fmt.Println("Printing/Listing Pods from the kube-system namespace")
	for _, pod := range podsList.Items {
		fmt.Printf("%s\n", pod.Name)
	}

}
