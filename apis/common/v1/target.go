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

package v1

import (
	"github.com/karimra/gnmic/target"
	"github.com/karimra/gnmic/types"
	"github.com/netw-device-driver/ndd-grpc/ndd"
)

type Target struct {
	Name   string
	Cfg    ndd.Config // this is the config client code
	Config *types.TargetConfig // this is gnmi based
	Target *target.Target // this is gnmi based
}

func (t *Target) IsTargetDeleted(ns []string) bool {
	for _, an := range ns {
		if an == t.Name {
			return false
		}
	}
	return true
}
