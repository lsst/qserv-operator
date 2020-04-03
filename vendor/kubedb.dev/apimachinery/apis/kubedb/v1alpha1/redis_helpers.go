/*
Copyright The KubeDB Authors.

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

package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ Redis) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedis))
}

var _ apis.ResourceInfo = &Redis{}

func (r Redis) OffshootName() string {
	return r.Name
}

func (r Redis) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: r.Name,
		LabelDatabaseKind: ResourceKindRedis,
	}
}

func (r Redis) OffshootLabels() map[string]string {
	out := r.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularRedis
	out[meta_util.VersionLabelKey] = string(r.Spec.Version)
	out[meta_util.InstanceLabelKey] = r.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, r.Labels)
}

func (r Redis) ResourceShortCode() string {
	return ResourceCodeRedis
}

func (r Redis) ResourceKind() string {
	return ResourceKindRedis
}

func (r Redis) ResourceSingular() string {
	return ResourceSingularRedis
}

func (r Redis) ResourcePlural() string {
	return ResourcePluralRedis
}

func (r Redis) ServiceName() string {
	return r.OffshootName()
}

func (r Redis) ConfigMapName() string {
	return r.OffshootName()
}

func (r Redis) BaseNameForShard() string {
	return fmt.Sprintf("%s-shard", r.OffshootName())
}

func (r Redis) StatefulSetNameWithShard(i int) string {
	return fmt.Sprintf("%s%d", r.BaseNameForShard(), i)
}

type redisApp struct {
	*Redis
}

func (r redisApp) Name() string {
	return r.Redis.Name
}

func (r redisApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedis))
}

func (r Redis) AppBindingMeta() appcat.AppBindingMeta {
	return &redisApp{&r}
}

type redisStatsService struct {
	*Redis
}

func (r redisStatsService) GetNamespace() string {
	return r.Redis.GetNamespace()
}

func (r redisStatsService) ServiceName() string {
	return r.OffshootName() + "-stats"
}

func (r redisStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", r.Namespace, r.Name)
}

func (r redisStatsService) Path() string {
	return DefaultStatsPath
}

func (r redisStatsService) Scheme() string {
	return ""
}

func (r Redis) StatsService() mona.StatsAccessor {
	return &redisStatsService{&r}
}

func (r Redis) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, r.OffshootSelectors(), r.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (r *Redis) GetMonitoringVendor() string {
	if r.Spec.Monitor != nil {
		return r.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (r *Redis) SetDefaults() {
	if r == nil {
		return
	}

	// perform defaulting
	if r.Spec.Mode == "" {
		r.Spec.Mode = RedisModeStandalone
	} else if r.Spec.Mode == RedisModeCluster {
		if r.Spec.Cluster == nil {
			r.Spec.Cluster = &RedisClusterSpec{}
		}
		if r.Spec.Cluster.Master == nil {
			r.Spec.Cluster.Master = types.Int32P(3)
		}
		if r.Spec.Cluster.Replicas == nil {
			r.Spec.Cluster.Replicas = types.Int32P(1)
		}
	}
	if r.Spec.StorageType == "" {
		r.Spec.StorageType = StorageTypeDurable
	}
	if r.Spec.UpdateStrategy.Type == "" {
		r.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if r.Spec.TerminationPolicy == "" {
		r.Spec.TerminationPolicy = TerminationPolicyDelete
	} else if r.Spec.TerminationPolicy == TerminationPolicyPause {
		r.Spec.TerminationPolicy = TerminationPolicyHalt
	}

	if r.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		r.Spec.PodTemplate.Spec.ServiceAccountName = r.OffshootName()
	}

	r.Spec.Monitor.SetDefaults()
}

func (e *RedisSpec) GetSecrets() []string {
	return nil
}
