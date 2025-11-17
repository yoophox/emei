package ckube

import (
  "context"
  "fmt"
  "strings"

  "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "github.com/yolksys/emei/cfg/source/inter"
)

type ckube struct {
  s *v1.Service
}

// Load ...
func Load(p string) (inter.Source, error) {
  config, err := rest.InClusterConfig()
  if err != nil {
    return nil, fmt.Errorf("fail:cfg source loas, source:kube,err:inclusterconfig")
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    return nil, fmt.Errorf("fail:cfg source loas, source:kube,err:newforconfig")
  }
  p_ := strings.Split(p, "/")
  if len(p_) != 2 {
    return nil, fmt.Errorf("fail:cfg source load, source:kube, path:%s", p)
  }
  switch p_[0] {
  case "service:cfg":
    t, err := clientset.CoreV1().Services("").Get(context.Background(), "", metav1.GetOptions{})
    if err != nil {
      return nil, fmt.Errorf("fail:cfg source load, source:kube, err:get, par:{%s}", err.Error())
    }

    return &ckube{t}, nil
  }
  return nil, fmt.Errorf("fail:cfg source load, source:kube, err:path type, path:%s", p)
}

func (s *ckube) Read() (any, error) {
  return s.s, nil
}

func (s *ckube) Write(v any) error {
  return nil
}
