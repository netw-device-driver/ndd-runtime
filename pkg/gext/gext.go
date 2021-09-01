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

package gext

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	errJSONMarshal   = "cannot marshal JSON object"
	errJSONUnMarshal = "cannot unmarshal JSON object"
)

type GEXTAction string

const (
	GEXTActionGet             GEXTAction = "get"
	GEXTActionDelete          GEXTAction = "delete"
	GEXTActionCreate          GEXTAction = "create"
	GEXTActionUpdate          GEXTAction = "Update"
	GEXTActionGetResourceName GEXTAction = "getresourcename"
)

func (c *GEXTAction) String() string {
	switch *c {
	case GEXTActionGet:
		return "get"
	case GEXTActionDelete:
		return "delete"
	case GEXTActionCreate:
		return "create"
	case GEXTActionGetResourceName:
		return "getresourcename"
	}
	return ""
}

type ResourceStatus string

const (
	ResourceStatusNone          ResourceStatus = "none"
	ResourceStatusSuccess       ResourceStatus = "success"
	ResourceStatusFailed        ResourceStatus = "failed"
	ResourceStatusCreatePending ResourceStatus = "createpending"
	ResourceStatusUpdatePending ResourceStatus = "updatepending"
	ResourceStatusDeletePending ResourceStatus = "deletepending"
)

func (c *ResourceStatus) String() string {
	switch *c {
	case ResourceStatusNone:
		return "none"
	case ResourceStatusSuccess:
		return "success"
	case ResourceStatusFailed:
		return "failed"
	case ResourceStatusCreatePending:
		return "createpending"
	case ResourceStatusUpdatePending:
		return "updatepending"
	case ResourceStatusDeletePending:
		return "deletepending"
	}
	return ""
}

// GEXT identifies an extension to GNMI to indicate the Actions and Parameters
// of a  managed resource
type GEXT struct {
	Action  GEXTAction     `json:"action,omitempty"`
	Name    string         `json:"name,omitempty"`
	Level   int            `json:"level,omitempty"`
	Status  ResourceStatus `json:"status,omitempty"`
	Exists  bool           `json:"exists,omitempty"`
	HasData bool           `json:"hasData,omitempty"`
}

func (gext *GEXT) String() (string, error) {
	d, err := json.Marshal(&gext)
	if err != nil {
		return "", errors.Wrap(err, errJSONMarshal)
	}
	return string(d), nil
}

func (gext *GEXT) GetAction() GEXTAction {
	return gext.Action
}

func (gext *GEXT) GetName() string {
	return gext.Name
}

func (gext *GEXT) GetLevel() int {
	return gext.Level
}

func (gext *GEXT) GetStatus() ResourceStatus {
	return gext.Status
}

func (gext *GEXT) GetExists() bool {
	return gext.Exists
}

func (gext *GEXT) GetHasData() bool {
	return gext.HasData
}

func (gext *GEXT) SetAction(s GEXTAction) {
	gext.Action = s
}

func (gext *GEXT) SetName(s string) {
	gext.Name = s
}

func (gext *GEXT) SetLevel(s int) {
	gext.Level = s
}

func (gext *GEXT) SetStatus(s ResourceStatus) {
	gext.Status = s
}

func (gext *GEXT) SetExists(s bool) {
	gext.Exists = s
}

func (gext *GEXT) SetHasData(s bool) {
	gext.HasData = s
}

func String2GEXT(s string) (*GEXT, error) {
	var gext GEXT
	fmt.Printf("String2GEXT: %s\n", s)
	if err := json.Unmarshal([]byte(s), &gext); err != nil {
		return nil, errors.Wrap(err, errJSONUnMarshal)
	}
	fmt.Printf("GEXT: %v\n", gext)
	return &gext, nil
}
