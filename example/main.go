package main

import (
	"fmt"

	kjClient "github.com/nnao45/client-k8s-job-go"
)

const manifest = `
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4`

func main() {
	results, err := kjClient.Involk(manifest)
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
