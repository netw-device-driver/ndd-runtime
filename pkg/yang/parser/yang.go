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
	"encoding/json"
	"fmt"
	"strings"

	config "github.com/netw-device-driver/ndd-grpc/config/configpb"
	"github.com/netw-device-driver/ndd-runtime/pkg/yang/container"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/stoewer/go-strcase"
)

func GetTypeName(e *yang.Entry) string {
	if e == nil || e.Type == nil {
		return ""
	}
	// Return our root's type name.
	// This is should be the builtin type-name
	// for this entry.
	return e.Type.Name
}

func GetTypeKind(e *yang.Entry) string {
	if e == nil || e.Type == nil {
		return ""
	}
	// Return our root's type name.
	// This is should be the builtin type-name
	// for this entry.
	return e.Type.Kind.String()
}

func InitializePathElem(e *yang.Entry) *config.PathElem {
	pathElem := &config.PathElem{
		Name:      e.Name,
		Key:       make(map[string]string),
		Attribute: new(config.Attribute),
	}

	if e.Key != "" {
		var keyType string
		switch GetTypeName(e.Dir[e.Key]) {
		case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
			keyType = GetTypeName(e.Dir[e.Key])
		case "boolean":
			keyType = "bool"
		case "enumeration":
			keyType = "string"
		default:
			keyType = "string"
		}
		pathElem.Key[e.Key] = keyType
	}
	return pathElem
}

func CreateContainerEntry(e *yang.Entry, next, prev *container.Container) *container.Entry {
	// Allocate a new Entry
	entry := container.NewEntry(e.Name)

	// initialize the Next pointer if relevant -> only relevant for list
	entry.Next = next
	entry.Prev = prev

	// process mandatory attribute
	switch e.Mandatory {
	case 1: // TSTrue
		entry.Mandatory = true
	default: // TSTrue
		entry.Mandatory = false
	}
	if e.Key != "" {
		entry.Mandatory = true
	}

	// process type attaribute
	switch GetTypeName(e) {
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
		entry.Type = GetTypeName(e)
	case "boolean":
		entry.Type = "bool"
	case "enumeration":
		entry.Enum = e.Type.Enum.Names()
		entry.Type = "string"
	default:
		switch GetTypeKind(e) {
		case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
			entry.Type = GetTypeKind(e)
		case "boolean":
			entry.Type = "bool"
		case "union":
			entry.Type = "string"
			entry.Union = true
			for _, t := range e.Type.Type {
				entry.Type = t.Root.Kind.String()
				if entry.Type == "enumeration" ||
					entry.Type == "leafref" ||
					entry.Type == "union" {
					entry.Type = "string"
				}
				entry.Pattern = append(entry.Pattern, t.Pattern...)

			}
		case "leafref":
			// The processing of leaf refs is handled in another function
			entry.Type = "string"
		default:
			entry.Type = "string"
		}
	}
	// process elementType for a Key
	if e.Key != "" {
		switch GetTypeName(e.Dir[e.Key]) {
		case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
			entry.Type = GetTypeName(e.Dir[e.Key])
		case "boolean":
			entry.Type = "bool"
		default:
			entry.Type = "string"
		}
	}
	// update the Type to reflect the reference to the proper struct
	if entry.Prev != nil {
		entry.Type = strcase.UpperCamelCase(entry.Prev.GetFullName() + "-" + e.Name)
	}

	if e.Type != nil {
		for _, ra := range e.Type.Range {
			entry.Range = append(entry.Range, int(ra.Min.Value))
			entry.Range = append(entry.Range, int(ra.Max.Value))
			//fmt.Printf("RANGE MIN: %d MAX: %d, TOTAL: %d\n", ra.Min.Value, ra.Max.Value, entry.Range)
		}

		for _, le := range e.Type.Length {
			entry.Length = append(entry.Length, int(le.Min.Value))
			entry.Length = append(entry.Length, int(le.Max.Value))
			//fmt.Printf("LENGTH MIN: %d MAX: %d, TOTAL: %d\n", le.Min.Value, le.Max.Value, entry.Length)
		}

		if e.Type.Pattern != nil {
			entry.Pattern = append(entry.Pattern, e.Type.Pattern...)
			//fmt.Printf("LEAF NAME: %s, PATTERN: %s\n", e.Name, entry.Pattern)

		}
		if e.Type.Kind.String() == "enumeration" {
			entry.Enum = e.Type.Enum.Names()
		}
		if e.Default != "" {
			entry.Default = e.Default
		}
	}

	// pattern post processing
	var pattern string
	for i, p := range entry.Pattern {
		if i == 0 {
			pattern += p
		} else {
			pattern += "|" + p
		}
	}
	if len(pattern) > 0 {
		//pattern = strings.ReplaceAll(pattern, "@", "")
		//pattern = strings.ReplaceAll(pattern, "#", "")
		//pattern = strings.ReplaceAll(pattern, "$", "")
		entry.PatternString = strings.ReplaceAll(pattern, "%", "")
	}

	// enum post processing
	for _, enum := range entry.Enum {
		entry.EnumString += "`" + enum + "`;"
	}
	if entry.EnumString != "" {
		entry.EnumString = strings.TrimRight(entry.EnumString, ";")
	}

	// key handling
	entry.Key = e.Key
	return entry
}

