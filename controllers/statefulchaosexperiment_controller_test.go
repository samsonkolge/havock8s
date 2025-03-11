package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	chaosv1alpha1 "github.com/statefulchaos/statefulchaos/api/v1alpha1"
)

var _ = Describe("StatefulChaosExperiment Controller", func() {
	const (
		ExperimentName      = "test-experiment"
		ExperimentNamespace = "default"
		Timeout             = time.Second * 10
		Interval            = time.Millisecond * 250
	)

	Context("When creating a StatefulChaosExperiment", func() {
		It("Should initialize the experiment status", func() {
			By("Creating a new StatefulChaosExperiment")
			ctx := context.Background()
			experiment := &chaosv1alpha1.StatefulChaosExperiment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chaos.statefulchaos.io/v1alpha1",
					Kind:       "StatefulChaosExperiment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExperimentName,
					Namespace: ExperimentNamespace,
				},
				Spec: chaosv1alpha1.StatefulChaosExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						Name:       "test-statefulset",
						TargetType: "StatefulSet",
					},
					ChaosType: "PodFailure",
					Duration:  "5m",
					Intensity: 0.5,
				},
			}
			Expect(k8sClient.Create(ctx, experiment)).Should(Succeed())

			// Look up the experiment after creation
			experimentLookupKey := types.NamespacedName{Name: ExperimentName, Namespace: ExperimentNamespace}
			createdExperiment := &chaosv1alpha1.StatefulChaosExperiment{}

			// We'll need to retry getting this newly created experiment, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, experimentLookupKey, createdExperiment)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// The experiment should have been initialized with Pending status
			Eventually(func() string {
				err := k8sClient.Get(ctx, experimentLookupKey, createdExperiment)
				if err != nil {
					return ""
				}
				return createdExperiment.Status.Phase
			}, Timeout, Interval).Should(Equal("Pending"))

			// The experiment should have a start time
			Expect(createdExperiment.Status.StartTime).ShouldNot(BeNil())

			// Clean up
			Expect(k8sClient.Delete(ctx, experiment)).Should(Succeed())
		})
	})

	Context("When the experiment duration is reached", func() {
		It("Should complete the experiment", func() {
			By("Creating a new StatefulChaosExperiment with a short duration")
			ctx := context.Background()
			experiment := &chaosv1alpha1.StatefulChaosExperiment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "chaos.statefulchaos.io/v1alpha1",
					Kind:       "StatefulChaosExperiment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ExperimentName + "-short",
					Namespace: ExperimentNamespace,
				},
				Spec: chaosv1alpha1.StatefulChaosExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						Name:       "test-statefulset",
						TargetType: "StatefulSet",
					},
					ChaosType: "PodFailure",
					Duration:  "1s", // Very short duration for testing
					Intensity: 0.5,
				},
			}
			Expect(k8sClient.Create(ctx, experiment)).Should(Succeed())

			// Look up the experiment after creation
			experimentLookupKey := types.NamespacedName{Name: ExperimentName + "-short", Namespace: ExperimentNamespace}
			createdExperiment := &chaosv1alpha1.StatefulChaosExperiment{}

			// We'll need to retry getting this newly created experiment, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, experimentLookupKey, createdExperiment)
				return err == nil
			}, Timeout, Interval).Should(BeTrue())

			// The experiment should eventually complete
			Eventually(func() string {
				err := k8sClient.Get(ctx, experimentLookupKey, createdExperiment)
				if err != nil {
					return ""
				}
				return createdExperiment.Status.Phase
			}, Timeout, Interval).Should(Or(Equal("Completed"), Equal("Failed")))

			// The experiment should have an end time
			Expect(createdExperiment.Status.EndTime).ShouldNot(BeNil())

			// Clean up
			Expect(k8sClient.Delete(ctx, experiment)).Should(Succeed())
		})
	})
})
