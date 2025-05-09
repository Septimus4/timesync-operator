package controller

import (
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HasTimeSyncSidecar", func() {
	It("returns true when sidecar present", func() {
		pod := &corev1.Pod{}
		pod.Spec.Containers = append(pod.Spec.Containers,
			corev1.Container{Name: "timesync"},
		)
		Expect(HasTimeSyncSidecar(pod)).To(BeTrue())
	})

	It("returns false otherwise", func() {
		pod := &corev1.Pod{}
		Expect(HasTimeSyncSidecar(pod)).To(BeFalse())
	})
})
