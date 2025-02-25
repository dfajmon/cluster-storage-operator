// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/openshift/api/operator/v1"
	operatorv1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1"
	typedoperatorv1 "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	gentype "k8s.io/client-go/gentype"
)

// fakeEtcds implements EtcdInterface
type fakeEtcds struct {
	*gentype.FakeClientWithListAndApply[*v1.Etcd, *v1.EtcdList, *operatorv1.EtcdApplyConfiguration]
	Fake *FakeOperatorV1
}

func newFakeEtcds(fake *FakeOperatorV1) typedoperatorv1.EtcdInterface {
	return &fakeEtcds{
		gentype.NewFakeClientWithListAndApply[*v1.Etcd, *v1.EtcdList, *operatorv1.EtcdApplyConfiguration](
			fake.Fake,
			"",
			v1.SchemeGroupVersion.WithResource("etcds"),
			v1.SchemeGroupVersion.WithKind("Etcd"),
			func() *v1.Etcd { return &v1.Etcd{} },
			func() *v1.EtcdList { return &v1.EtcdList{} },
			func(dst, src *v1.EtcdList) { dst.ListMeta = src.ListMeta },
			func(list *v1.EtcdList) []*v1.Etcd { return gentype.ToPointerSlice(list.Items) },
			func(list *v1.EtcdList, items []*v1.Etcd) { list.Items = gentype.FromPointerSlice(items) },
		),
		fake,
	}
}
