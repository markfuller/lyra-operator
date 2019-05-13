package workflow

import (
	"context"
	"fmt"
	"time"

	lyrav1alpha1 "github.com/lyraproj/lyra-operator/pkg/apis/lyra/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// AddUS adds a controller with a watch for UserScript and a secondary watch (not working) for WebServer
func AddUS(mgr manager.Manager, applicator Applicator) error {

	return addWatches(
		mgr,
		&ReconcileUserScript{
			client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		},
		"userscript-controller",
		&lyrav1alpha1.UserScript{},
		&lyrav1alpha1.UserScript{},
	)
}

// newReconciler returns a new reconcile.Reconciler
func newReconcilerLB(mgr manager.Manager, applicator Applicator) reconcile.Reconciler {
	return &ReconcileUserScript{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}

}

//ReconcileUserScript .
type ReconcileUserScript struct {
	client client.Client
	scheme *runtime.Scheme
}

//Reconcile .
func (r *ReconcileUserScript) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var log = logf.Log.WithName("loadbalancer.controller")
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling UserScript")

	// get the loadbalancer object
	instance := &lyrav1alpha1.UserScript{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("UserScript not found - probably deleted")
			return reconcile.Result{}, nil
		}
		reqLogger.Info("UserScript error trying to get object")
		return reconcile.Result{}, err
	}
	scriptName := instance.Spec.ScriptName
	if scriptName == "" {
		reqLogger.Info("Nothing to run, so we'll exit")
		return reconcile.Result{}, err
	}

	//TODO apply some validation here?

	reqLogger.Info("We are ready to create this UserScript in the real world")
	cmd := exec.Command(scriptName)
	output, err := cmd.Output()
	reqLogger.Info(fmt.Sprintf("output from script (%v) is \n%s\n", scriptName, output))

	// fail (incl. non-zero exit code), so requeue
	if err != nil {
		reqLogger.Error(err, "Problem running script, will requeue after 1 minute", "scriptName", scriptName)
		return reconcile.Result{RequeueAfter: 1 * time.Minute}, err
	}

	return reconcile.Result{}, nil
}

func handleScript(s lyrav1alpha1.UserScriptSpec) error {
	var err error
	if s.GitRepo != "" {
		err = runCommand("pwd")
		if err != nil {
			return err
		}
		err = runCommand("git")
		if err != nil {
			return err
		}
		err = runCommand("git", "--help")
		if err != nil {
			return err
		}
		err = runCommand("git", "clone", s.GitRepo)
		if err != nil {
			return err
		}
		if s.GitBranch != "" {
			err = runCommand("git", "checkout", s.GitBranch)
			if err != nil {
				return err
			}
		}
	}
	return runCommand(s.ScriptName)
}

func runCommand(s string, args ...string) error {
	var log = logf.Log.WithName("scripthandler")
	cmd := exec.Command(s, args...)
	output, err := cmd.Output()
	info := fmt.Sprintf("output from script (%v) is \n%s\n", s, output)
	log.Info(info)
	return err
}
