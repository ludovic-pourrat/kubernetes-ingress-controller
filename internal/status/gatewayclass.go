package status

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gateway_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

// SyncGatewayClass computes the current status of GatewayClass and updates status upon
// any changes since last sync.
func SyncGatewayClass(ctx context.Context, cli client.Client, gc *gateway_v1alpha1.GatewayClass, errs field.ErrorList) error {
	updated := gc.DeepCopy()

	updated.Status.Conditions = mergeConditions(updated.Status.Conditions, computeGatewayClassAdmittedCondition(errs))

	if !conditionsEqual(gc.Status.Conditions, updated.Status.Conditions) {
		if err := cli.Status().Update(ctx, updated); err != nil {
			return fmt.Errorf("failed to update gatewayclass %s status: %w", gc.Name, err)
		}
	}

	return nil
}
