package common

import (
	"os"

	"code.cloudfoundry.org/cli/actor/pluginaction"
	oldCmd "code.cloudfoundry.org/cli/cf/cmd"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/plugin/shared"
	"code.cloudfoundry.org/cli/util/configv3"
)

//go:generate counterfeiter . InstallPluginActor

type InstallPluginActor interface {
	FileExists(path string) bool
	ValidatePlugin(metadata pluginaction.PluginMetadata, commands pluginaction.CommandList, path string) (configv3.Plugin, error)
	UninstallPlugin(uninstaller pluginaction.PluginUninstaller, name string) error
	InstallPluginFromPath(path string, plugin configv3.Plugin) error
}

type InstallPluginCommand struct {
	OptionalArgs         flag.InstallPluginArgs `positional-args:"yes"`
	Force                bool                   `short:"f" description:"Force install of plugin without confirmation"`
	RegisteredRepository string                 `short:"r" description:"Name of a registered repository where the specified plugin is located"`
	usage                interface{}            `usage:"CF_NAME install-plugin (LOCAL-PATH/TO/PLUGIN | URL | -r REPO_NAME PLUGIN_NAME) [-f]\n\nEXAMPLES:\n   CF_NAME install-plugin ~/Downloads/plugin-foobar\n   CF_NAME install-plugin https://example.com/plugin-foobar_linux_amd64\n   CF_NAME install-plugin -r My-Repo plugin-echo"`
	relatedCommands      interface{}            `related_commands:"add-plugin-repo, list-plugin-repos, plugins"`

	UI     command.UI
	Config command.Config
	Actor  InstallPluginActor
}

func (cmd *InstallPluginCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.Actor = pluginaction.NewActor(config, nil)
	return nil
}

func (cmd InstallPluginCommand) Execute(args []string) error {
	//TODO: remove me
	if !cmd.Config.Experimental() {
		oldCmd.Main(os.Getenv("CF_TRACE"), os.Args)
		return nil
	}
	cmd.UI.DisplayText(command.ExperimentalWarning)

	var plugin configv3.Plugin
	pluginPath := string(cmd.OptionalArgs.LocalPath)

	if pluginPath != "" {
		if !cmd.Actor.FileExists(pluginPath) {
			return shared.FileNotFoundError{Path: pluginPath}
		}

		cmd.UI.DisplayText("Attention: Plugins are binaries written by potentially untrusted authors.")
		cmd.UI.DisplayText("Install and use plugins at your own risk.")

		if !cmd.Force {
			really, promptErr := cmd.UI.DisplayBoolPrompt(false, "Do you want to install the plugin {{.Path}}?", map[string]interface{}{
				"Path": pluginPath,
			})
			if promptErr != nil {
				return promptErr
			}
			if !really {
				return shared.PluginInstallationCancelled{}
			}
		}

		var err error
		plugin, err = cmd.Actor.ValidatePlugin(nil, Commands, pluginPath)
		if e, ok := err.(pluginaction.PluginAlreadyInstalledError); ok {
			if !cmd.Force {
				return shared.PluginAlreadyInstalledError{
					Name:              e.Name,
					Version:           e.Version,
					Path:              pluginPath,
					WrappedErrMessage: e.Error(),
				}
			}

			cmd.UI.DisplayText("Plugin {{.Name}} {{.Version}} is already installed. Uninstalling existing plugin...", map[string]interface{}{
				"Name":    plugin.Name,
				"Version": plugin.Version,
			})

			uninstallErr := cmd.Actor.UninstallPlugin(nil, plugin.Name)
			if uninstallErr != nil {
				return uninstallErr
			}

			cmd.UI.DisplayOK()
			cmd.UI.DisplayText("Plugin {{.Name}} successfully uninstalled.", map[string]interface{}{
				"Name": plugin.Name,
			})
		} else if err != nil {
			return shared.HandleError(err)
		}

		cmd.UI.DisplayTextWithFlavor("Installing plugin {{.Name}}...", map[string]interface{}{
			"Name": plugin.Name,
		})

		installErr := cmd.Actor.InstallPluginFromPath(pluginPath, plugin)
		if installErr != nil {
			return installErr
		}

		cmd.UI.DisplayOK()
		cmd.UI.DisplayTextWithFlavor("Plugin {{.Name}} {{.Version}} successfully installed.", map[string]interface{}{
			"Name":    plugin.Name,
			"Version": plugin.Version,
		})
	}

	return nil
}