// ProcessLeafRef processes the leafref and returns if a leafref localPath, remotePath and if the leafRef is local or external to the resource
func ProcessLeafRef(e *yang.Entry, resfullPath string, activeResPath *config.Path) (*config.Path, *config.Path, bool) {
	switch GetTypeName(e) {
	default:
		switch GetTypeKind(e) {
		case "leafref":
			//fmt.Println(e.Node.Statement().String())
			splitData := strings.Split(e.Node.Statement().String(), "\n")
			var path string
			var elem string
			var k string
			for _, s := range splitData {
				if strings.Contains(s, "path ") {
					// strip the junk from the leafref to get a plain xpath
					//fmt.Printf("LeafRef Path: %s\n", s)
					s = strings.ReplaceAll(s, "path ", "")
					s = strings.ReplaceAll(s, ";", "")
					s = strings.ReplaceAll(s, "\"", "")
					s = strings.ReplaceAll(s, " ", "")
					s = strings.ReplaceAll(s, "\t", "")
					//fmt.Printf("LeafRef Path: %s\n", s)

					// split the leafref per "/" and split the element and key from the path
					// last element is the key
					// 2nd last element is the element
					split2data := strings.Split(s, "/")
					//fmt.Printf("leafRef Len Split2 %d\n", len(split2data))

					for i, s2 := range split2data {
						switch i {
						case 0: // the first element in the leafref split is typically "", since the string before the "/" is empty
							if s2 != "" { // if not empty ensure we use the right data and split the string before ":" sign
								path += "/" + strings.Split(s2, ":")[len(strings.Split(s2, ":"))-1]

							}
						case (len(split2data) - 1): // last element is the key
							k = strings.Split(s2, ":")[len(strings.Split(s2, ":"))-1]
						case (len(split2data) - 2): // 2nd last element is the element
							elem = strings.Split(s2, ":")[len(strings.Split(s2, ":"))-1]
						default: // any other element gets added to the list
							path += "/" + strings.Split(s2, ":")[len(strings.Split(s2, ":"))-1]

						}
					}
					// if no path element exits we take the root "/" path
					if path == "" {
						path = "/"
					}
					// if the path contains /.. this is a relative leafref path
					relativeIndex := strings.Count(path, "/..")
					if relativeIndex > 0 {
						//fmt.Printf("leafRef Relative Path: %s, Element: %s, Key: %s, '/..' count %d\n", path, elem, k, relativeIndex)
						// check if the final p contains relative indirection to the resourcePath -> "/.."
						resSplitData := strings.Split(RemoveFirstEntry(resfullPath), "/")
						//fmt.Printf("ResPath Split Length: %d data: %v\n", len(resSplitData), resSplitData)
						var addString string
						for i := 1; i <= (len(resSplitData) - 1 - strings.Count(path, "/..")); i++ {
							addString += "/" + resSplitData[i]
						}
						//fmt.Printf("leafRef Absolute Path Add string: %s\n", addString)
						path = addString + strings.ReplaceAll(path, "/..", "")
					}
					//fmt.Printf("leafRef Absolute Path: %s, Element: %v, Key: %s, '/..' count %d\n", path, e, k, relativeIndex)

				}
			}
			//fmt.Printf("Path: %s, Elem: %s, Key: %s\n", path, elem, k)
			remotePath := XpathToGnmiPath(path, 0)
			InsertElemInPath(remotePath, elem, k)

			// build a gnmi path and remove the first entry since the yang contains a duplicate path
			localPath := XpathToGnmiPath(resfullPath, 1)
			// the last element hould be a key in the previous element
			localPath = TransformPathToLeafRefPath(localPath)

			if strings.Contains(*GnmiPathToXPath(remotePath, false), *GnmiPathToXPath(activeResPath, false)) {
				// this is a local leafref within the resource
				// make the localPath and remotePath relative to the resource
				//fmt.Printf("localPath: %v, remotePath %v, activePath %v\n", localPath, remotePath, activeResPath)
				localPath = TransformPathAsRelative2Resource(localPath, activeResPath)
				remotePath = TransformPathAsRelative2Resource(remotePath, activeResPath)
				//fmt.Printf("localPath: %v, remotePath %v\n", localPath, remotePath)
				return localPath, remotePath, true
			}
			// leafref is external to the resource
			//fmt.Printf("localPath: %v, remotePath %v, activePath %v\n", localPath, remotePath, activeResPath)
			// make the localPath relative to the resource
			localPath = TransformPathAsRelative2Resource(localPath, activeResPath)
			//fmt.Printf("localPath: %v, remotePath %v\n", localPath, remotePath)

			return localPath, remotePath, false
		}
	}
	return nil, nil, false
}

// RemoveHierarchicalKeys removes the hierarchical keys from the data
func RemoveHierarchicalKeys(d []byte, hids []string) ([]byte, error) {
	var x map[string]interface{}
	json.Unmarshal(d, &x)

	fmt.Printf("data before hierarchical key removal: %v\n", x)

	// we first go over hierarchical ids since when they are empty it optimizes the processing
	for _, h := range hids {
		for k := range x {
			if k == h {
				delete(x, k)
			}
		}
	}

	fmt.Printf("data after hierarchical key removal: %v\n", x)
	return json.Marshal(x)
}
