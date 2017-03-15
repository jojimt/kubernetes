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
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/extensions/validation"
)

// networkStrategy implements verification logic for Network.
type networkStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Network objects.
var Strategy = networkStrategy{api.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all Network objects need to be within a namespace.
func (networkStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of an Network before creation.
func (networkStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
	network := obj.(*extensions.Network)
	network.Generation = 1
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (networkStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
	newNetwork := obj.(*extensions.Network)
	oldNetwork := old.(*extensions.Network)

	// Any changes to the spec increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object.
	// See metav1.ObjectMeta description for more information on Generation.
	if !reflect.DeepEqual(oldNetwork.Spec, newNetwork.Spec) {
		newNetwork.Generation = oldNetwork.Generation + 1
	}
}

// Validate validates a new Network.
func (networkStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	network := obj.(*extensions.Network)
	return validation.ValidateNetwork(network)
}

// Canonicalize normalizes the object after validation.
func (networkStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for Network; this means you may not create one with a PUT request.
func (networkStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (networkStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	validationErrorList := validation.ValidateNetwork(obj.(*extensions.Network))
	updateErrorList := validation.ValidateNetworkUpdate(obj.(*extensions.Network), old.(*extensions.Network))
	return append(validationErrorList, updateErrorList...)
}

// AllowUnconditionalUpdate is the default update policy for Network objects.
func (networkStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// NetworkToSelectableFields returns a field set that represents the object.
func NetworkToSelectableFields(network *extensions.Network) fields.Set {
	return generic.ObjectMetaFieldsSet(&network.ObjectMeta, true)
}

// GetAttrs returns labels and fields of a given object for filtering purposes.
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	network, ok := obj.(*extensions.Network)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Network.")
	}
	return labels.Set(network.ObjectMeta.Labels), NetworkToSelectableFields(network), nil
}

// MatchNetwork is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchNetwork(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}
