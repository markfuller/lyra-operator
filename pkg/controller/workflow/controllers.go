package workflow

import (
	"context"

	lyrav1alpha1 "github.com/lyraproj/lyra-operator/pkg/apis/lyra/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// AddLB adds a controller with a watch for LoadBalancer and a secondary watch (not working) for WebServer
func AddLB(mgr manager.Manager, applicator Applicator) error {

	return addWatches(
		mgr,
		&ReconcileLoadBalancer{
			client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		},
		"loadbalancer-controller",
		&lyrav1alpha1.LoadBalancer{},
		&lyrav1alpha1.WebServer{},
	)
}

//ReconcileLoadBalancer .
type ReconcileLoadBalancer struct {
	client client.Client
	scheme *runtime.Scheme
	// applicator Applicator
}

//Reconcile .
func (r *ReconcileLoadBalancer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var log = logf.Log.WithName("loadbalancer.controller")
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling LoadBalancer")

	// get the loadbalancer object
	instance := &lyrav1alpha1.LoadBalancer{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("LoadBalancer not found - probably deleted")
			return reconcile.Result{}, nil
		}
		reqLogger.Info("LoadBalancer error trying to get object")
		return reconcile.Result{}, err
	}

	//how do we read another resource anyway?
	//is the web server id set
	//if no, requeue
	if instance.Spec.WebServerID == "" {
		reqLogger.Info("LoadBalancer's WebServerID is not set, so we'll exit")
		return reconcile.Result{}, err
	}

	//if so, act on this, for now just log the info
	reqLogger.Info("We are ready to create this LoadBalancer in the real world")

	return reconcile.Result{}, nil
}
