package vultr

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/vultr/govultr/v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInstancesV2_NodeInstanceAddressesRejectsEmptyInstance(t *testing.T) {
	instances := &instancesv2{}

	for _, instance := range []*govultr.Instance{nil, {}} {
		_, err := instances.nodeInstanceAddresses(instance)
		if err == nil {
			t.Errorf("expected an error for empty instance %v", instance)
		} else if !strings.Contains(err.Error(), "instance is empty") {
			t.Errorf("expected empty instance error, got %v", err)
		}
	}
}

func TestInstancesV2_NodeBareMetalAddressesRejectsEmptyServer(t *testing.T) {
	instances := &instancesv2{}

	for _, server := range []*govultr.BareMetalServer{nil, {}} {
		_, err := instances.nodeBareMetalAddresses(server)
		if err == nil {
			t.Errorf("expected an error for empty bare metal server %v", server)
		} else if !strings.Contains(err.Error(), "baremetal is empty") {
			t.Errorf("expected empty bare metal server error, got %v", err)
		}
	}
}

func TestInstances_InstanceExistsByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddressesByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_InstanceShutdownByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceShutdownByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual {
		t.Errorf("expcted %+v got %+v", "false", "true")
	}
}

func TestInstances_InstanceTypeByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceTypeByProviderID(context.TODO(), "vultr://576965")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "vc2-4c-8gb" {
		t.Errorf("expcted %+v got %+v", "204", actual)
	}
}

func TestInstances_NodeAddressesByProviderID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddresses(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_NodeAddresses(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	expected := []v1.NodeAddress{
		{
			Type:    v1.NodeHostName,
			Address: "ccm-test",
		},
		{
			Type:    v1.NodeInternalIP,
			Address: "10.1.95.4",
		},
		{
			Type:    v1.NodeExternalIP,
			Address: "149.28.225.110",
		},
	}

	actual, err := instances.NodeAddresses(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expcted %+v got %+v", expected, actual)
	}
}

func TestInstances_InstanceType(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceType(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "vc2-4c-8gb" {
		t.Errorf("expcted %+v got %+v", "vc2-4c-8gb", actual)
	}
}

func TestInstances_InstanceID(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.InstanceID(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "75b95d83-47e2-4c0f-b273-cc9ce2b456f8" {
		t.Errorf("expcted %+v got %+v", "75b95d83-47e2-4c0f-b273-cc9ce2b456f8", actual)
	}
}

func TestInstances_CurrentNodeName(t *testing.T) {
	client := newFakeClient()
	instances := newInstances(client)

	actual, err := instances.CurrentNodeName(context.TODO(), "ccm-test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if actual != "ccm-test" {
		t.Errorf("expcted %+v got %+v", "ccm-test", actual)
	}
}

// Nodes without a provider ID (on-prem nodes in a hybrid cluster) must never be
// reported as missing, or the cloud node lifecycle controller deletes them.
// The client has no Instance service, so any Vultr API lookup would panic.
func TestInstancesV2_EmptyProviderID(t *testing.T) {
	instances := newInstancesV2(&govultr.Client{})

	for _, node := range []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "home-node-1"}},
		{ObjectMeta: metav1.ObjectMeta{Name: "home-bm-1", Labels: map[string]string{"vultr.com/baremetal": "true"}}},
	} {
		t.Run(node.Name, func(t *testing.T) {
			exists, err := instances.InstanceExists(context.TODO(), node)
			if err != nil {
				t.Errorf("InstanceExists: unexpected error: %v", err)
			}
			if !exists {
				t.Errorf("InstanceExists: expected true for node without provider ID")
			}

			shutdown, err := instances.InstanceShutdown(context.TODO(), node)
			if err != nil {
				t.Errorf("InstanceShutdown: unexpected error: %v", err)
			}
			if shutdown {
				t.Errorf("InstanceShutdown: expected false for node without provider ID")
			}

			meta, err := instances.InstanceMetadata(context.TODO(), node)
			if err == nil {
				t.Errorf("InstanceMetadata: expected an error for node without provider ID, got %+v", meta)
			}
			if meta != nil {
				t.Errorf("InstanceMetadata: expected nil metadata, got %+v", meta)
			}
		})
	}
}

func TestInstancesV2_WithProviderIDStillQueriesVultr(t *testing.T) {
	instances := newInstancesV2(newFakeClient())
	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "k8s-node-7"},
		Spec:       v1.NodeSpec{ProviderID: "vultr://75b95d83-47e2-4c0f-b273-cc9ce2b456f8"},
	}

	meta, err := instances.InstanceMetadata(context.TODO(), node)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.ProviderID != "vultr://75b95d83-47e2-4c0f-b273-cc9ce2b456f8" {
		t.Errorf("unexpected provider ID %q", meta.ProviderID)
	}
}
