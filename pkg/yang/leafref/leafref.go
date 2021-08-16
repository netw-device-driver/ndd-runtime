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
	"reflect"
	"strconv"
	"strings"

	config "github.com/netw-device-driver/ndd-grpc/config/configpb"
)

type LeafRef struct {
	LocalPath  *config.Path `json:"localPath,omitempty"`
	RemotePath *config.Path `json:"remotePath,omitempty"`
}

type ResolvedLeafRef struct {
	LocalPath  *config.Path `json:"localPath,omitempty"`
	RemotePath *config.Path `json:"remotePath,omitempty"`
	Value      string       `json:"value,omitempty"`
	Resolved   bool         `json:"resolved,omitempty"`
}

func NewResolvedLeafRefCopy(in *ResolvedLeafRef) (out *ResolvedLeafRef) {
	out = new(ResolvedLeafRef)
	if in.LocalPath != nil {
		out.LocalPath = new(config.Path)
		out.LocalPath.Elem = make([]*config.PathElem, 0)
		for _, v := range in.LocalPath.GetElem() {
			elem := &config.PathElem{}
			elem.Name = v.Name
			if len(v.GetKey()) != 0 {
				elem.Key = v.Key
			}
			out.LocalPath.Elem = append(out.LocalPath.Elem, elem)
		}
	}
	if in.RemotePath != nil {
		out.RemotePath = new(config.Path)
		out.RemotePath.Elem = make([]*config.PathElem, 0)
		for _, v := range in.RemotePath.GetElem() {
			elem := &config.PathElem{}
			elem.Name = v.Name
			if len(v.GetKey()) != 0 {
				elem.Key = v.Key
			}
			out.RemotePath.Elem = append(out.LocalPath.Elem, elem)
		}
	}
	out.Resolved = in.Resolved
	out.Value = in.Value
	return out
}

/*
func (rlref *ResolvedLeafRef) DeepCopy() *ResolvedLeafRef {
	if rlref == nil {
		return nil
	}
	out := new(ResolvedLeafRef)
	rlref.DeepCopyInto(out)
	return out
}
*/

/*
func (rlref *ResolvedLeafRef) DeepCopyInto(out *ResolvedLeafRef){
	*out = *rlref
	if rlref.LocalPath != nil {
		in, out := &rlref.LocalPath, &out.LocalPath
		out = new(*config.Path)
		*out.Elem = make([]*config.PathElem, len(in.GetElem()))
		for _, v := range rlref.LocalPath.GetElem() {
			elem := &config.PathElem{}
			elem.Name = v.Name
			if len(v.GetKey()) != 0 {
				elem.Key = v.Key
			}
			out.LocalPath.Elem = append(out.LocalPath.Elem, elem)
		}
	}
	if rlref.RemotePath != nil {
		out.RemotePath = new(config.Path)
		out.RemotePath.Elem = make([]*config.PathElem, 0)
		for _, v := range rlref.RemotePath.GetElem() {
			elem := &config.PathElem{
				Name: v.Name,
			}
			if len(v.GetKey()) != 0 {
				elem.Key = v.Key
			}
			out.RemotePath.Elem = append(out.RemotePath.Elem, elem)
		}
	}
	out.Resolved = rlref.Resolved
	out.Value = rlref.Value
}
*/

func NewLeafReaf(lPath, rPath *config.Path) *LeafRef {
	return &LeafRef{
		LocalPath:  lPath,
		RemotePath: rPath,
	}
}

