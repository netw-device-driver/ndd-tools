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

package method

import (
	"go/token"
	"go/types"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/netw-device-driver/ndd-tools/internal/fields"
)

// New is a function that adds a method on the supplied object in the
// supplied file.
type New func(f *jen.File, o types.Object)

// A Set is a map of method names to the New functions that produce
// them.
type Set map[string]New

// Write the method Set for the supplied Object to the supplied file. Methods
// are filtered by the supplied Filter.
func (s Set) Write(f *jen.File, o types.Object, mf Filter) {
	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if mf(o, name) {
			continue
		}
		s[name](f, o)
	}
}

// A Filter is a function that determines whether a method should be written for
// the supplied object. It returns true if the method should be filtered.
type Filter func(o types.Object, methodName string) bool

// DefinedOutside returns a MethodFilter that returns true if the supplied
// object has a method with the supplied name that is not defined in the
// supplied filename. The object's filename is determined using the supplied
// FileSet.
func DefinedOutside(fs *token.FileSet, filename string) Filter {
	return func(o types.Object, name string) bool {
		s := types.NewMethodSet(types.NewPointer(o.Type()))
		for i := 0; i < s.Len(); i++ {
			mo := s.At(i).Obj()
			if mo.Name() != name {
				continue
			}
			if fs.Position(mo.Pos()).Filename != filename {
				return true
			}
		}
		return false
	}
}

// NewSetActive returns a NewMethod that writes a SetActive method for
// the supplied Object to the supplied file.
func NewSetActive(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetActive of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetActive").Params(jen.Id("b").Bool()).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("Active").Op("=").Id("b"),
		)
	}
}

// NewGetActive returns a NewMethod that writes a GetActive method for
// the supplied Object to the supplied file.
func NewGetActive(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetActive of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetActive").Params().Bool().Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("Active")),
		)
	}
}

// NewSetConditions returns a NewMethod that writes a SetConditions method for
// the supplied Object to the supplied file.
func NewSetConditions(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetConditions of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetConditions").Params(jen.Id("c").Op("...").Qual(runtime, "Condition")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetConditions").Call(jen.Id("c").Op("...")),
		)
	}
}

// NewGetCondition returns a NewMethod that writes a GetCondition method for
// the supplied Object to the supplied file.
func NewGetCondition(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetCondition of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetCondition").Params(jen.Id("ck").Qual(runtime, "ConditionKind")).Qual(runtime, "Condition").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetCondition").Call(jen.Id("ck"))),
		)
	}
}

// NewSetNetworkNodeReference returns a NewMethod that writes a SetNetworkNodeReference
// method for the supplied Object to the supplied file.
func NewSetNetworkNodeReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetNetworkNodeReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetNetworkNodeReference").Params(jen.Id("r").Op("*").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("NetworkNodeReference").Op("=").Id("r"),
		)
	}
}

// NewGetNetworkNodeReference returns a NewMethod that writes a GetNetworkNodeReference
// method for the supplied Object to the supplied file.
func NewGetNetworkNodeReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetNetworkNodeReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetNetworkNodeReference").Params().Op("*").Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("NetworkNodeReference")),
		)
	}
}

// NewSetDeletionPolicy returns a NewMethod that writes a SetDeletionPolicy
// method for the supplied Object to the supplied file.
func NewSetDeletionPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetDeletionPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetDeletionPolicy").Params(jen.Id("r").Qual(runtime, "DeletionPolicy")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("DeletionPolicy").Op("=").Id("r"),
		)
	}
}

// NewGetDeletionPolicy returns a NewMethod that writes a GetDeletionPolicy
// method for the supplied Object to the supplied file.
func NewGetDeletionPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetDeletionPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetDeletionPolicy").Params().Qual(runtime, "DeletionPolicy").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("DeletionPolicy")),
		)
	}
}

// NewInitializeTargetConditions returns a NewMethod that writes a InitializeTargetConditions
// method for the supplied Object to the supplied file.
func NewInitializeTargetConditions(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("InitializeTargetConditions of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("InitializeTargetConditions").Params().Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Op("=").Id("make").Params(jen.Id("map").Index(jen.String()).Op("*").Qual(runtime, "TargetConditions")),
		)
	}
}

// NewDeleteTargetCondition returns a NewMethod that writes a DeleteTargetCondition
// method for the supplied Object to the supplied file.
func NewDeleteTargetCondition(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("DeleteTargetCondition of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("DeleteTargetCondition").Params(jen.Id("target").String()).Block(
			jen.Id("delete").Params(jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions"), jen.Id("target")),
		)
	}
}

