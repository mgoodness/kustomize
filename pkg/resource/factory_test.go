/*
Copyright 2018 The Kubernetes Authors.

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

package resource

import (
	"reflect"
	"testing"

	"sigs.k8s.io/kustomize/pkg/internal/loadertest"
	"sigs.k8s.io/kustomize/pkg/patch"
)

func TestSliceFromPatches(t *testing.T) {

	patchGood1 := patch.StrategicMerge("/foo/patch1.yaml")
	patch1 := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pooh
`
	patchGood2 := patch.StrategicMerge("/foo/patch2.yaml")
	patch2 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: winnie
  namespace: hundred-acre-wood
---
# some comment
---
---
`
	patchBad := patch.StrategicMerge("/foo/patch3.yaml")
	patch3 := `
WOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOT: woot
`
	l := loadertest.NewFakeLoader("/foo")
	l.AddFile(string(patchGood1), []byte(patch1))
	l.AddFile(string(patchGood2), []byte(patch2))
	l.AddFile(string(patchBad), []byte(patch3))

	tests := []struct {
		name        string
		input       []patch.StrategicMerge
		expectedOut []*Resource
		expectedErr bool
	}{
		{
			name:        "happy",
			input:       []patch.StrategicMerge{patchGood1, patchGood2},
			expectedOut: []*Resource{testDeployment, testConfigMap},
			expectedErr: false,
		},
		{
			name:        "badFileName",
			input:       []patch.StrategicMerge{patchGood1, "doesNotExist"},
			expectedOut: []*Resource{},
			expectedErr: true,
		},
		{
			name:        "badData",
			input:       []patch.StrategicMerge{patchGood1, patchBad},
			expectedOut: []*Resource{},
			expectedErr: true,
		},
	}
	for _, test := range tests {
		rs, err := factory.SliceFromPatches(l, test.input)
		if test.expectedErr && err == nil {
			t.Fatalf("%v: should return error", test.name)
		}
		if !test.expectedErr && err != nil {
			t.Fatalf("%v: unexpected error: %s", test.name, err)
		}
		if len(rs) != len(test.expectedOut) {
			t.Fatalf("%s: length mismatch %d != %d",
				test.name, len(rs), len(test.expectedOut))
		}
		for i := range rs {
			if !reflect.DeepEqual(test.expectedOut[i], rs[i]) {
				t.Fatalf("%s: Got: %v\nexpected:%v",
					test.name, test.expectedOut[i], rs[i])
			}
		}
	}
}