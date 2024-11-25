package controller

import (
	"context"
	tsv1alpha1 "github.com/akyriako/typesense-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Definitions to manage status conditions
const (
	ConditionTypeReady                = "Ready"
	ConditionTypeTypesenseQuorumReady = "QuorumReady"

	ConditionReasonReconciliationInProgress = "ReconciliationInProgress"
	ConditionReasonSecretNotReady           = "SecretNotReady"
	ConditionReasonConfigMapNotReady        = "ConfigMapNotReady"
	ConditionReasonServicesNotReady         = "ServicesNotReady"
	ConditionReasonQuorumReady              = "QuorumReady"
	ConditionReasonQuorumNotReady           = "QuorumNotReady"
	ConditionReasonQuorumDegraded           = "QuorumDegraded"
	ConditionReasonQuorumNeedsAttention     = "QuorumNeedsAttention"
	ConditionReasonStatefulSetNotReady      = "StatefulSetNotReady"

	InitReconciliationMessage = "Starting reconciliation"
	UpdateStatusMessageFailed = "failed to update typesense cluster status"
)

func (r *TypesenseClusterReconciler) initConditions(ctx context.Context, ts *tsv1alpha1.TypesenseCluster) error {
	if ts.Status.Conditions == nil || len(ts.Status.Conditions) == 0 {
		if err := r.patchStatus(ctx, ts, func(status *tsv1alpha1.TypesenseClusterStatus) {
			meta.SetStatusCondition(&ts.Status.Conditions, metav1.Condition{Type: ConditionTypeReady, Status: metav1.ConditionUnknown, Reason: ConditionReasonReconciliationInProgress, Message: InitReconciliationMessage})
		}); err != nil {
			r.logger.Error(err, UpdateStatusMessageFailed)
			return err
		}
	}
	return nil
}

func (r *TypesenseClusterReconciler) setConditionNotReady(ctx context.Context, ts *tsv1alpha1.TypesenseCluster, reason string, err error) error {
	if err := r.patchStatus(ctx, ts, func(status *tsv1alpha1.TypesenseClusterStatus) {
		meta.SetStatusCondition(&ts.Status.Conditions, metav1.Condition{Type: ConditionTypeReady, Status: metav1.ConditionFalse, Reason: reason, Message: err.Error()})
	}); err != nil {
		return err
	}
	return nil
}

func (r *TypesenseClusterReconciler) setConditionReady(ctx context.Context, ts *tsv1alpha1.TypesenseCluster, reason string) error {
	if err := r.patchStatus(ctx, ts, func(status *tsv1alpha1.TypesenseClusterStatus) {
		meta.SetStatusCondition(&ts.Status.Conditions, metav1.Condition{Type: ConditionTypeReady, Status: metav1.ConditionTrue, Reason: reason, Message: "Cluster is Ready"})
	}); err != nil {
		return err
	}
	return nil
}
