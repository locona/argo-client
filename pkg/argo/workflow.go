package argo

import (
	"log"

	v1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/k0kubun/pp"
	"github.com/locona/argo/workflow/validate"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func toPtr(s string) *string {
	return &s
}

func (a *argo) List() (*v1alpha1.WorkflowList, error) {
	listOpts := metav1.ListOptions{}
	list, err := a.client.List(listOpts)
	if err != nil {
		pp.Println(err)
	}
	return list, err
}

func (a *argo) Watch() (watch.Interface, error) {
	listOpts := metav1.ListOptions{ResourceVersion: "0"}
	wi, err := a.client.Watch(listOpts)
	if err != nil {
		pp.Println(err)
	}
	return wi, err
}

func (a *argo) Create() (*v1alpha1.Workflow, error) {
	dagTemplates := &v1alpha1.DAGTemplate{
		Tasks: []v1alpha1.DAGTask{
			v1alpha1.DAGTask{
				Name:     "A",
				Template: "echo",
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						v1alpha1.Parameter{
							Name:  "message",
							Value: toPtr("A"),
						},
					},
				},
			},
			v1alpha1.DAGTask{
				Name:     "B",
				Template: "echo",
				Dependencies: []string{
					"A",
				},
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						v1alpha1.Parameter{
							Name:  "message",
							Value: toPtr("B"),
						},
					},
				},
			},
			v1alpha1.DAGTask{
				Name:     "C",
				Template: "echo",
				Dependencies: []string{
					"A",
				},
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						v1alpha1.Parameter{
							Name:  "message",
							Value: toPtr("C"),
						},
					},
				},
			},
			v1alpha1.DAGTask{
				Name:     "D",
				Template: "echo",
				Dependencies: []string{
					"B",
					"C",
				},
				Arguments: v1alpha1.Arguments{
					Parameters: []v1alpha1.Parameter{
						v1alpha1.Parameter{
							Name:  "message",
							Value: toPtr("D"),
						},
					},
				},
			},
		},
	}
	templates := []v1alpha1.Template{
		v1alpha1.Template{
			Name: "echo",
			Inputs: v1alpha1.Inputs{
				Parameters: []v1alpha1.Parameter{
					v1alpha1.Parameter{
						Name: "message",
					},
				},
			},
			Container: &apiv1.Container{
				Image: "alpine:3.6",
				Command: []string{
					"echo",
					"{{inputs.parameters.message}}",
				},
			},
		},
		v1alpha1.Template{
			Name: "diamond",
			DAG:  dagTemplates,
		},
	}
	typeMeta := metav1.TypeMeta{
		Kind:       "Workflow",
		APIVersion: "argoproj.io/v1alpha1",
	}

	meta := metav1.ObjectMeta{
		GenerateName: "dag-diamond-",
	}
	wf := &v1alpha1.Workflow{
		TypeMeta:   typeMeta,
		ObjectMeta: meta,
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint: "diamond",
			Templates:  templates,
		},
	}
	err := validate.ValidateWorkflow(wf)
	if err != nil {
		log.Fatal(err)
	}
	r, err := a.client.Create(wf)
	pp.Println(r, err)
	return r, err
}
