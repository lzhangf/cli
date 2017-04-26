package shared

import "code.cloudfoundry.org/cli/actor/pluginaction"

func HandleError(err error) error {
	switch e := err.(type) {
	case pluginaction.PluginNotFoundError:
		return PluginNotFoundError{Name: e.Name}
	case pluginaction.GettingPluginRepositoryError:
		return GettingPluginRepositoryError{Name: e.Name, Message: e.Message}
	case pluginaction.PluginInvalidError:
		return PluginInvalidError{Path: e.Path, WrappedErrMessage: e.Error()}
	case pluginaction.PluginCommandConflictError:
		return PluginCommandConflictError{PluginName: e.PluginName, PluginVersion: e.PluginVersion, CommandName: e.CommandName, WrappedErrMessage: e.Error()}
	case pluginaction.PluginAliasConflictError:
		return PluginAliasConflictError{PluginName: e.PluginName, PluginVersion: e.PluginVersion, CommandAlias: e.CommandAlias, WrappedErrMessage: e.Error()}
	}
	return err
}
