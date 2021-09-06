package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/internal/proxy"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayapi_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

type GatewayReconciler struct {
	ctx    context.Context
	Client client.Client
	Log    logr.Logger

	// gatewayClassControllerName is the configured controller of managed gatewayclasses.
	gatewayClassControllerName string
	Scheme                     *runtime.Scheme
	Proxy                      proxy.Proxy
	GatewayClassName           string
}

// hasMatchingController returns true if the provided object is a Gateway
// using a GatewayClass with a Spec.Controller string matching this Contour's
// controller string, or false otherwise.
func (r *GatewayReconciler) hasMatchingController(obj client.Object) bool {
	log := r.Log.WithName(obj.GetName()).WithName(obj.GetNamespace())

	gw, ok := obj.(*gatewayapi_v1alpha1.Gateway)
	if !ok {
		log.Info("invalid object, bypassing reconciliation.")
		return false
	}

	matches, err := r.hasKongOwnedClass(gw)
	if err != nil {
		r.Log.Error(err, "error matching controller message")
		return false
	}
	if matches {
		log.Info("enqueueing gateway")
		return true
	}

	log.Info("configured controllerName doesn't match an existing GatewayClass")
	return false
}

func (r *GatewayReconciler) hasKongOwnedClass(gw *gatewayapi_v1alpha1.Gateway) (bool, error) {
	gc := &gatewayapi_v1alpha1.GatewayClass{}
	if err := r.Client.Get(r.ctx, types.NamespacedName{Name: gw.Spec.GatewayClassName}, gc); err != nil {
		return false, fmt.Errorf("failed to get gatewayclass %s: %w", gw.Spec.GatewayClassName, err)
	}
	if gc.Spec.Controller != r.gatewayClassControllerName {
		return false, nil
	}
	return true, nil
}

func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&gatewayapi_v1alpha1.GatewayClass{}).Complete(r)
}

func (r *GatewayReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	r.Log.WithName(request.Namespace).WithName(request.Name).Info("reconciling gateway")

	// Fetch the Gateway.
	gw := &gatewayapi_v1alpha1.Gateway{}
	if err := r.Client.Get(ctx, request.NamespacedName, gw); errors.IsNotFound(err) {
		// Not-found error, so trigger an OnDelete.
		r.Log.WithName(request.Name).WithName(request.Namespace).Info("failed to find gateway")

		fmt.Printf("Need delete the gateway object.")
	} else if err != nil {
		// Error reading the object, so requeue the request.
		return reconcile.Result{}, fmt.Errorf("failed to get gateway %s/%s: %w", request.Namespace, request.Name, err)
	}

	// TODO: Ensure the gateway by creating manage infrastructure, i.e. the Envoy service.
	// xref: https://github.com/projectcontour/contour/issues/3545

	// Pass the gateway off to the eventHandler.
	//r.eventHandler.OnAdd(gw)

	return reconcile.Result{}, nil
}
