package csidriveroperator

import (
	"context"
	"io/fs"
	"os"
	"testing"
	"time"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	"github.com/openshift/cluster-storage-operator/pkg/csoclients"
	"github.com/openshift/cluster-storage-operator/pkg/operator/csidriveroperator/csioperatorclient"
	"github.com/openshift/library-go/pkg/operator/configobserver/featuregates"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/stretchr/testify/assert"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/errors"
	clocktesting "k8s.io/utils/clock/testing"
)

type RunControllerTest struct {
	name           string
	platformStatus *v1.PlatformStatus
	featureGate    featuregates.FeatureGate
	csiDriver      *storagev1.CSIDriver
	config         csioperatorclient.CSIOperatorConfig
	isInstalled    bool
	expectRun      bool
	expectError    bool
}

func TestShouldRunController(t *testing.T) {
	testingDefault := featuregates.NewFeatureGate(nil, []v1.FeatureGateName{features.FeatureGateExample})
	testingTechPreview := featuregates.NewFeatureGate([]v1.FeatureGateName{features.FeatureGateExample}, nil)
	customFeatureGate := featuregates.NewFeatureGate([]v1.FeatureGateName{"SomeOtherFeatureGate", features.FeatureGateExample, "YetAnotherGate"}, nil)
	customWithJustOther := featuregates.NewFeatureGate([]v1.FeatureGateName{"SomeOtherFeatureGate"}, nil)
	customWithNothing := featuregates.NewFeatureGate([]v1.FeatureGateName{}, nil)

	tests := []RunControllerTest{
		{
			name:           "tech preview example driver on AllPlatforms type",
			platformStatus: &v1.PlatformStatus{Type: v1.AWSPlatformType},
			featureGate:    testingTechPreview,
			csiDriver:      nil,
			config: csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           csioperatorclient.AllPlatforms,
				RequireFeatureGate: features.FeatureGateExample,
			},
			isInstalled: false,
			expectRun:   true,
			expectError: false,
		},
		{
			"tech preview example driver on AWSPlatformType",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			testingTechPreview,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			true,
			false,
		},
		{
			"tech preview example driver on GCPPlatformType",
			&v1.PlatformStatus{Type: v1.GCPPlatformType},
			testingTechPreview,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.GCPPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			true,
			false,
		},
		{
			"tech preview example driver on GCPPlatformType",
			&v1.PlatformStatus{Type: v1.VSpherePlatformType},
			testingTechPreview,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.VSpherePlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			true,
			false,
		},
		{
			"GA CSI driver on matching platform",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			testingDefault,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "ebs.csi.aws.com",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: "",
			},
			false,
			true,
			false,
		},
		{
			"GA CSI driver on non-matching platform",
			&v1.PlatformStatus{Type: v1.GCPPlatformType},
			testingDefault,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "ebs.csi.aws.com",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: "",
			},
			false,
			false,
			false,
		},
		{
			"GA CSI driver with StatusFilter returning true",
			&v1.PlatformStatus{Type: v1.IBMCloudPlatformType},
			testingDefault,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "vpc.block.csi.ibm.io",
				Platform:           v1.IBMCloudPlatformType,
				RequireFeatureGate: "",
				StatusFilter: func(*v1.InfrastructureStatus, bool) bool {
					return true
				},
			},
			false,
			true,
			false,
		},
		{
			"GA CSI driver with StatusFilter returning false",
			&v1.PlatformStatus{Type: v1.IBMCloudPlatformType},
			testingDefault,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "vpc.block.csi.ibm.io",
				Platform:           v1.IBMCloudPlatformType,
				RequireFeatureGate: "",
				StatusFilter: func(*v1.InfrastructureStatus, bool) bool {
					return false
				},
			},
			false,
			false,
			false,
		},
		{
			"tech preview example driver with positive custom featureGate",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			customFeatureGate,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			true,
			false,
		},
		{
			"tech preview example driver with negative custom featureGate",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			customWithJustOther,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			false,
			false,
		},
		{
			"tech preview example driver with empty custom featureGate",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			customWithNothing,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			false,
			false,
		},
		{
			"tech preview example driver with nil custom featureGate",
			&v1.PlatformStatus{Type: v1.AWSPlatformType},
			customWithNothing,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "csi.test.openshift.io",
				Platform:           v1.AWSPlatformType,
				RequireFeatureGate: features.FeatureGateExample,
			},
			false,
			false,
			false,
		},
		{
			"Azure File driver should not run on Azure StackHub if not already installed",
			&v1.PlatformStatus{Type: v1.AzurePlatformType, Azure: &v1.AzurePlatformStatus{CloudName: v1.AzureStackCloud}},
			customWithNothing,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "file.csi.azure.com",
				Platform:           v1.AzurePlatformType,
				StatusFilter:       csioperatorclient.IsNotAzueStackCloud,
				RequireFeatureGate: "",
			},
			false,
			false,
			false,
		},
		{
			"Azure File driver should keep running on Azure StackHub if already installed",
			&v1.PlatformStatus{Type: v1.AzurePlatformType, Azure: &v1.AzurePlatformStatus{CloudName: v1.AzureStackCloud}},
			customWithNothing,
			nil,
			csioperatorclient.CSIOperatorConfig{
				CSIDriverName:      "file.csi.azure.com",
				Platform:           v1.AzurePlatformType,
				StatusFilter:       csioperatorclient.IsNotAzueStackCloud,
				RequireFeatureGate: "",
			},
			true,
			true,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			infra := NewTestInfra().WithStatus(test.platformStatus)

			res, err := shouldRunController(test.config, infra, test.featureGate, test.csiDriver, test.isInstalled)
			if res != test.expectRun {
				t.Errorf("Expected run %t, got %t", test.expectRun, res)
			}
			gotError := err != nil
			if gotError != test.expectError {
				t.Errorf("Expected error %t, got %t: %s", test.expectError, res, err)
			}
		})
	}
}

