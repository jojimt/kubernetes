/*
Copyright 2016 The Kubernetes Authors.

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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func TestNetworkStrategy(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	if !Strategy.NamespaceScoped() {
		t.Errorf("Network must be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("Network should not allow create on update")
	}

	nw := &extensions.Network{
		ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
		Spec: extensions.NetworkSpec{
			Plugin: "xyz",
		},
	}

	Strategy.PrepareForCreate(ctx, nw)
	errs := Strategy.Validate(ctx, nw)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}

	invalidNw := &extensions.Network{
		ObjectMeta: metav1.ObjectMeta{Name: "bar", ResourceVersion: "4"},
	}
	Strategy.PrepareForUpdate(ctx, invalidNw, nw)
	errs = Strategy.ValidateUpdate(ctx, invalidNw, nw)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
	if invalidNw.ResourceVersion != "4" {
		t.Errorf("Incoming resource version on update should not be mutated")
	}
}
