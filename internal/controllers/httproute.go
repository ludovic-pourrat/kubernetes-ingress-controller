package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/internal/proxy"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayapi_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

type HttpRouteReconciler struct {
	Client client.Client
	//eventHandler cache.ResourceEventHandler  // what really matters is how we handle the different resources
	Log logr.Logger

	Scheme           *runtime.Scheme // why we set all those the same ???
	Proxy            proxy.Proxy     // why we set all those the same ???
	GatewayClassName string
}

// jm: should use this
// NewHTTPRouteController creates the httproute controller from mgr. The controller will be pre-configured
// to watch for HTTPRoute objects across all namespaces.
// func NewHTTPRouteController(mgr manager.Manager, eventHandler cache.ResourceEventHandler, log logrus.FieldLogger) (controller.Controller, error) {
// 	r := &HttpRouteReconciler{
// 		client:       mgr.GetClient(),
// 		eventHandler: eventHandler,
// 		FieldLogger:  log,
// 	}
// 	c, err := controller.New("httproute-controller", mgr, controller.Options{Reconciler: r})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := c.Watch(&source.Kind{Type: &gatewayapi_v1alpha1.HTTPRoute{}}, &handler.EnqueueRequestForObject{}); err != nil {
// 		return nil, err
// 	}
// 	return c, nil
// }

func (r *HttpRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&gatewayapi_v1alpha1.HTTPRoute{}).Complete(r)
}

func (r *HttpRouteReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	httpRoute := &gatewayapi_v1alpha1.HTTPRoute{}
	err := r.Client.Get(ctx, request.NamespacedName, httpRoute)
	if errors.IsNotFound(err) {
		fmt.Printf("add httproute deletion logic here !")
	}

	fmt.Printf("First Pass, we should detect the HTTP Routes.")
	// Pass the new changed object off to the eventHandler.
	//r.eventHandler.OnAdd(httpRoute)
	fmt.Printf("add httproute process logic here.")

	return reconcile.Result{}, nil
}
