package pluginaction

import (
	"os"
	"path/filepath"

	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/gofileutils/fileutils"
)

//go:generate counterfeiter . PluginMetadata

type PluginMetadata interface {
	GetMetadata(pluginPath string) (configv3.Plugin, error)
}

//go:generate counterfeiter . CommandList

type CommandList interface {
	HasCommand(string) bool
	HasAlias(string) bool
}

// PluginInvalidError is returned with a plugin is invalid because it is
// missing a name or has 0 commands.
type PluginInvalidError struct {
	Path string
}

func (e PluginInvalidError) Error() string {
	return "File {{.Path}} is not a valid cf CLI plugin binary."
}

// PluginCommandConflictError is returned when a plugin command name conflicts
// with a core or existing plugin command name.
type PluginCommandConflictError struct {
	PluginName    string
	PluginVersion string
	CommandName   string
}

func (e PluginCommandConflictError) Error() string {
	return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with names that are already used: {{.CommandName}}."
}

// PluginAliasConflictError is returned when a plugin command alias conflicts
// with a core or existing plugin command alias.
type PluginAliasConflictError struct {
	PluginName    string
	PluginVersion string
	CommandAlias  string
}

func (e PluginAliasConflictError) Error() string {
	return "Plugin {{.PluginName}} v{{.PluginVersion}} could not be installed as it contains commands with aliases that are already used: {{.CommandAlias}}."
}

// PluginAlreadyInstalledError is returned when a plugin with the same name is
// already installed.
type PluginAlreadyInstalledError struct {
	Name    string
	Version string
}

func (e PluginAlreadyInstalledError) Error() string {
	return "Plugin {{.Name}} {{.Version}} could not be installed. A plugin with that name is already installed."
}

// FileExists returns true if the file exists. It returns false if the file
// doesn't exist or there is an error checking.
func (actor Actor) FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func (actor Actor) ValidatePlugin(pluginMetadata PluginMetadata, commandList CommandList, path string) (configv3.Plugin, error) {
	plugin, err := pluginMetadata.GetMetadata(path)
	if err != nil {
		return configv3.Plugin{}, err
	}

	if plugin.Name == "" || len(plugin.Commands) == 0 {
		return configv3.Plugin{}, PluginInvalidError{Path: path}
	}

	pluginName := plugin.Name
	pluginVersion := plugin.Version.String()
	installedPlugins := actor.config.Plugins()

	for _, command := range plugin.Commands {
		commandName := command.Name
		commandAlias := command.Alias

		if commandList.HasCommand(commandName) {
			return configv3.Plugin{}, PluginCommandConflictError{
				PluginName:    pluginName,
				PluginVersion: pluginVersion,
				CommandName:   commandName,
			}
		}

		if commandList.HasAlias(commandAlias) {
			return configv3.Plugin{}, PluginAliasConflictError{
				PluginName:    pluginName,
				PluginVersion: pluginVersion,
				CommandAlias:  commandAlias,
			}
		}

		for _, installedPlugin := range installedPlugins {
			for _, installedCommand := range installedPlugin.Commands {
				if commandName == installedCommand.Name {
					return configv3.Plugin{}, PluginCommandConflictError{
						PluginName:    pluginName,
						PluginVersion: pluginVersion,
						CommandName:   commandName,
					}
				}

				if commandAlias != "" && commandAlias == installedCommand.Alias {
					return configv3.Plugin{}, PluginAliasConflictError{
						PluginName:    pluginName,
						PluginVersion: pluginVersion,
						CommandAlias:  commandAlias,
					}
				}
			}
		}
	}

	if _, exist := actor.config.GetPlugin(pluginName); exist {
		return configv3.Plugin{}, PluginAlreadyInstalledError{Name: pluginName, Version: pluginVersion}
	}

	return plugin, nil
}

func (actor Actor) InstallPluginFromPath(path string, plugin configv3.Plugin) error {
	installPath := filepath.Join(actor.config.PluginHome(), filepath.Base(path))
	err := fileutils.CopyPathToPath(path, installPath)
	if err != nil {
		return err
	}

	actor.config.AddPlugin(plugin)

	err = actor.config.WritePluginConfig()
	if err != nil {
		return err
	}

	return nil
}
