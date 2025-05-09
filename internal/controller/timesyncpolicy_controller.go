/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	syncv1alpha1 "github.com/Septimus4/timesync-operator/api/v1alpha1"
)

// TimeSyncPolicyReconciler reconciles a TimeSyncPolicy object
type TimeSyncPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sync.example.com,resources=timesyncpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sync.example.com,resources=timesyncpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=sync.example.com,resources=timesyncpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the TimeSyncPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *TimeSyncPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var policy syncv1alpha1.TimeSyncPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log.Info("Reconciling TimeSyncPolicy", "name", policy.Name)

	var namespaces corev1.NamespaceList
	if err := r.List(ctx, &namespaces); err != nil {
		return ctrl.Result{}, err
	}

	selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.NamespaceSelector)
	if err != nil {
		log.Error(err, "Invalid namespaceSelector")
		return ctrl.Result{}, nil
	}

	matchCount := 0
	for _, ns := range namespaces.Items {
		if selector.Matches(labels.Set(ns.Labels)) {
			matchCount++
		}
	}

	// Optional: update .status with the match count
	if policy.Status.MatchedNamespaces != matchCount {
		policy.Status.MatchedNamespaces = matchCount
		if err := r.Status().Update(ctx, &policy); err != nil {
			log.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
	}

	log.Info("TimeSyncPolicy reconciled", "matchedNamespaces", matchCount)
	return ctrl.Result{}, nil
}

// map a *Namespace event to the TimeSyncPolicies it matches
func (r *TimeSyncPolicyReconciler) mapNamespaceToPolicies(
	ctx context.Context,
	obj client.Object,
) []reconcile.Request {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return nil
	}

	var policies syncv1alpha1.TimeSyncPolicyList
	if err := r.List(ctx, &policies); err != nil {
		return nil
	}

	var reqs []reconcile.Request
	for _, policy := range policies.Items {
		selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.NamespaceSelector)
		if err != nil {
			continue
		}
		if selector.Matches(labels.Set(ns.Labels)) {
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{Name: policy.Name},
			})
		}
	}
	return reqs
}

// SetupWithManager wires the controller
func (r *TimeSyncPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&syncv1alpha1.TimeSyncPolicy{}).
		Watches(
			&corev1.Namespace{},
			handler.TypedEnqueueRequestsFromMapFunc[client.Object](r.mapNamespaceToPolicies),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}
