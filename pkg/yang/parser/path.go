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

import (
	"strings"

	config "github.com/netw-device-driver/ndd-grpc/config/configpb"
	"github.com/netw-device-driver/ndd-runtime/pkg/utils"
)

func GnmiPathToName(path *config.Path) string {
	sb := strings.Builder{}
	for i, pElem := range path.GetElem() {
		pes := strings.Split(pElem.GetName(), ":")
		var pe string
		if len(pes) > 1 {
			pe = pes[1]
		} else {
			pe = pes[0]
		}
		sb.WriteString(pe)

		if i+1 != len(path.GetElem()) {
			sb.WriteString("-")
		}
	}
	return sb.String()
}

func GnmiPathToXPath(path *config.Path, keys bool) *string {
	sb := strings.Builder{}
	for i, pElem := range path.GetElem() {
		pes := strings.Split(pElem.GetName(), ":")
		var pe string
		if len(pes) > 1 {
			pe = pes[1]
		} else {
			pe = pes[0]
		}
		sb.WriteString(pe)
		if keys {
			for k, v := range pElem.GetKey() {
				sb.WriteString("[")
				sb.WriteString(k)
				sb.WriteString("=")
				sb.WriteString(v)
				sb.WriteString("]")
			}
		}
		if i+1 != len(path.GetElem()) {
			sb.WriteString("/")
		}
	}
	return utils.StringPtr("/" + sb.String())
}

func XpathToGnmiPath(p string, offset int) (path *config.Path) {
	split := strings.Split(p, "/")
	for i, element := range split {
		// ignore the first element
		//fmt.Printf("i = %d, element = %s\n", i, element)
		if i == 0 {
			path = &config.Path{
				Elem: make([]*config.PathElem, 0),
			}
		} else {
			// offset is used to ignore an element from the path
			if i > offset {
				pathElem := &config.PathElem{}
			if element != "" {
				if strings.Contains(element, "[") {
					s1 := strings.Split(element, "[")
					pathElem.Key = make(map[string]string)
					pathElem.Name = s1[0]
					s2 := strings.Split(s1[1], ",")
					for _, eWithKey := range s2 {
						// TODO if there is a "/" in the name of the path
						s := strings.Split(eWithKey, "=")
						var v string
						if strings.Contains(s[1], "]") {
							v = strings.Trim(s[1], "]")
						} else {
							v = s[1]
						}
						// trim blanks from the final element/key/value
						pathElem.Key[strings.Trim(s[0], " ")] = strings.Trim(v, " ")
					}
				} else {
					// trim blanks from the final element/key/value
					pathElem.Name = strings.Trim(element, " ")
				}
				path.Elem = append(path.Elem, pathElem)
			}
			}
			
		}
	}
	return path
}

func InsertElemInPath(path *config.Path, name, key string) *config.Path {
	pathElem := &config.PathElem{
		Name: name,
	}
	if key != "" {
		pathElem.Key = make(map[string]string)
		pathElem.Key[key] = ""
	}

	path.Elem = append(path.Elem, pathElem)
	return path
}

// for a leafRef path the last entry of the name should be a key in the previous element
func TransformPathToLeafRefPath(path *config.Path) *config.Path {
	key := path.GetElem()[len(path.GetElem())-1].Name
	path.Elem = path.Elem[:(len(path.GetElem())-1)]
	path.GetElem()[len(path.GetElem())-1].Key = make(map[string]string)
	path.GetElem()[len(path.GetElem())-1].Key[key] = ""
	return path
}

// RemoveFirstEntry removes the first entry of the xpath, so it trims the first element of the /
func RemoveFirstEntry(s string) string {
	split := strings.Split(s, "/")
	var p string
	for i, s := range split {
		if i > 1 {
			p += "/" + s
		}
	}
	return p
}

func TransformPathAsRelative2Resource(localPath, activeResPath *config.Path) *config.Path {
	localPath.Elem = localPath.Elem[(len(activeResPath.GetElem())-1):(len(localPath.GetElem()))]
	return localPath
}