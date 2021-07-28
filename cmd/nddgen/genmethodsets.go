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

package nddgen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/netw-device-driver/ndd-tools/internal/comments"
	"github.com/netw-device-driver/ndd-tools/internal/generate"
	"github.com/netw-device-driver/ndd-tools/internal/match"
	"github.com/netw-device-driver/ndd-tools/internal/method"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

const (
	// LoadMode used to load all packages.
	LoadMode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax

	// DisableMarker used to disable generation of managed resource methods for
	// a type that otherwise appears to be a managed resource that is missing a
	// subset of its methods.
	DisableMarker = "ndd:generate:methods"
)
const (
	CoreAlias  = "corev1"
	CoreImport = "k8s.io/api/core/v1"

	RuntimeAlias  = "nddv1"
	RuntimeImport = "github.com/netw-device-driver/ndd-runtime/apis/common/v1"

	ResourceAlias  = "resource"
	ResourceImport = "github.com/netw-device-driver/ndd-runtime/pkg/resource"
)

const (
	errLoadPackages                     = "cannot load packages"
	errReadheaderFile                   = "cannot read header file"
	errWriteManagedResourceMethod       = "cannot write managed resource method set for package"
	errWriteManagedResourceListMethod   = "cannot write managed resource list method set for package"
	errLoadingPackages                  = "error loading packages using pattern"
	errWriteNetworkNodeMethod          = "cannot write network node methods"
	errWriteNetworkNodeUsageMethod     = "cannot write network node usage methods"
	errWriteNetworkNodeUsageListMethod = "cannot write network node usage list methods"
)

var (
	headerFile          string
	filenameManaged     string
	filenameManagedList string
	filenameNN          string
	filenameNNU         string
	filenameNNUList     string
	pattern             string
)

// startCmd represents the start command for the network device driver
var genmethodsetCmd = &cobra.Command{
	Use:          "generate-methodsets",
	Short:        "generate a ndd method sets.",
	Long:         "generate a ndd method sets.",
	Aliases:      []string{"gen-methodsets"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("ndd-gen started ...")
		pkgs, err := packages.Load(&packages.Config{Mode: LoadMode}, pattern)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s : %s", errLoadPackages, pattern))
		}

		header := ""
		if headerFile != "" {
			h, err := ioutil.ReadFile(headerFile)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", errReadheaderFile, headerFile))
			}
			header = string(h)
		}

		for _, pkg := range pkgs {
			for _, err := range pkg.Errors {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", errLoadingPackages, pattern))
			}
			if err := GenerateManaged(filenameManaged, header, pkg); err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", err, pkg.PkgPath))
			}
			if err := GenerateManagedList(filenameManagedList, header, pkg); err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", err, pkg.PkgPath))
			}
			if err := GenerateNetworkNode(filenameNN, header, pkg); err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", err, pkg.PkgPath))
			}
			if err := GenerateNetworkNodeUsage(filenameNNU, header, pkg); err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", err, pkg.PkgPath))
			}
			if err := GenerateNetworkNodeUsageList(filenameNNUList, header, pkg); err != nil {
				return errors.Wrap(err, fmt.Sprintf("%s : %s", err, pkg.PkgPath))
			}
		}
		fmt.Println("ndd-gen finished ...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genmethodsetCmd)
	genmethodsetCmd.Flags().StringVarP(&headerFile, "header-file", "", "", "The contents of this file will be added to the top of all generated files.")
	genmethodsetCmd.Flags().StringVarP(&filenameManaged, "filename-managed", "", "zz_generated.managed.go", "The filename of generated managed resource files.")
	genmethodsetCmd.Flags().StringVarP(&filenameManagedList, "filename-managed-list", "", "zz_generated.managedlist.go", "The filename of generated managed list resource files.")
	genmethodsetCmd.Flags().StringVarP(&filenameNN, "filename-nn", "", "zz_generated.nn.go", "The filename of generated NetworkNode files.")
	genmethodsetCmd.Flags().StringVarP(&filenameNNU, "filename-nnu", "", "zz_generated.nnu.go", "The filename of generated NetworkNode usage files.")
	genmethodsetCmd.Flags().StringVarP(&filenameNNUList, "filename-nnu-list", "", "zz_generated.nnulist.go", "The filename of generated NetworkNode list usage files.")
	genmethodsetCmd.Flags().StringVarP(&pattern, "paths", "", "", "Package(s) for which to generate methods, for example github.com/netw-device-driver/ndd-core/apis/...")
}

// GenerateManaged generates the resource.Managed method set.
func GenerateManaged(filename, header string, p *packages.Package) error {
	receiver := "mg"

	methods := method.Set{
		"SetActive":                  method.NewSetActive(receiver, RuntimeImport),
		"GetActive":                  method.NewGetActive(receiver, RuntimeImport),
		"SetConditions":              method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":               method.NewGetCondition(receiver, RuntimeImport),
		"GetNetworkNodeReference":   method.NewGetNetworkNodeReference(receiver, RuntimeImport),
		"SetNetworkNodeReference":   method.NewSetNetworkNodeReference(receiver, RuntimeImport),
		"SetDeletionPolicy":          method.NewSetDeletionPolicy(receiver, RuntimeImport),
		"GetDeletionPolicy":          method.NewGetDeletionPolicy(receiver, RuntimeImport),
		"InitializeTargetConditions": method.NewInitializeTargetConditions(receiver, RuntimeImport),
		"DeleteTargetCondition":      method.NewDeleteTargetCondition(receiver, RuntimeImport),
		"GetTargetCondition":         method.NewGetTargetCondition(receiver, RuntimeImport),
		"SetTargetConditions":        method.NewSetTargetConditions(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			CoreImport:    CoreAlias,
			RuntimeImport: RuntimeAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.Managed(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, errWriteManagedResourceMethod)
}

// GenerateManagedList generates the resource.ManagedList method set.
func GenerateManagedList(filename, header string, p *packages.Package) error {
	receiver := "l"

	methods := method.Set{
		"GetItems": method.NewManagedGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			ResourceImport: ResourceAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.ManagedList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, errWriteManagedResourceListMethod)
}

// GenerateNetworkNode generates the resource.NetworkNode method set.
func GenerateNetworkNode(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"SetUsers":      method.NewSetUsers(receiver),
		"GetUsers":      method.NewGetUsers(receiver),
		"SetConditions": method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":  method.NewGetCondition(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.NetworkNode(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, errWriteNetworkNodeMethod)
}

// GenerateNetworkNodeUsage generates the resource.NetworkNodeUsage method set.
func GenerateNetworkNodeUsage(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"SetNetworkNodeReference": method.NewSetRootNetworkNodeReference(receiver, RuntimeImport),
		"GetNetworkNodeReference": method.NewGetRootNetworkNodeReference(receiver, RuntimeImport),
		"SetResourceReference":     method.NewSetRootResourceReference(receiver, RuntimeImport),
		"GetResourceReference":     method.NewGetRootResourceReference(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.NetworkNodeUsage(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, errWriteNetworkNodeUsageMethod)
}

// GenerateNetworkNodeUsageList generates the
// resource.NetworkNodeUsageList method set.
func GenerateNetworkNodeUsageList(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"GetItems": method.NewNetworkNodeUsageGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.NetworkNodeUsageList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, errWriteNetworkNodeUsageListMethod)
}
