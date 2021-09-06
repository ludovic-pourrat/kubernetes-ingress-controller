package status

import (
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/internal/errors"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	gatewayapi_v1alpha1 "sigs.k8s.io/gateway-api/apis/v1alpha1"
)

type GatewayClassReasonType string

const reasonValidGatewayClass = "Valid"
const reasonInvalidGatewayClass = "Invalid"

// computeGatewayClassAdmittedCondition computes the GatewayClass Admitted status
// condition based on errs.
func computeGatewayClassAdmittedCondition(errs field.ErrorList) metav1.Condition {
	c := metav1.Condition{
		Type:    string(gatewayapi_v1alpha1.GatewayClassConditionStatusAdmitted),
		Status:  metav1.ConditionTrue,
		Reason:  reasonValidGatewayClass,
		Message: "Valid GatewayClass",
	}

	if errs != nil {
		c.Status = metav1.ConditionFalse
		c.Reason = reasonInvalidGatewayClass
		c.Message = fmt.Sprintf("Invalid GatewayClass: %s.", errors.ParseFieldErrors(errs))
	}

	return c
}

func mergeConditions(conditions []metav1.Condition, updates ...metav1.Condition) []metav1.Condition {
	now := metav1.NewTime(clock.Now())
	var additions []metav1.Condition
	for i, update := range updates {
		add := true
		for j, cond := range conditions {
			if cond.Type == update.Type {
				add = false
				if conditionChanged(cond, update) {
					conditions[j].Status = update.Status
					conditions[j].Reason = update.Reason
					conditions[j].Message = update.Message
					if cond.Status != update.Status {
						conditions[j].LastTransitionTime = now
					}
					break
				}
			}
		}
		if add {
			updates[i].LastTransitionTime = now
			additions = append(additions, updates[i])
		}
	}
	conditions = append(conditions, additions...)
	return conditions
}

func conditionChanged(a, b metav1.Condition) bool {
	return a.Status != b.Status || a.Reason != b.Reason || a.Message != b.Message
}

func conditionsEqual(a, b []metav1.Condition) bool {
	return apiequality.Semantic.DeepEqual(a, b)
}