// ResolveLeafRefWithJSONObject resolved the leafref in the data/object supplied via the path and returns the leafref values and leafref path augmnted with the data of the leaf reference
// we can have multiple leafrefs in an object and hence we return a list with the leafref values and the leafref path (witht the Populateed data of the object)
func (l *LeafRef) ResolveLeafRefWithJSONObject(x1 interface{}, idx int, lridx int, resolvedLeafRefs []*ResolvedLeafRef) []*ResolvedLeafRef {
	fmt.Printf("ResolveLeafRefWithJSONObject entry, idx %d, lridx: %d\n  x1: %v\n", idx, lridx, x1)
	for _, resolvedLeafRef := range resolvedLeafRefs {
		fmt.Printf("  resolvedLeafRef: %v\n", *resolvedLeafRef)
	}
	switch x := x1.(type) {
	case map[string]interface{}:
		for k, x2 := range x {
			fmt.Printf("ResolveLeafRefWithJSONObject map[string]interface{}, idx %d, lridx: %d\n l.LocalPath: %v\n k: %s, x2: %v \n", idx, lridx, l.LocalPath, k, x2)
			for _, resolvedLeafRef := range resolvedLeafRefs {
				fmt.Printf("  resolvedLeafRef: %v\n", *resolvedLeafRef)
			}
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
						/// return the value we have found the leafref
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
		// cop

		resolvedLeafRefsOrig := NewResolvedLeafRefCopy(resolvedLeafRefs[lridx]) 
		fmt.Printf("resolvedLeafRefsOrig: %v\n", resolvedLeafRefsOrig)
		for n, v := range x {
			switch x2 := v.(type) {
			case map[string]interface{}:
				for k3, x3 := range x2 {
					fmt.Printf("ResolveLeafRefWithJSONObject []interface{}, idx %d, lridx: %d\n l.LocalPath: %v\n n: %d, k3: %s, x3: %v\n ", idx, lridx, l.LocalPath, n, k3, x3)
					for _, resolvedLeafRef := range resolvedLeafRefs {
						fmt.Printf("  resolvedLeafRef: %v\n", *resolvedLeafRef)
					}
					if len(l.LocalPath.GetElem()[idx].GetKey()) != 0 {
						for k := range l.LocalPath.GetElem()[idx].GetKey() {
							if k == k3 {
								if n > 0 {
									resolvedLeafRefs = append(resolvedLeafRefs, resolvedLeafRefsOrig)
									lridx++
								}
								// check if this is the last element/index in the path
								if idx == len(l.LocalPath.GetElem())-1 {
									// return the value we have found the leafref
									fmt.Printf("ResolveLeafRefWithJSONObject []interface{} last entry in localPath\n")
									resolvedLeafRefs[lridx].PopulateLocalLeafRefValue(x3, idx)
									// we use the generic return statement to return
								} else {
									// the value is always a string since it is part of map[string]interface{}
									// since we are not at the end of the path we dont have leafRefValues and hence we dont need to Populate them
									resolvedLeafRefs[lridx].PopulateLocalLeafRefKey(x3, idx)
									// given that we can have multiple entries in the list we initialize a new index to increment independently
									i := idx
									i++
									fmt.Printf("ResolveLeafRefWithJSONObject []interface{} NOT last entry in localPath\n")
									resolvedLeafRefs = l.ResolveLeafRefWithJSONObject(x2, i, lridx, resolvedLeafRefs)
									// we use the generic return statement to return
								}
							}
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
		// the value is typically resolved using rlref.Value
		if len(rlref.LocalPath.GetElem()[idx].GetKey()) != 0 {
			for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
				rlref.LocalPath.GetElem()[idx].GetKey()[k] = x1
			}
		}
		rlref.Resolved = true
	case int:
		rlref.Value = strconv.Itoa(int(x1))
		// the value is typically resolved using rlref.Value
		if len(rlref.LocalPath.GetElem()[idx].GetKey()) != 0 {
			for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
				rlref.LocalPath.GetElem()[idx].GetKey()[k] = strconv.Itoa(int(x1))
			}
		}
		rlref.Resolved = true
	case float64:
		rlref.Value = fmt.Sprintf("%.0f", x1)
		// the value is typically resolved using rlref.Value
		if len(rlref.LocalPath.GetElem()[idx].GetKey()) != 0 {
			for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
				rlref.LocalPath.GetElem()[idx].GetKey()[k] = fmt.Sprintf("%.0f", x1)
			}
		}
		rlref.Resolved = true
	default:
		// maybe there is another type we still have to process, TBD
		fmt.Printf("PopulateLocalLeafRefValue: new type %v", x1)
	}
}

func (rlref *ResolvedLeafRef) PopulateLocalLeafRefKey(x interface{}, idx int) {
	switch x1 := x.(type) {
	case string:
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
		rlref.Value = fmt.Sprintf("%.0f", x1)
		// a leaf ref can only have 1 value, this is why this works
		for k := range rlref.LocalPath.GetElem()[idx].GetKey() {
			rlref.LocalPath.GetElem()[idx].GetKey()[k] = fmt.Sprintf("%.0f", x1)
		}
	default:
		// maybe there is another type we still have to process, TBD
		fmt.Printf("PopulateLocalLeafRefKey: new type %v", x1)
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
	fmt.Printf("FindRemoteLeafRef Entry: idx: %d, data: %v\n, remotePath: %v\n", idx, x1, rlref.RemotePath)
	switch x := x1.(type) {
	case map[string]interface{}:
		for k, x2 := range x {
			fmt.Printf("FindRemoteLeafRef map[string]interface{} List with Key: idx: %d, data: %v\n", idx, x2)
			if k == rlref.RemotePath.GetElem()[idx].GetName() {
				// check if this is the last element/index in the path
				if idx == len(rlref.RemotePath.GetElem())-1 {
					// check if the element in the path has a key
					if len(rlref.RemotePath.GetElem()[idx].GetKey()) != 0 {

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
					fmt.Printf("FindRemoteLeafRef []interface{} 1 idx: %d, k3: %v, x3: %v\n", idx, k3, x3)
					if len(rlref.RemotePath.GetElem()[idx].GetKey()) != 0 {
						for k := range rlref.RemotePath.GetElem()[idx].GetKey() {
							fmt.Printf("FindRemoteLeafRef []interface{} 2 idx: %d, k3: %v, k: %v\n", idx, k3, k)
							if k3 == k {
								// check if this is the last element/index in the path
								if idx == len(rlref.RemotePath.GetElem())-1 {
									// return the value
									fmt.Printf("FindRemoteLeafRef []interface{} 3 idx: %d, k3: %v, x3: %v, rlref.Value: %v\n", idx, k3, x3, rlref.Value)
									switch x := x3.(type) {
									case string:
										if string(x) == rlref.Value {
											return true
										}
									case uint32:
										if strconv.Itoa(int(x)) == rlref.Value {
											return true
										}
									case float64:
										fmt.Printf("a: %s, b: %s\n", fmt.Sprintf("%.0f", x), rlref.Value)
										if fmt.Sprintf("%.0f", x) == rlref.Value {
											return true
										}
									default:
										fmt.Println(reflect.TypeOf(x))
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
			}
		}
		//we will use the genric return here
	case nil:
		return false
	}
	return false
}
