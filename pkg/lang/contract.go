package lang

import (
	"github.com/Aptomi/aptomi/pkg/lang/expression"
	"github.com/Aptomi/aptomi/pkg/lang/template"
	"github.com/Aptomi/aptomi/pkg/runtime"
)

// ContractObject is an informational data structure with Kind and Constructor for Contract
var ContractObject = &runtime.Info{
	Kind:        "contract",
	Storable:    true,
	Versioned:   true,
	Constructor: func() runtime.Object { return &Contract{} },
}

// Contract is an object, which allows you to define a contract for a service, as well as a set of specific
// implementations. For example, contract can be a 'database', with specific service contexts represented
// by 'MySQL', 'MariaDB', 'SQLite'.
//
// When dependencies get declared, they always get declared on a contract (not on a specific service).
type Contract struct {
	runtime.TypeKind `yaml:",inline"`
	Metadata         `validate:"required"`

	// ChangeLabels defines how current set of labels will get changed/transformed in case
	// the contract gets matched
	ChangeLabels LabelOperations `yaml:"change-labels" validate:"labelOperations"`

	// Contexts contains an ordered list of contexts within a contract. When allocating an instance, Aptomi will pick
	// and instantiate the first context which matches the criteria
	Contexts []*Context `validate:"dive"`
}

// Context represents a single context within a service contract.
// It's essentially a service instance for a given of class of use cases, a given set of consumers, etc.
type Context struct {
	// Name defines context name in the policy
	Name string `validate:"identifier"`

	// Criteria - if it gets evaluated to true during policy resolution, then contract
	// will get fulfilled by allocating this service context. It's an optional field, so if it's nil then
	// it is considered to be evaluated to true automatically
	Criteria *Criteria `validate:"omitempty"`

	// ChangeLabels defines how current set of labels will get changed/transformed in case
	// the context gets matched
	ChangeLabels LabelOperations `yaml:"change-labels" validate:"labelOperations"`

	// Allocation defines how the context will get allocated (which service to allocate and which unique key to use)
	Allocation *Allocation `validate:"required"`
}

// Allocation determines which service should be allocated for by the given context
// and which additional keys should be added to component instance key
type Allocation struct {
	// Service defined which service to allocated. It can be in form of 'serviceName', referring to service within
	// current namespace. Or it can be in form of 'namespace/serviceName', referring to service in a different
	// namespace
	Service string `validate:"required"`

	// Keys define a set of unique keys that define this allocation. If keys are not defined, then allocation will
	// always correspond to a single instance. If keys are defined, it will allow to create different service instances
	// based on labels. Different keys values resolved during policy processing will result in different service
	// instances created by Aptomi. For example, if key is set to {{.User.Labels.team}}, it will get dynamically
	// resolved into a user's team name. And, since users from different teams will have different keys, every team
	// will get their own service instance from Aptomi
	Keys []string `validate:"dive,template"`
}

// Matches checks if context criteria is satisfied
func (context *Context) Matches(params *expression.Parameters, cache *expression.Cache) (bool, error) {
	if context.Criteria == nil {
		return true, nil
	}
	return context.Criteria.allows(params, cache)
}

// ResolveKeys resolves dynamic allocation keys, which later get added to component instance key
func (context *Context) ResolveKeys(params *template.Parameters, cache *template.Cache) ([]string, error) {
	if cache == nil {
		cache = template.NewCache()
	}
	// Resolve allocation keys (they can be dynamic, depending on user labels)
	result := []string{}
	for _, key := range context.Allocation.Keys {
		keyResolved, err := cache.Evaluate(key, params)
		if err != nil {
			return nil, err
		}
		result = append(result, keyResolved)
	}
	return result, nil
}
