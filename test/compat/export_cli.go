//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"os"
	"sort"

	plycmd "github.com/devdimensionlab/plybuild/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type treeContract struct {
	SchemaVersion int               `json:"schema_version"`
	Root          string            `json:"root"`
	Commands      []commandContract `json:"commands"`
}

type commandContract struct {
	Path                  string            `json:"path"`
	Use                   string            `json:"use"`
	Aliases               []string          `json:"aliases"`
	SuggestFor            []string          `json:"suggest_for"`
	Short                 string            `json:"short"`
	Long                  string            `json:"long"`
	Example               string            `json:"example"`
	Hidden                bool              `json:"hidden"`
	Deprecated            string            `json:"deprecated"`
	DisableFlagParsing    bool              `json:"disable_flag_parsing"`
	DisableFlagsInUseLine bool              `json:"disable_flags_in_use_line"`
	TraverseChildren      bool              `json:"traverse_children"`
	Runnable              bool              `json:"runnable"`
	HasArgsValidator      bool              `json:"has_args_validator"`
	HasPreRun             bool              `json:"has_pre_run"`
	HasPersistentPreRun   bool              `json:"has_persistent_pre_run"`
	HasPostRun            bool              `json:"has_post_run"`
	HasPersistentPostRun  bool              `json:"has_persistent_post_run"`
	Annotations           map[string]string `json:"annotations"`
	Flags                 []flagContract    `json:"flags"`
}

type flagContract struct {
	Scope               string              `json:"scope"`
	Name                string              `json:"name"`
	Shorthand           string              `json:"shorthand"`
	Type                string              `json:"type"`
	Default             string              `json:"default"`
	NoOptDefault        string              `json:"no_opt_default"`
	Usage               string              `json:"usage"`
	Hidden              bool                `json:"hidden"`
	Deprecated          string              `json:"deprecated"`
	ShorthandDeprecated string              `json:"shorthand_deprecated"`
	Annotations         map[string][]string `json:"annotations"`
}

func main() {
	root := plycmd.RootCmd
	root.InitDefaultHelpCmd()
	initializeHelpFlags(root)
	tree := treeContract{SchemaVersion: 1, Root: root.Name()}
	appendCommands(root, &tree.Commands)
	sort.Slice(tree.Commands, func(i, j int) bool {
		return tree.Commands[i].Path < tree.Commands[j].Path
	})
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tree); err != nil {
		panic(err)
	}
}

func initializeHelpFlags(command *cobra.Command) {
	command.InitDefaultHelpFlag()
	for _, child := range command.Commands() {
		initializeHelpFlags(child)
	}
}

func appendCommands(command *cobra.Command, output *[]commandContract) {
	aliases := append([]string(nil), command.Aliases...)
	suggestFor := append([]string(nil), command.SuggestFor...)
	sort.Strings(aliases)
	sort.Strings(suggestFor)
	contract := commandContract{
		Path:                  command.CommandPath(),
		Use:                   command.Use,
		Aliases:               aliases,
		SuggestFor:            suggestFor,
		Short:                 command.Short,
		Long:                  command.Long,
		Example:               command.Example,
		Hidden:                command.Hidden,
		Deprecated:            command.Deprecated,
		DisableFlagParsing:    command.DisableFlagParsing,
		DisableFlagsInUseLine: command.DisableFlagsInUseLine,
		TraverseChildren:      command.TraverseChildren,
		Runnable:              command.Runnable(),
		HasArgsValidator:      command.Args != nil,
		HasPreRun:             command.PreRun != nil || command.PreRunE != nil,
		HasPersistentPreRun:   command.PersistentPreRun != nil || command.PersistentPreRunE != nil,
		HasPostRun:            command.PostRun != nil || command.PostRunE != nil,
		HasPersistentPostRun:  command.PersistentPostRun != nil || command.PersistentPostRunE != nil,
		Annotations:           copyStringMap(command.Annotations),
	}
	appendFlags(command.LocalNonPersistentFlags(), "local", &contract.Flags)
	appendFlags(command.PersistentFlags(), "persistent", &contract.Flags)
	sort.Slice(contract.Flags, func(i, j int) bool {
		if contract.Flags[i].Scope != contract.Flags[j].Scope {
			return contract.Flags[i].Scope < contract.Flags[j].Scope
		}
		return contract.Flags[i].Name < contract.Flags[j].Name
	})
	*output = append(*output, contract)
	for _, child := range command.Commands() {
		appendCommands(child, output)
	}
}

func appendFlags(flags *pflag.FlagSet, scope string, output *[]flagContract) {
	flags.VisitAll(func(flag *pflag.Flag) {
		annotations := make(map[string][]string, len(flag.Annotations))
		for key, values := range flag.Annotations {
			copied := append([]string(nil), values...)
			sort.Strings(copied)
			annotations[key] = copied
		}
		*output = append(*output, flagContract{
			Scope:               scope,
			Name:                flag.Name,
			Shorthand:           flag.Shorthand,
			Type:                flag.Value.Type(),
			Default:             flag.DefValue,
			NoOptDefault:        flag.NoOptDefVal,
			Usage:               flag.Usage,
			Hidden:              flag.Hidden,
			Deprecated:          flag.Deprecated,
			ShorthandDeprecated: flag.ShorthandDeprecated,
			Annotations:         annotations,
		})
	})
}

func copyStringMap(input map[string]string) map[string]string {
	if input == nil {
		return map[string]string{}
	}
	output := make(map[string]string, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
