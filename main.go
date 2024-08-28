package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "location to your kubeconfig file")
	namespace := flag.String("n", "", "namespace to query, empty means all namespaces")
	flag.Parse()

	// Build the Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Get the pods in the specified namespace or across all namespaces
	pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	// Create a map to hold namespace summaries
	summary := make(map[string]*namespaceSummary)

	// Summarize pod information
	for _, pod := range pods.Items {
		ns := pod.Namespace
		if _, exists := summary[ns]; !exists {
			summary[ns] = &namespaceSummary{}
		}

		summary[ns].TotalPods++
		switch pod.Status.Phase {
		case corev1.PodRunning:
			summary[ns].RunningPods++
		case corev1.PodPending:
			summary[ns].PendingPods++
		case corev1.PodFailed:
			summary[ns].FailedPods++
		}

		// Summarize resource usage (simplified)
		for _, container := range pod.Spec.Containers {
			usage := summary[ns]
			if container.Resources.Requests != nil {
				cpu := container.Resources.Requests.Cpu().MilliValue()
				mem := container.Resources.Requests.Memory().Value()
				usage.TotalCPU += cpu
				usage.TotalMemory += mem
			}
		}
	}

	// Print the summary
	fmt.Printf("%-20s %5s %8s %8s %8s %8s %10s\n", "NAMESPACE", "PODS", "RUNNING", "PENDING", "FAILED", "CPU(m)", "MEMORY(Mi)")
	fmt.Println("-------------------------------------------------------------------------------")

	for ns, nsSummary := range summary {
		fmt.Printf("%-20s %5d %8d %8d %8d %8d %10d\n",
			ns,
			nsSummary.TotalPods,
			nsSummary.RunningPods,
			nsSummary.PendingPods,
			nsSummary.FailedPods,
			nsSummary.TotalCPU,
			nsSummary.TotalMemory/(1024*1024),
		)
	}
}

// namespaceSummary holds the summary information for a namespace
type namespaceSummary struct {
	TotalPods   int
	RunningPods int
	PendingPods int
	FailedPods  int
	TotalCPU    int64 // in millicores
	TotalMemory int64 // in bytes
}
