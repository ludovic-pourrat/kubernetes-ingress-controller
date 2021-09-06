package controllers

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/internal/errors"
	"github.com/kong/kubernetes-ingress-controller/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/internal/status"
	"github.com/kong/kubernetes-ingress-controller/internal/validation"

	"github.com/go-logr/logr"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayapi_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

type GatewayClassReconciler struct {
	Client       client.Client
	eventHandler cache.ResourceEventHandler
	Log          logr.Logger
	Controller   string
	Scheme       *runtime.Scheme
	Proxy        proxy.Proxy
}

func (r *GatewayClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&gatewayapi_v1alpha1.GatewayClass{}).Complete(r)
}

// hasMatchingController returns true if the provided object is a GatewayClass
// with a Spec.Controller string matching this Contour's controller string,
// or false otherwise.
func (r *GatewayClassReconciler) hasMatchingController(obj client.Object) bool {

	log := r.Log.WithName(obj.GetName())

	gc, ok := obj.(*gatewayapi_v1alpha1.GatewayClass)
	if !ok {
		log.Info("invalid object, bypassing reconciliation.")
		return false
	}

	if gc.Spec.Controller == r.Controller {
		log.Info("persisting gatewayclass")
		return true
	}

	log.Info("controller is not %s; bypassing reconciliation", r.Controller)
	return false
}

func (r *GatewayClassReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	r.Log.WithName(request.Name).Info("reconciling gatewayclass")

	// Fetch the Gateway from the cache.
	gc := &gatewayapi_v1alpha1.GatewayClass{}
	if err := r.Client.Get(ctx, request.NamespacedName, gc); api_errors.IsNotFound(err) {
		// Not-found error, so trigger an OnDelete.
		r.Log.WithName(request.Name).Info("failed to find gatewayclass")

		fmt.Printf("delete the gatewayclass objects")
		return reconcile.Result{}, nil
	} else if err != nil {
		// Error reading the object, so requeue the request.
		return reconcile.Result{}, fmt.Errorf("failed to get gatewayclass %q: %w", request.Name, err)
	} else if !r.hasMatchingController(gc) {

		fmt.Printf("delete the gatewayclass objects because it does not belong to this controller.")
		return reconcile.Result{}, nil
	}

	// The gatewayclass is safe to process, so check if it's valid.
	errs := validation.ValidateGatewayClass(ctx, r.Client, gc, r.Controller)
	if errs != nil {
		r.Log.WithName(request.Name).Error(nil, "invalid gatewayclass: ", errors.ParseFieldErrors(errs))
	}

	// Pass the new changed object off to the eventHandler.
	r.eventHandler.OnAdd(gc)

	if err := status.SyncGatewayClass(ctx, r.Client, gc, errs); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to sync gatewayclass %q status: %w", gc.Name, err)
	}
	r.Log.WithName(request.Name).Info("synced gatewayclass status")

	return reconcile.Result{}, nil
}
