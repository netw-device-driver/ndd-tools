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
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetCondition").Params(jen.Id("ct").Qual(runtime, "ConditionKind")).Qual(runtime, "Condition").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetCondition").Call(jen.Id("ct"))),
		)
	}
}

// NewSetProviderConfigReference returns a NewMethod that writes a SetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewSetProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetProviderConfigReference").Params(jen.Id("r").Op("*").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderConfigReference").Op("=").Id("r"),
		)
	}
}

// NewGetProviderConfigReference returns a NewMethod that writes a GetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewGetProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetProviderConfigReference").Params().Op("*").Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderConfigReference")),
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
