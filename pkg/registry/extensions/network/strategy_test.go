/*
Copyright 2017 The Kubernetes Authors.

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

package network

import (
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	metav1 "k8s.io/kubernetes/pkg/apis/meta/v1"
)

func TestNetworkStrategy(t *testing.T) {
	ctx := api.NewDefaultContext()
	if !Strategy.NamespaceScoped() {
		t.Errorf("Network must be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("NetworkPolicy should not allow create on update")
	}

	validMatchLabels := map[string]string{"a": "b"}
	nw := &extensions.Network{
		ObjectMeta: api.ObjectMeta{Name: "net1", Namespace: api.NamespaceDefault},
		Spec: extensions.NetworkSpec{
			Name: "net1",
			Plugin: "vendorA",
			Args: []{"arg1=aaa", "arg2=bbb"},
			CIDRMask: "11.22.192.00/18",
			Gateway: "11.22.192.1",
		},
	}

	Strategy.PrepareForCreate(ctx, nw)
	errs := Strategy.Validate(ctx, nw)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %+v - %v", nw, errs)
	}

	invalidNp := &extensions.NetworkPolicy{
		ObjectMeta: api.ObjectMeta{Name: "bar", ResourceVersion: "4"},
	}
	Strategy.PrepareForUpdate(ctx, invalidNp, np)
	errs = Strategy.ValidateUpdate(ctx, invalidNp, np)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
	if invalidNp.ResourceVersion != "4" {
		t.Errorf("Incoming resource version on update should not be mutated")
	}
}