type TestInfra struct {
	infra *v1.Infrastructure
}

func NewTestInfra() *TestInfra {
	testInfra := TestInfra{}
	testInfra.infra = &v1.Infrastructure{
		Status: v1.InfrastructureStatus{},
	}

	return &testInfra
}

func (i *TestInfra) WithStatus(status ...*v1.PlatformStatus) *v1.Infrastructure {
	if len(status) > 0 {
		i.infra.Status.PlatformStatus = status[0]
	}

	return i.infra
}

func csiDriver(csiDriverName string, annotations map[string]string) *storagev1.CSIDriver {
	return &storagev1.CSIDriver{
		ObjectMeta: metav1.ObjectMeta{
			Name:        csiDriverName,
			Annotations: annotations,
		},
		Spec: storagev1.CSIDriverSpec{},
	}
}

func TestIsNoMatchError(t *testing.T) {
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "foos",
	}
	gvk := schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Foo",
	}
	gk := schema.GroupKind{
		Group: "",
		Kind:  "foo",
	}

	tests := []struct {
		name          string
		err           error
		expectNoMatch bool
	}{
		{
			name:          "no error",
			err:           nil,
			expectNoMatch: false,
		},
		{
			name:          "NoResourceMatch",
			err:           &meta.NoResourceMatchError{PartialResource: gvr},
			expectNoMatch: true,
		},
		{
			name:          "NoKindMatch",
			err:           &meta.NoKindMatchError{GroupKind: gk},
			expectNoMatch: true,
		},
		{
			name: "unrelated error",
			err: &meta.AmbiguousKindError{
				PartialKind:       gvk,
				MatchingResources: nil,
				MatchingKinds:     nil,
			},
			expectNoMatch: false,
		},
		{
			name:          "aggregated NoResourceMatch",
			err:           errors.NewAggregate([]error{&meta.NoResourceMatchError{PartialResource: gvr}, os.ErrPermission}),
			expectNoMatch: true,
		},
		{
			name:          "aggregated NoKindMatch",
			err:           errors.NewAggregate([]error{os.ErrPermission, &meta.NoKindMatchError{GroupKind: gk}}),
			expectNoMatch: true,
		},
		{
			name:          "aggregated unrelated errors",
			err:           errors.NewAggregate([]error{os.ErrPermission, os.ErrExist, fs.ErrClosed}),
			expectNoMatch: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ret := isNoMatchError(test.err)
			if ret != test.expectNoMatch {
				t.Errorf("expected %t, got %t", test.expectNoMatch, ret)
			}
		})
	}
}

func getInfrastructure(platformType v1.PlatformType) *v1.Infrastructure {
	return &v1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name: infraConfigName,
		},
		Status: v1.InfrastructureStatus{
			PlatformStatus: &v1.PlatformStatus{
				Type: platformType,
			},
		},
	}
}

func getDefaultFeatureGate() *v1.FeatureGate {
	return &v1.FeatureGate{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Spec: v1.FeatureGateSpec{},
	}
}

func TestStandAloneStarter(t *testing.T) {
	initialObjects := &csoclients.FakeTestObjects{}
	initialObjects.OperatorObjects = append(initialObjects.OperatorObjects, csoclients.GetCR())
	initialObjects.ConfigObjects = append(initialObjects.ConfigObjects, getInfrastructure(v1.AWSPlatformType))
	initialObjects.ConfigObjects = append(initialObjects.ConfigObjects, getDefaultFeatureGate())

	finish, cancel := context.WithCancel(context.TODO())
	defer cancel()

	clients := csoclients.NewFakeClients(initialObjects)

	storageInformer := clients.OperatorInformers.Operator().V1().Storages().Informer()
	storageInformer.GetStore().Add(csoclients.GetCR())

	infrInformer := clients.ConfigInformers.Config().V1().Infrastructures().Informer()
	infrInformer.GetStore().Add(getInfrastructure(v1.AWSPlatformType))
	testingDefault := featuregates.NewFeatureGate(nil, []v1.FeatureGateName{features.FeatureGateExample})

	csoclients.StartInformers(clients, finish.Done())
	csoclients.WaitForSync(clients, finish.Done())

	awsConfig := []csioperatorclient.CSIOperatorConfig{csioperatorclient.GetAWSEBSCSIOperatorConfig(false)}

	fc, standAloneStarter := NewStandaloneDriverStarter(clients,
		testingDefault,
		20*time.Minute,
		status.NewVersionGetter(),
		"",
		events.NewInMemoryRecorder(csiDriverControllerName, clocktesting.NewFakePassiveClock(time.Now())),
		awsConfig)

	if fc == nil {
		t.Errorf("error creating the controller, gotl for the controller")
	}
	if len(standAloneStarter.controllers) != 1 {
		t.Errorf("expected at least one configured controller got %d", len(standAloneStarter.controllers))
	}

	csiController := standAloneStarter.controllers[0]
	assert.Equal(t, csiController.operatorConfig.CSIDriverName, "ebs.csi.aws.com")
	assert.Equal(t, csiController.operatorConfig.DeploymentAsset, "csidriveroperators/aws-ebs/standalone/generated/apps_v1_deployment_aws-ebs-csi-driver-operator.yaml")
	assert.NotEmpty(t, csiController.operatorConfig.StaticAssets)
}
