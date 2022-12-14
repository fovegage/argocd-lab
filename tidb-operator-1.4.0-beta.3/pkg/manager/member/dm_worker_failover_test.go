// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package member

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	"github.com/pingcap/tidb-operator/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

func TestWorkerFailoverFailover(t *testing.T) {
	tests := []struct {
		name     string
		update   func(*v1alpha1.DMCluster)
		err      bool
		expectFn func(t *testing.T, dc *v1alpha1.DMCluster)
	}{
		{
			name: "normal",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-1": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-1",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-2": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-2",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(2))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())

			},
		},
		{
			name: "dm-worker stage is not Offline",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"1": {Stage: v1alpha1.DMWorkerStateBound, Name: "dm-worker-1"},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(0))
				g.Expect(dc.Status.Worker.FailoverUID).To(BeEmpty())

			},
		},
		{
			name: "deadline not exceed",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-1": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-1",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-30 * time.Minute)},
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(0))
				g.Expect(dc.Status.Worker.FailoverUID).To(BeEmpty())

			},
		},
		{
			name: "lastTransitionTime is zero",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-1": {
						Stage: v1alpha1.DMWorkerStateOffline,
						Name:  "dm-worker-1",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(0))
				g.Expect(dc.Status.Worker.FailoverUID).To(BeEmpty())

			},
		},
		{
			name: "exist in failureStores",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-1": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-1",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
				}
				dc.Status.Worker.FailureMembers = map[string]v1alpha1.WorkerFailureMember{
					"dm-worker-1": {
						PodName: "dm-worker-1",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(1))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())
			},
		},
		{
			name: "not exceed max failover count",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-0": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-0",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-4": {
						Stage:              v1alpha1.DMWorkerStateFree,
						Name:               "dm-worker-4",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-5": {
						Stage:              v1alpha1.DMWorkerStateFree,
						Name:               "dm-worker-5",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
				}
				dc.Status.Worker.FailureMembers = map[string]v1alpha1.WorkerFailureMember{
					"dm-worker-1": {
						PodName: "dm-worker-1",
					},
					"dm-worker-2": {
						PodName: "dm-worker-2",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(3))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())
			},
		},
		{
			name: "exceed max failover count1",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-3": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-3",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-4": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-4",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-5": {
						Stage:              v1alpha1.DMWorkerStateFree,
						Name:               "dm-worker-5",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
				}
				dc.Status.Worker.FailureMembers = map[string]v1alpha1.WorkerFailureMember{
					"dm-worker-1": {
						PodName: "dm-worker-1",
					},
					"dm-worker-2": {
						PodName: "dm-worker-2",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(3))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())
			},
		},
		{
			name: "exceed max failover count2",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-0": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-0",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-4": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-4",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
					"dm-worker-5": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-5",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
				}
				dc.Status.Worker.FailureMembers = map[string]v1alpha1.WorkerFailureMember{
					"dm-worker-1": {
						PodName: "dm-worker-1",
					},
					"dm-worker-2": {
						PodName: "dm-worker-2",
					},
					"dm-worker-3": {
						PodName: "dm-worker-3",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(3))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())

			},
		},
		{
			name: "exceed max failover count2 but maxFailoverCount = 0",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Spec.Worker.MaxFailoverCount = pointer.Int32Ptr(0)
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-12": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-12",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-13": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-13",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
					"dm-worker-14": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-14",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
				}
				dc.Status.Worker.FailureMembers = map[string]v1alpha1.WorkerFailureMember{
					"dm-worker-1": {
						PodName: "dm-worker-1",
					},
					"dm-worker-2": {
						PodName: "dm-worker-2",
					},
					"dm-worker-3": {
						PodName: "dm-worker-3",
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(3))
				g.Expect(dc.Status.Worker.FailoverUID).To(BeEmpty())
			},
		},
		{
			name: "already have failoverUID",
			update: func(dc *v1alpha1.DMCluster) {
				dc.Status.Worker.FailoverUID = "failover-uid-test"
				dc.Status.Worker.Members = map[string]v1alpha1.WorkerMember{
					"dm-worker-1": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-1",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-70 * time.Minute)},
					},
					"dm-worker-2": {
						Stage:              v1alpha1.DMWorkerStateOffline,
						Name:               "dm-worker-2",
						LastTransitionTime: metav1.Time{Time: time.Now().Add(-61 * time.Minute)},
					},
				}
			},
			err: false,
			expectFn: func(t *testing.T, dc *v1alpha1.DMCluster) {
				g := NewGomegaWithT(t)
				g.Expect(len(dc.Status.Worker.FailureMembers)).To(Equal(2))
				g.Expect(dc.Status.Worker.FailoverUID).NotTo(BeEmpty())
				g.Expect(dc.Status.Worker.FailoverUID).To(Equal(types.UID("failover-uid-test")))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			dc := newDMClusterForMaster()
			dc.Spec.Worker.Replicas = 6
			dc.Spec.Worker.MaxFailoverCount = pointer.Int32Ptr(3)
			tt.update(dc)

			fakeDeps := controller.NewFakeDependencies()
			fakeDeps.CLIConfig.WorkerFailoverPeriod = 1 * time.Hour
			workerFailover := &workerFailover{deps: fakeDeps}

			err := workerFailover.Failover(dc)
			if tt.err {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
			tt.expectFn(t, dc)
		})
	}
}
