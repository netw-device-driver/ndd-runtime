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

package leafref

import (
	"fmt"
	"strconv"
	"strings"

	config "github.com/netw-device-driver/ndd-grpc/config/configpb"
)

type LeafRef struct {
	LocalPath  *config.Path
	RemotePath *config.Path
}

type ResolvedLeafRef struct {
	LocalPath  *config.Path
	RemotePath *config.Path
	Value      string
	Resolved   bool
}

func NewLeafReaf(lPath, rPath *config.Path) *LeafRef {
	return &LeafRef{
		LocalPath:  lPath,
		RemotePath: rPath,
	}
}

// ResolveLeafRefWithJSONObject resolved the leafref in the data/object supplied via the path and returns the leafref values and leafref path augmnted with the data of the leaf reference
// we can have multiple leafrefs in an object and hence we return a list with the leafref values and the leafref path (witht the Populateed data of the object)
func (l *LeafRef) ResolveLeafRefWithJSONObject(x1 interface{}, idx int, lridx int, resolvedLeafRefs []*ResolvedLeafRef) []*ResolvedLeafRef {
	switch x := x1.(type) {
	case map[string]interface{}:
		for k, x2 := range x {
			if k == l.LocalPath.GetElem()[idx].GetName() {
				// check if this is the last element/index in the path
				if idx == len(l.LocalPath.GetElem())-1 {
					// check if the element in the path has a key
					if len(l.LocalPath.GetElem()[idx].GetKey()) != 0 {
						// given the element has a key we need to go through the []interface{} part
						resolvedLeafRefs = l.ResolveLeafRefWithJSONObject(x2, idx, lridx, resolvedLeafRefs)
						// we use the generic return statement to return
					} else {
						// element in the path has no key
						// return the value
						resolvedLeafRefs[lridx].PopulateLocalLeafRefValue(x2, idx)
						// we use the generic return statement to return
					}
				} else {
					// not last element/index in path
					// check if the element in the path has a key
					if len(l.LocalPath.GetElem()[idx].GetKey()) != 0 {
						// given the element has a key we need to go through the []interface{} part
						resolvedLeafRefs = l.ResolveLeafRefWithJSONObject(x2, idx, lridx, resolvedLeafRefs)
						// we use the generic return statement to return
					} else {
						// increment the index
						idx++
						resolvedLeafRefs = l.ResolveLeafRefWithJSONObject(x2, idx, lridx, resolvedLeafRefs)
						// we use the generic return statement to return
					}
				}
			}
		}
	case []interface{}:
		resolvedLeafRefsOrig := resolvedLeafRefs[lridx]
		for n, v := range x {
			switch x2 := v.(type) {
			case map[string]interface{}:
				for k3, x3 := range x2 {
					if k3 == l.LocalPath.GetElem()[idx].GetName() {
						if n > 0 {

							resolvedLeafRefs = append(resolvedLeafRefs, resolvedLeafRefsOrig)
							lridx++
						}
						// check if this is the last element/index in the path
						if idx == len(l.LocalPath.GetElem())-1 {
							// return the value
							resolvedLeafRefs[lridx].PopulateLocalLeafRefValue(x3, idx)
							// we use the generic return statement to return
						} else {
							// the value is always a string since it is part of map[string]interface{}
							// since we are not at the end of the path we dont have leafRefValues and hence we dont need to Populate them
							resolvedLeafRefs[lridx].PopulateLocalLeafRefKey(x3, idx)
							// given that we can have multiple entries in the list we initialize a new index to increment independently
							i := idx
							i++
							resolvedLeafRefs = l.ResolveLeafRefWithJSONObject(x2, i, lridx, resolvedLeafRefs)
							// we use the generic return statement to return
						}
					}
				}
			}
		}
		//we will use the geenric return here
	case nil:
		// When we come here we have no resolution
	}
	return resolvedLeafRefs
}

