package shared

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
	return "File not found locally, make sure the file exists at given path {{.filepath}}"
}

func (e FileNotFoundError) Translate(translate func(string, ...interface{}) string) string {
	return translate(e.Error(), map[string]interface{}{
		"filepath": e.Path,
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
