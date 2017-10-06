package cluster

import (
	"github.com/Aptomi/aptomi/pkg/slinga/engine/apply/action"
	"github.com/Aptomi/aptomi/pkg/slinga/object"
)

// PostProcessActionObject is an informational data structure with Kind and Constructor for the action
var PostProcessActionObject = &object.Info{
	Kind:        "action-clusters-post-process",
	Constructor: func() object.Base { return &PostProcessAction{} },
}

// PostProcessAction is a global post-processing action which gets called once after all components have been processed by the engine
type PostProcessAction struct {
	*action.Metadata
}

// NewClustersPostProcessAction creates new PostProcessAction
func NewClustersPostProcessAction(revision object.Generation) *PostProcessAction {
	return &PostProcessAction{
		Metadata: action.NewMetadata(revision, PostProcessActionObject.Kind),
	}
}

// GetName returns action name
func (a *PostProcessAction) GetName() string {
	return "Clusters post process"
}

// Apply applies the action
func (a *PostProcessAction) Apply(context *action.Context) error {
	return nil
}