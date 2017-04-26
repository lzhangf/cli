package shared

import "fmt"

type PluginNotFoundError struct {
	Name string
}

func (e PluginNotFoundError) Error() string {
	return "Plugin {{.Name}} does not exist."
}

func (e PluginNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Name": e.Name,
	})
}

type NoPluginRepositoriesError struct{}

func (e NoPluginRepositoriesError) Error() string {
	return "No plugin repositories registered to search for plugin updates."
}

func (e NoPluginRepositoriesError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error())
}

// GettingPluginRepositoryError is returned when there's an error
// accessing the plugin repository
type GettingPluginRepositoryError struct {
	Name    string
	Message string
}

func (e GettingPluginRepositoryError) Error() string {
	return "Could not get plugin repository '{{.RepositoryName}}': {{.ErrorMessage}}"
}

func (e GettingPluginRepositoryError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{"RepositoryName": e.Name, "ErrorMessage": e.Message})
}

// FileNotFoundError is returned when a local plugin binary is not found during
// installation.
type FileNotFoundError struct {
	Path string
}

func (e FileNotFoundError) Error() string {
	return "File not found locally, make sure the file exists at given path {{.FilePath}}"
}

func (e FileNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"FilePath": e.Path,
	})
}

// PluginInstallationCancelled is returned when the user chooses no during
// plugin installation confirmation.
type PluginInstallationCancelled struct {
}

func (e PluginInstallationCancelled) Error() string {
	return "Plugin installation cancelled"
}

func (e PluginInstallationCancelled) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error())
}

// PluginInvalidError is returned with a plugin is invalid because it is
// missing a name or has 0 commands.
type PluginInvalidError struct {
	Path              string
	WrappedErrMessage string
}

func (e PluginInvalidError) Error() string {
	return e.WrappedErrMessage
}

func (e PluginInvalidError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Path": e.Path,
	})
}

// PluginCommandConflictError is returned when a plugin command name conflicts
// with a native or existing plugin command name.
type PluginCommandConflictError struct {
	PluginName        string
	PluginVersion     string
	CommandName       string
	WrappedErrMessage string
}

func (e PluginCommandConflictError) Error() string {
	return e.WrappedErrMessage
}

func (e PluginCommandConflictError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"PluginName":    e.PluginName,
		"PluginVersion": e.PluginVersion,
		"CommandName":   e.CommandName,
	})
}

// PluginAliasConflictError is returned when a plugin command alias conflicts
// with a core or existing plugin command alias.
type PluginAliasConflictError struct {
	PluginName        string
	PluginVersion     string
	CommandAlias      string
	WrappedErrMessage string
}

func (e PluginAliasConflictError) Error() string {
	return e.WrappedErrMessage
}

func (e PluginAliasConflictError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"PluginName":    e.PluginName,
		"PluginVersion": e.PluginVersion,
		"CommandAlias":  e.CommandAlias,
	})
}

// PluginAlreadyInstalledError is returned when the plugin has the same name as
// an installed plugin.
type PluginAlreadyInstalledError struct {
	Name              string
	Version           string
	Path              string
	WrappedErrMessage string
}

func (e PluginAlreadyInstalledError) Error() string {
	return fmt.Sprintf("%s\nTIP: Use '{{.Command}}' to force a reinstall.", e.WrappedErrMessage)
}

func (e PluginAlreadyInstalledError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"Name":    e.Name,
		"Version": e.Version,
		"Command": fmt.Sprintf("cf install-plugin %s -f", e.Path),
	})
}
