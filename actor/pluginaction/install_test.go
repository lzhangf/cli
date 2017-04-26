package pluginaction_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "code.cloudfoundry.org/cli/actor/pluginaction"
	"code.cloudfoundry.org/cli/actor/pluginaction/pluginactionfakes"
	"code.cloudfoundry.org/cli/util/configv3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("install actions", func() {
	var (
		actor      Actor
		fakeConfig *pluginactionfakes.FakeConfig
		pluginPath string
	)

	BeforeEach(func() {
		fakeConfig = new(pluginactionfakes.FakeConfig)
		actor = NewActor(fakeConfig, nil)

		pluginFile, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())
		err = pluginFile.Close()
		Expect(err).NotTo(HaveOccurred())

		pluginPath = pluginFile.Name()
	})

	AfterEach(func() {
		err := os.Remove(pluginPath)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("FileExists", func() {
		Context("when the file exists", func() {
			It("returns true", func() {
				Expect(actor.FileExists(pluginPath)).To(BeTrue())
			})
		})

		Context("when the file does not exist", func() {
			It("returns false", func() {
				Expect(actor.FileExists("/some/path/that/does/not/exist")).To(BeFalse())
			})
		})
	})

	Describe("ValidatePlugin", func() {
		var (
			fakePluginMetadata *pluginactionfakes.FakePluginMetadata
			fakeCommandList    *pluginactionfakes.FakeCommandList
			plugin             configv3.Plugin
			validateErr        error
		)

		BeforeEach(func() {
			fakePluginMetadata = new(pluginactionfakes.FakePluginMetadata)
			fakeCommandList = new(pluginactionfakes.FakeCommandList)
		})

		JustBeforeEach(func() {
			plugin, validateErr = actor.ValidatePlugin(fakePluginMetadata, fakeCommandList, pluginPath)
		})

		Context("when getting the plugin metadata returns an error", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("error getting metadata")
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{}, expectedErr)
			})

			It("returns the error", func() {
				Expect(validateErr).To(MatchError(expectedErr))
			})
		})

		Context("when the plugin name is missing", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{}, nil)
			})

			It("returns a PluginInvalidError", func() {
				Expect(validateErr).To(MatchError(PluginInvalidError{Path: pluginPath}))
			})
		})

		Context("when the plugin does not have any commands", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{Name: "some-plugin"}, nil)
			})

			It("returns a PluginInvalidError", func() {
				Expect(validateErr).To(MatchError(PluginInvalidError{Path: pluginPath}))
			})
		})

		Context("when the plugin has a command name that conflicts with a native command", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{Name: "some-command"},
						{Name: "version"},
					},
				}, nil)

				count := 0
				fakeCommandList.HasCommandStub = func(commandName string) bool {
					if count == 0 {
						count++
						return false
					}
					return true
				}
			})

			It("returns a PluginCommandConflictError", func() {
				Expect(validateErr).To(MatchError(PluginCommandConflictError{
					PluginName:    "some-plugin",
					PluginVersion: "1.1.1",
					CommandName:   "version",
				}))
			})
		})

		Context("when the plugin has a command alias that conflicts with a native command alias", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{
							Name:  "some-command",
							Alias: "sc",
						},
						{
							Name:  "version",
							Alias: "v",
						},
					},
				}, nil)

				count := 0
				fakeCommandList.HasAliasStub = func(aliasName string) bool {
					if count == 0 {
						count++
						return false
					}
					return true
				}
			})

			It("returns a PluginAliasConflictError", func() {
				Expect(validateErr).To(MatchError(PluginAliasConflictError{
					PluginName:    "some-plugin",
					PluginVersion: "1.1.1",
					CommandAlias:  "v",
				}))
			})
		})

		Context("when the plugin has a command name that conflicts with an existing plugin command name", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{Name: "some-command"},
						{Name: "duplicate-command"},
					},
				}, nil)

				fakeConfig.PluginsReturns([]configv3.Plugin{
					{
						Name: "installed-plugin-1",
						Commands: []configv3.PluginCommand{
							{Name: "unique-command-1"},
						},
					},
					{
						Name: "installed-plugin-2",
						Commands: []configv3.PluginCommand{
							{Name: "unique-command-2"},
							{Name: "duplicate-command"},
						},
					},
				})
			})

			It("returns a PluginCommandConflictError", func() {
				Expect(validateErr).To(MatchError(PluginCommandConflictError{
					PluginName:    "some-plugin",
					PluginVersion: "1.1.1",
					CommandName:   "duplicate-command",
				}))
			})
		})

		Context("when the plugin has a command alias that conflicts with an existing plugin command alias", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{
							Name:  "some-command",
							Alias: "sc",
						},
						{
							Name:  "non-unique-alias",
							Alias: "uc3",
						},
					},
				}, nil)

				fakeConfig.PluginsReturns([]configv3.Plugin{
					{
						Name: "installed-plugin-1",
						Commands: []configv3.PluginCommand{
							{
								Name:  "unique-command-1",
								Alias: "uc1",
							},
						},
					},
					{
						Name: "installed-plugin-2",
						Commands: []configv3.PluginCommand{
							{
								Name:  "unique-command-2",
								Alias: "uc2",
							},
							{
								Name:  "unique-command-3",
								Alias: "uc3",
							},
						},
					},
				})
			})

			It("returns a PluginAliasConflictError", func() {
				Expect(validateErr).To(MatchError(PluginAliasConflictError{
					PluginName:    "some-plugin",
					PluginVersion: "1.1.1",
					CommandAlias:  "uc3",
				}))
			})
		})

		Context("when the plugin name conflicts with an existing plugin name", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-existing-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{
							Name:  "some-command",
							Alias: "sc",
						},
					},
				}, nil)

				fakeConfig.GetPluginReturns(configv3.Plugin{Name: "some-existing-plugin"}, true)
			})

			It("returns a PluginAlreadyInstalledError", func() {
				Expect(validateErr).To(MatchError(PluginAlreadyInstalledError{
					Name:    "some-existing-plugin",
					Version: "1.1.1",
				}))
			})
		})

		Context("when the plugin is valid", func() {
			BeforeEach(func() {
				fakePluginMetadata.GetMetadataReturns(configv3.Plugin{
					Name: "some-plugin",
					Version: configv3.PluginVersion{
						Major: 1,
						Minor: 1,
						Build: 1,
					},
					Commands: []configv3.PluginCommand{
						{
							Name:  "some-command",
							Alias: "sc",
						},
						{
							Name:  "some-other-command",
							Alias: "soc",
						},
					},
				}, nil)

				fakeConfig.PluginsReturns([]configv3.Plugin{
					{
						Name: "installed-plugin-1",
						Commands: []configv3.PluginCommand{
							{
								Name:  "unique-command-1",
								Alias: "uc1",
							},
						},
					},
					{
						Name: "installed-plugin-2",
						Commands: []configv3.PluginCommand{
							{
								Name:  "unique-command-2",
								Alias: "uc2",
							},
							{
								Name:  "unique-command-3",
								Alias: "uc3",
							},
						},
					},
				})
			})

			It("returns nil", func() {
				Expect(validateErr).ToNot(HaveOccurred())

				Expect(fakePluginMetadata.GetMetadataCallCount()).To(Equal(1))
				Expect(fakePluginMetadata.GetMetadataArgsForCall(0)).To(Equal(pluginPath))

				Expect(fakeCommandList.HasCommandCallCount()).To(Equal(2))
				Expect(fakeCommandList.HasCommandArgsForCall(0)).To(Equal("some-command"))
				Expect(fakeCommandList.HasCommandArgsForCall(1)).To(Equal("some-other-command"))

				Expect(fakeCommandList.HasAliasCallCount()).To(Equal(2))
				Expect(fakeCommandList.HasAliasArgsForCall(0)).To(Equal("sc"))
				Expect(fakeCommandList.HasAliasArgsForCall(1)).To(Equal("soc"))

				Expect(fakeConfig.PluginsCallCount()).To(Equal(1))

				Expect(fakeConfig.GetPluginCallCount()).To(Equal(1))
				Expect(fakeConfig.GetPluginArgsForCall(0)).To(Equal("some-plugin"))
			})
		})
	})

	Describe("InstallPluginFromLocalPath", func() {
		var (
			plugin     configv3.Plugin
			installErr error
		)

		BeforeEach(func() {
			plugin = configv3.Plugin{
				Name: "some-plugin",
				Commands: []configv3.PluginCommand{
					{Name: "some-command"},
				},
			}
		})

		JustBeforeEach(func() {
			installErr = actor.InstallPluginFromPath(pluginPath, plugin)
		})

		Context("when an error is encountered copying the plugin to the plugin directory", func() {
			BeforeEach(func() {
				fakeConfig.PluginHomeReturns(pluginPath)
			})

			It("returns the error", func() {
				Expect(installErr).To(MatchError(MatchRegexp("not a directory")))
			})
		})

		Context("when an error is encountered writing the plugin config to disk", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = errors.New("write config error")
				fakeConfig.WritePluginConfigReturns(expectedErr)
			})

			It("returns the error", func() {
				Expect(installErr).To(MatchError(expectedErr))
			})
		})

		Context("when no errors are encountered", func() {
			var (
				tempDir       string
				pluginHomeDir string
			)

			BeforeEach(func() {
				tempDir, err := ioutil.TempDir("", "")
				Expect(err).ToNot(HaveOccurred())

				pluginHomeDir = filepath.Join(tempDir, ".cf", "plugin")
				fakeConfig.PluginHomeReturns(pluginHomeDir)
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDir)
				Expect(err).ToNot(HaveOccurred())
			})

			It("makes an executable copy of the plugin file in the plugin directory, updates the plugin config, and writes the config to disk", func() {
				Expect(installErr).ToNot(HaveOccurred())

				pluginInfo, err := os.Stat(filepath.Join(pluginHomeDir, filepath.Base(pluginPath)))
				Expect(err).ToNot(HaveOccurred())
				// The filemode of the original plugin file is 0600
				Expect(pluginInfo.Mode()).To(Equal(os.FileMode(0700)))

				Expect(fakeConfig.PluginHomeCallCount()).To(Equal(1))

				Expect(fakeConfig.AddPluginCallCount()).To(Equal(1))
				Expect(fakeConfig.AddPluginArgsForCall(0)).To(Equal(plugin))

				Expect(fakeConfig.WritePluginConfigCallCount()).To(Equal(1))
			})
		})
	})
})
