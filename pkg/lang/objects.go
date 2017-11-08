package lang

import "github.com/Aptomi/aptomi/pkg/runtime"

var (
	PolicyObjects = []*runtime.Info{
		ServiceObject,
		ContractObject,
		DependencyObject,
		ClusterObject,
		RuleObject,
		ACLRuleObject,
	}

	policyObjectsMap = make(map[runtime.Kind]bool)
)

func init() {
	for _, obj := range PolicyObjects {
		policyObjectsMap[obj.Kind] = true
	}
}

func IsPolicyObject(obj runtime.Object) bool {
	return policyObjectsMap[obj.GetKind()]
}
