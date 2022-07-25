package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/Users/sandeep/.kube/config", "location of the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Printf("Error: %v", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Error getting incluster config: %v", err)
			return
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//informerfactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	tweakOptions := func(options *metav1.ListOptions) {
		options.LabelSelector = "app=k8s-dev"

	}

	informerfactory := informers.NewFilteredSharedInformerFactory(clientset, 10*time.Minute, "default", tweakOptions)
	informer := informerfactory.Apps().V1().Deployments()

	//stopper := make(chan struct{})
	//defer close(stopper)

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("New Deployment is created \n")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Printf("Deployment is updated \n")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("Deployment Delete triggered \n")
		},
	})

	//informer.Informer().Run(stopper)

	informerfactory.Start(wait.NeverStop)
	informerfactory.WaitForCacheSync(wait.NeverStop)
	deployment, err := informer.Lister().Deployments("default").Get("k8s-dev")
	if err != nil {
		panic(err)
	}
	fmt.Printf("deployment: %v\n", deployment.Name)

	//To make sure it never stops
	<-wait.NeverStop

}
