// Copyright PingCAP, Inc.
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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1alpha1 "github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=pingcap.com, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("backups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().Backups().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("backupschedules"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().BackupSchedules().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("dmclusters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().DMClusters().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("dataresources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().DataResources().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("restores"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().Restores().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("tidbclusters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().TidbClusters().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("tidbclusterautoscalers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().TidbClusterAutoScalers().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("tidbinitializers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().TidbInitializers().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("tidbmonitors"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().TidbMonitors().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("tidbngmonitorings"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Pingcap().V1alpha1().TidbNGMonitorings().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}