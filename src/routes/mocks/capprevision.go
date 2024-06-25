package mocks

import (
	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCappRevision returns a mock CappRevision object.
func GetCappRevision(name, namespace string, labels, annotations map[string]string) cappv1alpha1.CappRevision {
	cappRevision := cappv1alpha1.CappRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: cappv1alpha1.CappRevisionSpec{
			RevisionNumber: 1,
			CappTemplate: cappv1alpha1.CappTemplate{
				Spec: getCappSpec(),
			},
		},
		Status: cappv1alpha1.CappRevisionStatus{},
	}

	return cappRevision
}

// GetCappRevisionType returns a mock CappRevision type object.
func GetCappRevisionType(name, namespace string, labels, annotations map[string]string) types.CappRevision {
	cappRevision := types.CappRevision{
		Annotations: utils.ConvertMapToKeyValue(annotations),
		Labels:      utils.ConvertMapToKeyValue(labels),
		Metadata: types.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cappv1alpha1.CappRevisionSpec{
			RevisionNumber: 1,
			CappTemplate: cappv1alpha1.CappTemplate{
				Spec: getCappSpec(),
			},
		},
		Status: cappv1alpha1.CappRevisionStatus{},
	}

	return cappRevision
}
