package sync

// NewNetworkPolicySyncer generates Network Policies specification for Qserv
// func NewNetworkPolicySyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
// 	cm := qserv.GenerateNetworkPolicy(r, controllerLabels)
// 	objectName := fmt.Sprintf("%s%sConfigMap", strings.Title(string(container)), strings.Title(subpath))
// 	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func(existing runtime.Object) error {
// 		return nil
// 	})
// }
