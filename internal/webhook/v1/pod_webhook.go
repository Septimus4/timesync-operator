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

package v1

import (
	"context"
	"fmt"
	syncv1alpha1 "github.com/Septimus4/timesync-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// SetupPodWebhookWithManager registers the webhook for Pod in the manager.
func SetupPodWebhookWithManager(mgr ctrl.Manager) error {
	k8sClient = mgr.GetClient()
	return ctrl.NewWebhookManagedBy(mgr).For(&corev1.Pod{}).
		WithDefaulter(&PodCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate--v1-pod,mutating=true,failurePolicy=fail,sideEffects=None,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod-v1.kb.io,admissionReviewVersions=v1

// PodCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Pod when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodCustomDefaulter struct {
}

var _ webhook.CustomDefaulter = &PodCustomDefaulter{}
var k8sClient client.Client

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Pod.
func (d *PodCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod object but got %T", obj)
	}

	logger := logf.FromContext(ctx)
	logger.Info("Webhook triggered for Pod", "name", pod.GetName(), "namespace", pod.GetNamespace())

	for _, c := range pod.Spec.Containers {
		if c.Name == "timesync" {
			logger.Info("Timesync sidecar already present; skipping")
			return nil
		}
	}

	ns := &corev1.Namespace{}
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: pod.Namespace}, ns); err != nil {
		logger.Error(err, "Failed to get namespace")
		return nil
	}

	policies := &syncv1alpha1.TimeSyncPolicyList{}
	if err := k8sClient.List(ctx, policies); err != nil {
		logger.Error(err, "Failed to list TimeSyncPolicies")
		return nil
	}

	for _, policy := range policies.Items {
		selector, err := metav1.LabelSelectorAsSelector(&policy.Spec.NamespaceSelector)
		if err != nil {
			continue
		}

		if selector.Matches(labels.Set(ns.Labels)) && policy.Spec.Enable {
			logger.Info("Injecting timesync sidecar from policy", "policy", policy.Name)

			sidecar := corev1.Container{
				Name:  "timesync",
				Image: policy.Spec.Image,
				Args:  []string{"sleep", "infinity"},
			}

			pod.Spec.Containers = append(pod.Spec.Containers, sidecar)
			break
		}
	}

	return nil
}
