package kjClient

import (
	"flag"
	"fmt"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Dummy() {
	fmt.Println("hello")
}

type KjcResult struct {
	AppliedJob *batchv1.Job
}

func Involk(manifests ...string) ([]KjcResult, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return []KjcResult{}, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return []KjcResult{}, err
	}

	var results = make([]KjcResult, 0, len(manifests))
	var eg errgroup.Group
	for _, manifest := range manifests {
		eg.Go(func() error {
			decode := scheme.Codecs.UniversalDeserializer().Decode
			object, _, err := decode([]byte(manifest), nil, nil)
			if err != nil {
				return err
			}
			switch obj := object.(type) {
			case *batchv1.Job:
				if obj.GetNamespace() == "" {
					obj.ObjectMeta.Namespace = "default"
				}
				var result KjcResult
				jobsClient := clientset.BatchV1().Jobs(obj.GetNamespace())
				result.AppliedJob, err = jobsClient.Create(obj)
				if err != nil {
					return err
				}
				results = append(results, result)
				return nil
			default:
				return err
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return []KjcResult{}, err
	}

	return results, nil
}
