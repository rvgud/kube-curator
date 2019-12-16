package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Checking in-cluster config\n")
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created kubernetes client\n")
	fmt.Println("Parsing YAML file")
	var fileName string
	flag.StringVar(&fileName, "f", "", "YAML file to parse.")
	flag.Parse()
	fmt.Println(fileName)
	if fileName == "" {
		fmt.Println("Please provide yaml file by using -f option")
		return
	}
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}
	type YamlConfig struct {
		Namespaces map[string]string
	}
	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	fmt.Printf("Result: %#v\n", yamlConfig)
	for k := range yamlConfig.Namespaces {
		namespace := k
		labelSelector := fmt.Sprintf(yamlConfig.Namespaces[k])
		fmt.Printf("Namespace: %#v\n", namespace)
		fmt.Printf("labelSelector: %#v\n", labelSelector)
		err = clientset.CoreV1().Pods(namespace).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Deleted all pods from %s namespace and with specfied labels\n", namespace)
	}

	/**/

	os.Exit(0)
}
