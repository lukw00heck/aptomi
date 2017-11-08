package lang

import (
	"github.com/Aptomi/aptomi/pkg/runtime"
)

// ClusterObject is an informational data structure with Kind and Constructor for Cluster
var ClusterObject = &runtime.Info{
	Kind:        "cluster",
	Storable:    true,
	Versioned:   true,
	Constructor: func() runtime.Object { return &Cluster{} },
}

// Cluster defines an individual cluster where containers get deployed.
// Various cloud providers are supported via setting a cluster type (k8s, Amazon ECS, GKE, etc).
type Cluster struct {
	runtime.TypeKind `yaml:",inline"`
	Metadata         `validate:"required"`

	// Type is a cluster type. Based on its type, the appropriate deployment plugin will be called to deploy containers.
	Type string `validate:"clustertype"`

	// Labels is a set of labels attached to the cluster
	Labels map[string]string `validate:"omitempty,labels"`

	// Config for a given cluster type
	Config ClusterConfig `validate:"required"`
}

// ClusterConfig defines config for a k8s cluster with Helm
type ClusterConfig struct {
	KubeContext     string `validate:"required"`
	TillerNamespace string `validate:"omitempty"`
	Namespace       string `validate:"required"`
}
