package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	namespace string
	conf      string
)

var rootCmd = &cobra.Command{
	Use:   "kubectl cronlogs",
	Short: "Utility for getting logs of cronjobs",
	Long:  "Utility for logs of cronjobs",
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
		fzf := exec.Command("fzf", "--ansi", "--no-preview")
		var out bytes.Buffer
		fzf.Stdin = os.Stdin
		fzf.Stderr = os.Stderr
		fzf.Stdout = &out
		original := strings.Join(os.Args, " ")
		command := original + " cronjoblist"
		choice := runfzf(command)
		if choice == "" {
			panic("BUDDYYYYYYYY")
		}
		command = original + " joblist -j " + choice
		choice = runfzf(command)
		pods, err := clientset.CoreV1().Pods(namespace).List(cmd.Context(), v1.ListOptions{LabelSelector: "job-name=" + choice})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			print(getPodLogs(pod, clientset))
		}

	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	config, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	namespace = config.Contexts[config.CurrentContext].Namespace

	if namespace == "" {
		namespace = "default"
	} else if namespace == "*" {
		namespace = ""
	}
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", namespace, "Namespace to get logs from cronjob within")
	rootCmd.PersistentFlags().StringVar(&conf, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "Kubeconfig to use, defaults to $HOME/.kube/config")
}