// NewGetTargetCondition returns a NewMethod that writes a GetTargetCondition
// method for the supplied Object to the supplied file.
func NewGetTargetCondition(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetTargetCondition of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetTargetCondition").Params(jen.Id("target").String(), jen.Id("ck").Qual(runtime, "ConditionKind")).Qual(runtime, "Condition").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Index(jen.Id("target")).Dot("GetCondition").Params(jen.Id("ck"))),
		)
	}
}

// NewSetTargetConditions returns a NewMethod that writes a SetTargetConditions
// method for the supplied Object to the supplied file.
func NewSetTargetConditions(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetTargetConditions of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetTargetConditions").Params(jen.Id("target").String(), jen.Id("c").Op("...").Qual(runtime, "Condition")).Block(
			jen.If(jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Index(jen.Id("target")).Op("==").Nil()).Block(
				jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Index(jen.Id("target")).Op("=").New(jen.Qual(runtime, "TargetConditions")),
			//	jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Index(jen.Id("target").Op("=").New(jen.Qual(runtime, "TargetConditions")),
			),
			jen.Id(receiver).Dot(fields.NameStatus).Dot("TargetConditions").Index(jen.Id("target")).Dot("SetConditions").Params(jen.Id("c").Op("...")),
		)
	}
}

// NewSetUsers returns a NewMethod that writes a SetUsers method for the
// supplied Object to the supplied file.
func NewSetUsers(receiver string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetUsers of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetUsers").Params(jen.Id("i").Int64()).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("Users").Op("=").Id("i"),
		)
	}
}

// NewGetUsers returns a NewMethod that writes a GetUsers method for the
// supplied Object to the supplied file.
func NewGetUsers(receiver string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetUsers of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetUsers").Params().Int64().Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("Users")),
		)
	}
}

// NewManagedGetItems returns a New that writes a GetItems method for the
// supplied object to the supplied file.
func NewManagedGetItems(receiver, resource string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetItems of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetItems").Params().Index().Qual(resource, "Managed").Block(
			jen.Id("items").Op(":=").Make(jen.Index().Qual(resource, "Managed"), jen.Len(jen.Id(receiver).Dot("Items"))),
			jen.For(jen.Id("i").Op(":=").Range().Id(receiver).Dot("Items")).Block(
				jen.Id("items").Index(jen.Id("i")).Op("=").Op("&").Id(receiver).Dot("Items").Index(jen.Id("i")),
			),
			jen.Return(jen.Id("items")),
		)
	}
}

// NewSetRootNetworkNodeReference returns a NewMethod that writes a
// SetNetworkNodeReference method for the supplied Object to the supplied
// file. Note that unlike NewSetNetworkNodeReference the generated method
// expects the NetworkNodeReference to be at the root of the struct, not
// under its Spec field.
func NewSetRootNetworkNodeReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetNetworkNodeReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetNetworkNodeReference").Params(jen.Id("r").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot("NetworkNodeReference").Op("=").Id("r"),
		)
	}
}

// NewGetRootNetworkNodeReference returns a NewMethod that writes a
// GetNetworkNodeReference method for the supplied Object to the supplied
// file. file. Note that unlike NewGetNetworkNodeReference the generated
// method expects the NetworkNodeReference to be at the root of the struct,
// not under its Spec field.
func NewGetRootNetworkNodeReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetNetworkNodeReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetNetworkNodeReference").Params().Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot("NetworkNodeReference")),
		)
	}
}

// NewSetRootResourceReference returns a NewMethod that writes a
// SetRootResourceReference method for the supplied Object to the supplied file.
func NewSetRootResourceReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetResourceReference").Params(jen.Id("r").Qual(runtime, "TypedReference")).Block(
			jen.Id(receiver).Dot("ResourceReference").Op("=").Id("r"),
		)
	}
}

// NewGetRootResourceReference returns a NewMethod that writes a
// GetRootResourceReference method for the supplied Object to the supplied file.
func NewGetRootResourceReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetResourceReference").Params().Qual(runtime, "TypedReference").Block(
			jen.Return(jen.Id(receiver).Dot("ResourceReference")),
		)
	}
}

// NewNetworkNodeUsageGetItems returns a New that writes a GetItems method for the
// supplied object to the supplied file.
func NewNetworkNodeUsageGetItems(receiver, resource string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetItems of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetItems").Params().Index().Qual(resource, "NetworkNodeUsage").Block(
			jen.Id("items").Op(":=").Make(jen.Index().Qual(resource, "NetworkNodeUsage"), jen.Len(jen.Id(receiver).Dot("Items"))),
			jen.For(jen.Id("i").Op(":=").Range().Id(receiver).Dot("Items")).Block(
				jen.Id("items").Index(jen.Id("i")).Op("=").Op("&").Id(receiver).Dot("Items").Index(jen.Id("i")),
			),
			jen.Return(jen.Id("items")),
		)
	}
}