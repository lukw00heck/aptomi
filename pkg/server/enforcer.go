package server

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/engine/apply"
	"github.com/Aptomi/aptomi/pkg/engine/diff"
	"github.com/Aptomi/aptomi/pkg/engine/progress"
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/event"
	"github.com/Aptomi/aptomi/pkg/lang"
	"github.com/Aptomi/aptomi/pkg/plugin"
	"github.com/Aptomi/aptomi/pkg/plugin/helm"
	"github.com/Aptomi/aptomi/pkg/runtime"
	log "github.com/Sirupsen/logrus"
	"runtime/debug"
	"time"
)

func logError(err interface{}) {
	log.Errorf("Error while enforcing policy: %s", err)

	// todo make configurable
	debug.PrintStack()
}

func (server *Server) enforceLoop() error {
	for {
		err := server.enforce()
		if err != nil {
			logError(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (server *Server) enforce() error {
	defer func() {
		if err := recover(); err != nil {
			logError(err)
		}
	}()

	desiredPolicy, desiredPolicyGen, err := server.store.GetPolicy(runtime.LastGen)
	if err != nil {
		return fmt.Errorf("error while getting desiredPolicy: %s", err)
	}

	// skip policy enforcement if no policy found
	if desiredPolicy == nil {
		//todo log
		return nil
	}

	actualState, err := server.store.GetActualState()
	if err != nil {
		return fmt.Errorf("error while getting actual state: %s", err)
	}

	resolver := resolve.NewPolicyResolver(desiredPolicy, server.externalData)
	desiredState, eventLog, err := resolver.ResolveAllDependencies()
	if err != nil {
		return fmt.Errorf("cannot resolve desiredPolicy: %v %v %v", err, desiredState, actualState)
	}

	eventLog.Save(&event.HookConsole{})

	nextRevision, err := server.store.NewRevision(desiredPolicyGen)
	if err != nil {
		return fmt.Errorf("unable to get next revision: %s", err)
	}

	stateDiff := diff.NewPolicyResolutionDiff(desiredState, actualState, nextRevision.GetGeneration())

	// todo add check that policy gen not changed (always create new revision if policy gen changed)
	if len(stateDiff.Actions) <= 0 {
		// todo
		log.Infof("No changes")
		return nil
	}
	//todo
	log.Infof("Changes")
	// todo if policy gen changed, we still need to save revision but with progress == done

	//todo remove debug log
	for _, action := range stateDiff.Actions {
		fmt.Println(action)
	}

	// Save revision
	err = server.store.SaveRevision(nextRevision)
	if err != nil {
		return fmt.Errorf("error while saving new revision: %s", err)
	}

	// todo generate diagrams
	//prefDiagram := visualization.NewDiagram(actualPolicy, actualState, externalData)
	//newDiagram := visualization.NewDiagram(desiredPolicy, desiredState, externalData)
	//deltaDiagram := visualization.NewDiagramDelta(desiredPolicy, desiredState, actualPolicy, actualState, externalData)
	//visualization.CreateImage(...) for all diagrams

	// Build plugins
	helmIstio := helm.NewPlugin(server.cfg.Helm)
	plugins := plugin.NewRegistry(
		[]plugin.DeployPlugin{helmIstio},
		[]plugin.PostProcessPlugin{helmIstio},
	)

	actualPolicy, err := server.getActualPolicy()
	if err != nil {
		return fmt.Errorf("error while getting actual policy: %s", err)
	}

	applier := apply.NewEngineApply(desiredPolicy, desiredState, actualPolicy, actualState, server.store.GetActualStateUpdater(), server.externalData, plugins, stateDiff.Actions, progress.NewConsole())
	resolution, eventLog, err := applier.Apply()

	eventLog.Save(&event.HookConsole{})

	if err != nil {
		return fmt.Errorf("error while applying new revision: %s", err)
	}
	log.Infof("Applied new revision with resolution: %v", resolution)

	return nil
}

func (server *Server) getActualPolicy() (*lang.Policy, error) {
	currRevision, err := server.store.GetRevision(runtime.LastGen)
	if err != nil {
		return nil, fmt.Errorf("unable to get current revision: %s", err)
	}

	// it's just a first revision
	if currRevision == nil {
		return lang.NewPolicy(), nil
	}

	actualPolicy, _, err := server.store.GetPolicy(currRevision.Policy)
	if err != nil {
		return nil, fmt.Errorf("unable to get actual policy: %s", err)
	}

	return actualPolicy, nil
}
