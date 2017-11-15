package api

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/engine/diff"
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/event"
	"github.com/Aptomi/aptomi/pkg/runtime"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (api *coreAPI) handlePolicyGet(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	gen := params.ByName("gen")

	if len(gen) == 0 {
		gen = strconv.Itoa(int(runtime.LastGen))
	}

	policyData, err := api.store.GetPolicyData(runtime.ParseGeneration(gen))
	if err != nil {
		panic(fmt.Sprintf("error while getting requested policy: %s", err))
	}

	if policyData == nil {
		api.contentType.WriteOneWithStatus(writer, request, nil, http.StatusNotFound)
	} else {
		api.contentType.WriteOne(writer, request, policyData)
	}
}

func (api *coreAPI) handlePolicyObjectGet(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	gen := params.ByName("gen")

	if len(gen) == 0 {
		gen = strconv.Itoa(int(runtime.LastGen))
	}

	policy, _, err := api.store.GetPolicy(runtime.ParseGeneration(gen))
	if err != nil {
		panic(fmt.Sprintf("error while getting requested policy: %s", err))
	}

	ns := params.ByName("ns")
	kind := params.ByName("kind")
	name := params.ByName("name")

	obj, err := policy.GetObject(kind, name, ns)
	if err != nil {
		panic(fmt.Sprintf("error while getting object %s/%s/%s in policy #%s", ns, kind, name, gen))
	}
	if obj == nil {
		api.contentType.WriteOneWithStatus(writer, request, nil, http.StatusNotFound)
	}

	api.contentType.WriteOne(writer, request, obj)
}

// PolicyUpdateResultObject is an informational data structure with Kind and Constructor for PolicyUpdateResult
var PolicyUpdateResultObject = &runtime.Info{
	Kind:        "policy-update-result",
	Constructor: func() runtime.Object { return &PolicyUpdateResult{} },
}

// PolicyUpdateResult represents results for the policy update request (estimated list of actions to be executed to
// update existing actual state to the desired state)
type PolicyUpdateResult struct {
	runtime.TypeKind `yaml:",inline"`
	PolicyGeneration runtime.Generation
	Actions          []string
}

func (api *coreAPI) handlePolicyUpdate(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	objects := api.readLang(request)

	user := api.getUserRequired(request)

	// Verify ACL for updated objects
	currentPolicy, currentPolicyGeneration, err := api.store.GetPolicy(runtime.LastGen)
	if err != nil {
		panic(fmt.Sprintf("Error while loading current policy: %s", err))
	}
	for _, obj := range objects {
		errAdd := currentPolicy.AddObject(obj)
		if errAdd != nil {
			panic(fmt.Sprintf("Error while adding updated object to policy: %s", errAdd))
		}
		errManage := currentPolicy.View(user).ManageObject(obj)
		if errManage != nil {
			panic(fmt.Sprintf("Error while adding updated object to policy: %s", errManage))
		}
	}

	err = currentPolicy.Validate()
	if err != nil {
		panic(fmt.Sprintf("Updated policy is invalid: %s", err))
	}

	// todo(slukjanov): handle deleted
	deleted := make([]runtime.Key, 0)
	changed, policyData, err := api.store.UpdatePolicy(objects, deleted, user.Name)
	if err != nil {
		panic(fmt.Sprintf("Error while updating policy: %s", err))
	}

	if !changed {
		api.contentType.WriteOne(writer, request, &PolicyUpdateResult{
			TypeKind:         PolicyUpdateResultObject.GetTypeKind(),
			PolicyGeneration: currentPolicyGeneration,
			Actions:          nil,
		})

		return
	}

	desiredPolicyGen := policyData.GetGeneration()
	desiredPolicy, _, err := api.store.GetPolicy(desiredPolicyGen)
	if err != nil {
		panic(fmt.Sprintf("Error while getting desiredPolicy: %s", err))
	}
	if desiredPolicy == nil {
		panic(fmt.Sprintf("Can't read policy right after updating it"))
	}

	actualState, err := api.store.GetActualState()
	if err != nil {
		panic(fmt.Sprintf("Error while getting actual state: %s", err))
	}

	// todo we should resolve before saving policy => add Mutex for this method to make sure it's safe
	// todo: add request id to the event log scope
	eventLog := event.NewLog("api-policy-update", true)
	resolver := resolve.NewPolicyResolver(desiredPolicy, api.externalData, eventLog)
	desiredState, err := resolver.ResolveAllDependencies()

	if err != nil {
		panic(fmt.Sprintf("Cannot resolve desiredPolicy: %v %v %v", err, desiredState, actualState))
	}

	stateDiff := diff.NewPolicyResolutionDiff(desiredState, actualState)

	actions := make([]string, len(stateDiff.Actions))
	for idx, action := range stateDiff.Actions {
		actions[idx] = action.GetName()
	}

	api.contentType.WriteOne(writer, request, &PolicyUpdateResult{
		TypeKind:         PolicyUpdateResultObject.GetTypeKind(),
		PolicyGeneration: desiredPolicyGen,
		Actions:          actions,
	})
}
