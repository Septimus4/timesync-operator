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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	syncv1alpha1 "github.com/Septimus4/timesync-operator/api/v1alpha1"
)

var _ = Describe("TimeSyncPolicy Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		timesyncpolicy := &syncv1alpha1.TimeSyncPolicy{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind TimeSyncPolicy")
			err := k8sClient.Get(ctx, typeNamespacedName, timesyncpolicy)
			if err != nil && errors.IsNotFound(err) {
				resource := &syncv1alpha1.TimeSyncPolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			By("Cleanup the specific resource instance TimeSyncPolicy")
			resource := &syncv1alpha1.TimeSyncPolicy{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err != nil {
				if errors.IsNotFound(err) {
					return
				}
				Expect(err).NotTo(HaveOccurred())
			}
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &TimeSyncPolicyReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Verifying the status of the reconciled resource")
			err = k8sClient.Get(ctx, typeNamespacedName, timesyncpolicy)
			Expect(err).NotTo(HaveOccurred())
			Expect(timesyncpolicy.Status.MatchedNamespaces).To(BeNumerically(">=", 0))
		})
	})
})
