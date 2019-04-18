package kjClient

import (
	"flag"
	"fmt"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Dummy() {
	fmt.Println("hello")
}

func DummyClinet(manifests ...string) error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	_ = clientset

	for _, manifest := range manifests {
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode([]byte(manifest), nil, nil)
		if err != nil {
			return err
		}
		switch o := obj.(type) {
		case *batchv1.Job:
			fmt.Println(o)
		default:
			return err
		}
	}

	/*
		jobsClient := clientset.BatchV1().Jobs("default")
		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "demo-job",
				Namespace: "gitlab",
			},
			Spec: batchv1.JobSpec{
				Template: apiv1.PodTemplateSpec{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "demo",
								Image: "myimage",
							},
						},
					},
				},
			},
		}*/
	return nil
}
