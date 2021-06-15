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

// PackageLinterFn lints an entire package. If function applies a check for
// multiple objects, consider using an ObjectLinterFn.
type PackageLinterFn func(*Package) error

// ObjectLinterFn lints an object in a package.
type ObjectLinterFn func(runtime.Object) error

// PackageLinter lints packages by applying package and object linter functions
// to it.
type PackageLinter struct {
	pre       []PackageLinterFn
	perMeta   []ObjectLinterFn
	perObject []ObjectLinterFn
}

// NewPackageLinter creates a new PackageLinter.
func NewPackageLinter(pre []PackageLinterFn, perMeta, perObject []ObjectLinterFn) *PackageLinter {
	return &PackageLinter{
		pre:       pre,
		perMeta:   perMeta,
		perObject: perObject,
	}
}
