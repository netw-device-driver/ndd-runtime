/*
Copyright 2020 Wim Henderickx.

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

package parser

import "k8s.io/apimachinery/pkg/runtime"

// ObjectCreaterTyper know how to create and determine the type of objects.
type ObjectCreaterTyper interface {
	runtime.ObjectCreater
	runtime.ObjectTyper
}

// Package is the set of metadata and objects in a package.
type Package struct {
	meta    []runtime.Object
	objects []runtime.Object
}

// PackageParser is a Parser implementation for parsing packages.
type PackageParser struct {
	metaScheme ObjectCreaterTyper
	objScheme  ObjectCreaterTyper
}

// New returns a new PackageParser.
func New(meta, obj ObjectCreaterTyper) *PackageParser {
	return &PackageParser{
		metaScheme: meta,
		objScheme:  obj,
	}
}
