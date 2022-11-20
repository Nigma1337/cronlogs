package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	BatchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getCronjobList(namespace string, clientset *kubernetes.Clientset) *BatchV1.CronJobList {
	cronjobs, err := clientset.BatchV1().CronJobs(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	return cronjobs
}

func getJobList(namespace string, clientset *kubernetes.Clientset, name string) *[]BatchV1.Job {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), v1.ListOptions{})
	var result []BatchV1.Job
	for _, job := range jobs.Items {
		for _, reference := range job.OwnerReferences {
			if reference.Name == name {
				result = append(result, job)
			}
		}
	}
	if err != nil {
		panic(err.Error())
	}
	return &result
}

func runfzf(command string) string {
	fzf := exec.Command("fzf", "--ansi", "--no-preview")
	var out bytes.Buffer
	fzf.Stdin = os.Stdin
	fzf.Stderr = os.Stderr
	fzf.Stdout = &out
	// This is needed for * to not be expanded below.
	command = strings.Replace(command, "*", "\"*\"", -1)
	fzf.Env = append(os.Environ(),
		fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", command),
	)
	if err := fzf.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic(err)
		}
	}
	return strings.TrimSpace(out.String())
}

func getPodLogs(pod corev1.Pod, clientset *kubernetes.Clientset) string {
	podLogOpts := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()
	return str
}