// PopulateLeafRefWithValue fills out the keys in the path and Populateed the Values if the leafRefValues is not nil
func (rlref *ResolvedLeafRef) PopulateLocalLeafRefValue(x interface{}, idx int) {
	switch x1 := x.(type) {
	case string:
		rlref.Value = x1
		// a leaf ref can only have 1 value, this is why this works
		for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
			rlref.LocalPath.GetElem()[idx].GetKey()[k] = x1
		}
	case int:
		rlref.Value = strconv.Itoa(int(x1))
		// a leaf ref can only have 1 value, this is why this works
		for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
			rlref.LocalPath.GetElem()[idx].GetKey()[k] = strconv.Itoa(int(x1))
		}

	case float64:
		rlref.Value = fmt.Sprintf("%f", x1)
		// a leaf ref can only have 1 value, this is why this works
		for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
			rlref.LocalPath.GetElem()[idx].GetKey()[k] = fmt.Sprintf("%f", x1)
		}
	default:
		// maybe there is another type we still have to process, TBD
	}
}

func (rlref *ResolvedLeafRef) PopulateLocalLeafRefKey(x interface{}, idx int) {
	switch x1 := x.(type) {
	case string:
		// a leaf ref can only have 1 value, this is why this works
		for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
			rlref.LocalPath.GetElem()[idx].GetKey()[k] = x1
		}
	default:
		// a key only has a string tyope, so we should never come here
	}
}

func (rlref *ResolvedLeafRef) PopulateRemoteLeafRefKey() {
	split := strings.Split(rlref.Value, ".")
	n := 0
	for _, PathElem := range rlref.RemotePath.GetElem() {
		if len(PathElem.GetKey()) != 0 {
			for k := range PathElem.GetKey() {
				PathElem.GetKey()[k] = split[n]
				n++
			}
		}
	}
}

func (rlref *ResolvedLeafRef) FindRemoteLeafRef(x1 interface{}, idx int) (found bool) {
	fmt.Printf("FindRemoteLeafRef Entry: idx: %d, data: %v", idx, x1)
	switch x := x1.(type) {
	case map[string]interface{}:
		for k, x2 := range x {
			if k == rlref.RemotePath.GetElem()[idx].GetName() {
				// check if this is the last element/index in the path
				if idx == len(rlref.RemotePath.GetElem())-1 {
					// check if the element in the path has a key
					if len(rlref.RemotePath.GetElem()[idx].GetKey()) != 0 {
						fmt.Printf("FindRemoteLeafRef map[string]interface{} List with Key: idx: %d, data: %v", idx, x2)
						// given the element has a key we need to go through the []interface{} part
						return rlref.FindRemoteLeafRef(x2, idx)
					} else {
						// element in the path has no key
						if x2 == rlref.Value {
							return true
						} else {
							return false
						}
					}
				} else {
					// not last element/index in path
					// check if the element in the path has a key
					if len(rlref.RemotePath.GetElem()[idx].GetKey()) != 0 {
						// given the element has a key we need to go through the []interface{} part
						return rlref.FindRemoteLeafRef(x2, idx)
						// we use the generic return statement to return
					} else {
						// increment the index
						idx++
						return rlref.FindRemoteLeafRef(x2, idx)
					}
				}
			}
		}
	case []interface{}:
		for _, v := range x {
			switch x2 := v.(type) {
			case map[string]interface{}:
				for k3, x3 := range x2 {
					fmt.Printf("FindRemoteLeafRef []interface{} idx: %d, k3: %v, x3: %v", idx, k3, x3)
					if k3 == rlref.RemotePath.GetElem()[idx].GetName() {
						// check if this is the last element/index in the path
						if idx == len(rlref.RemotePath.GetElem())-1 {
							// return the value
							if x3 == rlref.Value {
								return true
							}
							// we should not return here since there can be multiple elements in the
							// list and we need to exercise them all, the geenric return will take care of it
						} else {
							idx++
							return rlref.FindRemoteLeafRef(x2, idx)
						}
					}
				}
			}
		}
		//we will use the genric return here
	case nil:
		return false
	}
	return false
}
