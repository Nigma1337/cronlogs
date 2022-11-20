package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	names []string
)

var cronjoblistCmd = &cobra.Command{
	Use:   "cronjoblist",
	Short: "Lists different cronjobs",
	Long:  "Lists different cronjobs, to be used with fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		if namespace == "*" {
			namespace = ""
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		cronjobs := getCronjobList(namespace, clientset)
		for _, cronjob := range cronjobs.Items {
			names = append(names, cronjob.Name)
		}
		fmt.Println(strings.Join(names, "\n"))
	},
}

func init() {
	rootCmd.AddCommand(cronjoblistCmd)
}
