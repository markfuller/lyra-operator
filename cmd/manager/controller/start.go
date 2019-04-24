package controller

import (
	"context"
	"fmt"

	"github.com/lyraproj/lyra-operator/pkg/apis"
	"github.com/lyraproj/lyra-operator/pkg/controller/workflow"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/ready"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// Start the Kubernetes controller running
func Start(namespace string, applicator workflow.Applicator) error {

	var log = logf.Log.WithName("start.controllers")

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	// Become the leader before proceeding
	leader.Become(context.TODO(), "lyra-operator-lock")

	r := ready.NewFileReady()
	err = r.Set()
	if err != nil {
		return fmt.Errorf("failed to get ready: %v", err)
	}
	defer r.Unset()

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{Namespace: namespace})
	if err != nil {
		return fmt.Errorf("failed to create manager: %v", err)
	}

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		return fmt.Errorf("failed to add scheme: %v", err)
	}

	// set up workflow controller
	if err := workflow.Add(mgr, applicator); err != nil {
		//only log the error so that we can continue if the kind is not found (i.e. if we haven't created the CRD in kube yet)
		log.Error(err, "failed to setup workflow controller")
	}

	// set up load balancer controller
	// TODO markf refactor signatures (e.g. wrong package here) if it eventually works
	if err := workflow.AddLB(mgr, applicator); err != nil {
		//only log the error so that we can continue if the kind is not found (i.e. if we haven't created the CRD in kube yet)
		log.Error(err, "failed to setup loadbalancer controller")
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		return fmt.Errorf("failed to start manager: %v", err)
	}

	log.Info("Completed controller setup")
	return nil
}
