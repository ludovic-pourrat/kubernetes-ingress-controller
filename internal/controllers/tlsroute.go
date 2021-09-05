package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/internal/proxy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gatewayapi_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

type TlsRouteReconciler struct {
	// jm : we do not need expose those details
	Client client.Client
	//eventHandler cache.ResourceEventHandler // this is more generic
	Log logr.Logger

	Scheme           *runtime.Scheme // why we set all those the same ???
	Proxy            proxy.Proxy     // why we set all those the same ???
	GatewayClassName string
}

// NewTLSRouteController creates the tlsroute controller from mgr. The controller will be pre-configured
// to watch for TLSRoute objects across all namespaces.
// func NewTLSRouteController(mgr manager.Manager, eventHandler cache.ResourceEventHandler, log logrus.FieldLogger) (controller.Controller, error) {
// 	r := &TlsRouteReconciler{
// 		client:       mgr.GetClient(),
// 		eventHandler: eventHandler,
// 		FieldLogger:  log,
// 	}
// 	c, err := controller.New("tlsroute-controller", mgr, controller.Options{Reconciler: r})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := c.Watch(&source.Kind{Type: &gatewayapi_v1alpha1.TLSRoute{}}, &handler.EnqueueRequestForObject{}); err != nil {
// 		return nil, err
// 	}
// 	return c, nil
// }

// SetupWithManager sets up the controller with the Manager.
func (r *TlsRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&gatewayapi_v1alpha1.TLSRoute{}).Complete(r)
}

func (r *TlsRouteReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {

	// Fetch the TLSRoute from the cache.
	tlsroute := &gatewayapi_v1alpha1.TLSRoute{}
	err := r.client.Get(ctx, request.NamespacedName, tlsroute)
	if errors.IsNotFound(err) {
		r.eventHandler.OnDelete(&gatewayapi_v1alpha1.TLSRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      request.Name,
				Namespace: request.Namespace,
			},
		})
		return reconcile.Result{}, nil
	}

	// Pass the new changed object off to the eventHandler.
	r.eventHandler.OnAdd(tlsroute)

	return reconcile.Result{}, nil
}
