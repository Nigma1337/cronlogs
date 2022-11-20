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
	name string
)

var joblistCmd = &cobra.Command{
	Use:   "joblist",
	Short: "Lists different jobs",
	Long:  "Lists different jobs, to be used with fzf",
	PreRun: func(cmd *cobra.Command, args []string) {
		if namespace == "*" {
			namespace = ""
		}
		if name == "" {
			panic("BUDDY YOU GOTTA PICK A JOb")
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
		jobs := getJobList(namespace, clientset, name)
		for _, job := range *jobs {
			names = append(names, job.Name)
		}
		fmt.Println(strings.Join(names, "\n"))
	},
}

func init() {
	rootCmd.AddCommand(joblistCmd)
	joblistCmd.PersistentFlags().StringVarP(&name, "job", "j", "", "Namespace to get logs from cronjob within")
}
