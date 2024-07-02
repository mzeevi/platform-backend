package mocks

import (
	"fmt"
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knativeapis "knative.dev/pkg/apis"
	knativev1 "knative.dev/serving/pkg/apis/serving/v1"
	knativev1beta1 "knative.dev/serving/pkg/apis/serving/v1beta1"
)

const (
	CappImage     = "ghcr.io/dana-team/capp-gin-app:v0.2.0"
	containerName = "capp-container"
	Domain        = "dana-team.io"
	Hostname      = "custom-capp"
)

// PrepareCapp returns a mock Capp object.
func PrepareCapp(name, namespace string, labels, annotations map[string]string) cappv1alpha1.Capp {
	cappRevision := cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappSpec(),
		Status: PrepareCappStatus(name, namespace),
	}

	return cappRevision
}

// PrepareCappWithHostname returns a mock Capp object with Hostname set in the spec.
func PrepareCappWithHostname(name, namespace string, labels, annotations map[string]string) cappv1alpha1.Capp {
	cappRevision := cappv1alpha1.Capp{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec:   PrepareCappSpecWithHostname(),
		Status: PrepareCappStatusWithHostname(name, namespace),
	}

	return cappRevision
}

// PrepareCappSpec returns a mock Capp spec.
func PrepareCappSpec() cappv1alpha1.CappSpec {
	return cappv1alpha1.CappSpec{
		ConfigurationSpec: knativev1.ConfigurationSpec{
			Template: knativev1.RevisionTemplateSpec{
				Spec: knativev1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: CappImage,
								Name:  containerName,
							},
						},
					},
				},
			},
		},
	}
}

// PrepareCappSpecWithHostname returns a mock Capp spec with Hostname set.
func PrepareCappSpecWithHostname() cappv1alpha1.CappSpec {
	return cappv1alpha1.CappSpec{
		ConfigurationSpec: knativev1.ConfigurationSpec{
			Template: knativev1.RevisionTemplateSpec{
				Spec: knativev1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: CappImage,
								Name:  containerName,
							},
						},
					},
				},
			},
		},
		RouteSpec: cappv1alpha1.RouteSpec{
			Hostname: Hostname + "." + Domain,
		},
	}
}

// PrepareCappStatus returns a mock Capp status.
func PrepareCappStatus(name, namespace string) cappv1alpha1.CappStatus {
	return cappv1alpha1.CappStatus{
		KnativeObjectStatus: knativev1.ServiceStatus{
			RouteStatusFields: knativev1.RouteStatusFields{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s-%s.%s", name, namespace, Domain)),
			},
		},
	}
}

// PrepareCappStatusWithHostname returns a mock Capp status with Hostname set.
func PrepareCappStatusWithHostname(name, namespace string) cappv1alpha1.CappStatus {
	return cappv1alpha1.CappStatus{
		KnativeObjectStatus: knativev1.ServiceStatus{
			RouteStatusFields: knativev1.RouteStatusFields{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s-%s.%s", name, namespace, Domain)),
			},
		},
		RouteStatus: cappv1alpha1.RouteStatus{
			DomainMappingObjectStatus: knativev1beta1.DomainMappingStatus{
				URL: knativeapis.HTTPS(fmt.Sprintf("%s.%s", Hostname, Domain)),
			},
		},
	}
}

// PrepareCappType returns a Capp type object.
func PrepareCappType(name, namespace string, labels, annotations map[string]string) types.Capp {
	cappRevision := types.Capp{
		Annotations: utils.ConvertMapToKeyValue(annotations),
		Labels:      utils.ConvertMapToKeyValue(labels),
		Metadata: types.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   PrepareCappSpec(),
		Status: PrepareCappStatus(name, namespace),
	}

	return cappRevision
}

// PrepareCappTypeWithoutStatus returns a Capp type object without the Status set.
func PrepareCappTypeWithoutStatus(name, namespace string, labels, annotations map[string]string) types.Capp {
	return types.Capp{
		Annotations: utils.ConvertMapToKeyValue(annotations),
		Labels:      utils.ConvertMapToKeyValue(labels),
		Metadata: types.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   PrepareCappSpec(),
		Status: cappv1alpha1.CappStatus{},
	}
}