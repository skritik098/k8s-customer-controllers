package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Create a channel to record the interrupt signals
func createSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		fmt.Printf("Signal handler: received signal %s\n", sig)
		close(stop)
	}()
	return stop
}

func addPod(new interface{}) {
	fmt.Println("inside add function")
}

func deletePod(obj interface{}) {
	fmt.Println("inside delete function")
}
func updatePod(old, new interface{}) {
	fmt.Println("inside update function")
}

func main() {
	kuberconfig := flag.String("kubeconfig", "/Users/kritiksachdeva/.kube/config", "Locati on to k8s configuration")
	config, err := clientcmd.BuildConfigFromFlags("", *kuberconfig)
	if err != nil {
		// Handle error reading the kubeconfig
		fmt.Printf("error %s reading kubeconfig", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error %s reading config from service account", err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error %s reading clientset", err.Error())
	}

	// Initialize the informer resource. Here we will be using sharedinformer factory instead of simple informers
	// because in case if we need to query / watch multiple Group versions

	// NewSharedInformerFactory will creates a new ShareInformerFactory for "all namespaces"
	informerfactory := informers.NewSharedInformerFactory(clientset, 30*time.Second) // 30*time.Second is the resyc period to update the in-memory cache of informer

	// From this informerfactory we can create specific informers for every group version resource

	podinformer := informerfactory.Core().V1().Pods()

	// Once we have our informer, we need to start it to initialize the in-memory for that informer by calling the List() method
	// But before that we need to set the event handler function
	podinformer.Informer().AddEventHandler(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    addPod,
			DeleteFunc: deletePod,
			UpdateFunc: updatePod,
		})
	// Now we can start our informer, but before starting we need to finalize
	// how we are going to stop this informer/controller!

	//stopChannel := createSignalHandler()

	//informerfactory.Start(stopChannel)

	informerfactory.Start(wait.NeverStop)

	// Since the in-memory store/cache to get fully initialized and that might require some waiting time
	// if podinformer.Informer().HasSynced(){}

	// informerfactory.WaitForCacheSync(stopChannel)

	informerfactory.WaitForCacheSync(wait.NeverStop)

	// Once the in-store has been initialised, informer will be going to use the watch() for subsequent calls to api Server

	fmt.Println("Running informer lister")

	pod, err := podinformer.Lister().Pods("kube-system").Get("etcd-minikube")
	if err != nil {
		fmt.Printf("Error %s reading informer list", err.Error())
	}

	fmt.Printf(pod.Name)
}
