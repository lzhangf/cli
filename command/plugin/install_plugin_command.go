package plugin

import (
	"os"

	"code.cloudfoundry.org/cli/actor/pluginaction"
	oldCmd "code.cloudfoundry.org/cli/cf/cmd"

	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/plugin/shared"
)

//go:generate counterfeiter . InstallPluginActor
type InstallPluginActor interface {
	FileExists(path string) bool
	InstallPluginFromPath(path string) (pluginaction.Plugin, error)
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

	return nil
}

func (cmd InstallPluginCommand) Execute(args []string) error {
	if !cmd.Config.Experimental() {
		oldCmd.Main(os.Getenv("CF_TRACE"), os.Args)
		return nil
	}

	cmd.UI.DisplayText(command.ExperimentalWarning)

	var (
		plugin pluginaction.Plugin
		err    error
	)

	if cmd.OptionalArgs.LocalPath != "" {
		if !cmd.Actor.FileExists(string(cmd.OptionalArgs.LocalPath)) {
			return shared.FileNotFoundError{
				Path: string(cmd.OptionalArgs.LocalPath),
			}
		}

		cmd.UI.DisplayText("Attention: Plugins are binaries written by potentially untrusted authors.")
		cmd.UI.DisplayText("Install and use plugins at your own risk.")

		if !cmd.Force {
			really, promptErr := cmd.UI.DisplayBoolPrompt(false, "Do you want to install the plugin {{.Path}}?", map[string]interface{}{
				"Path": cmd.OptionalArgs.LocalPath,
			})
			if promptErr != nil {
				return promptErr
			}

			if !really {
				return shared.PluginInstallationCancelled{}
			}
		}

		cmd.UI.DisplayTextWithFlavor("Installing plugin {{.PluginPath}}...", map[string]interface{}{
			"PluginPath": cmd.OptionalArgs.LocalPath,
		})
		plugin, err = cmd.Actor.InstallPluginFromPath(string(cmd.OptionalArgs.LocalPath))
		if err != nil {
			return err
		}
	}

	cmd.UI.DisplayOK()
	cmd.UI.DisplayTextWithFlavor("Plugin {{.Name}} {{.Version}} successfully installed.", map[string]interface{}{
		"Name":    plugin.Name,
		"Version": plugin.Version,
	})

	return nil
}
