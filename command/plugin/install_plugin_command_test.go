package plugin_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/pluginaction"
	"code.cloudfoundry.org/cli/command/commandfakes"
	. "code.cloudfoundry.org/cli/command/plugin"
	"code.cloudfoundry.org/cli/command/plugin/pluginfakes"
	"code.cloudfoundry.org/cli/command/plugin/shared"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("install-plugin command", func() {
	var (
		cmd        InstallPluginCommand
		testUI     *ui.UI
		input      *Buffer
		fakeConfig *commandfakes.FakeConfig
		fakeActor  *pluginfakes.FakeInstallPluginActor
		executeErr error
	)

	BeforeEach(func() {
		input = NewBuffer()
		testUI = ui.NewTestUI(input, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeActor = new(pluginfakes.FakeInstallPluginActor)

		cmd = InstallPluginCommand{
			UI:     testUI,
			Config: fakeConfig,
			Actor:  fakeActor,
		}

		fakeConfig.ExperimentalReturns(true)
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	Context("installing from a local file", func() {
		BeforeEach(func() {
			cmd.OptionalArgs.LocalPath = "some-path"
		})

		Context("when the -f argument is given", func() {
			BeforeEach(func() {
				cmd.Force = true
			})

			Context("when the file exists", func() {
				BeforeEach(func() {
					fakeActor.FileExistsReturns(true)
				})

				Context("when the actor does not return an error", func() {
					BeforeEach(func() {
						fakeActor.InstallPluginFromPathReturns(pluginaction.Plugin{
							Name:    "some-plugin",
							Version: "1.0.0",
						}, nil)
					})

					It("prints the warnings installs the plugin", func() {
						Expect(executeErr).ToNot(HaveOccurred())

						Expect(testUI.Out).To(Say("Attention: Plugins are binaries written by potentially untrusted authors\\."))
						Expect(testUI.Out).To(Say("Install and use plugins at your own risk\\."))
						Expect(testUI.Out).To(Say("Installing plugin some-path..."))
						Expect(testUI.Out).To(Say("OK"))
						Expect(testUI.Out).To(Say("Plugin some-plugin 1\\.0\\.0 successfully installed\\."))

						Expect(fakeActor.FileExistsCallCount()).To(Equal(1))
						Expect(fakeActor.FileExistsArgsForCall(0)).To(Equal("some-path"))
						Expect(fakeActor.InstallPluginFromPathCallCount()).To(Equal(1))
						Expect(fakeActor.InstallPluginFromPathArgsForCall(0)).To(Equal("some-path"))
					})
				})

				Context("when the actor returns an error", func() {
					var expectedErr error

					BeforeEach(func() {
						expectedErr = errors.New("some error")
						fakeActor.InstallPluginFromPathReturns(pluginaction.Plugin{}, expectedErr)
					})

					It("prints warnings and returns the error", func() {
						Expect(executeErr).To(MatchError(expectedErr))

						Expect(testUI.Out).To(Say("Attention: Plugins are binaries written by potentially untrusted authors\\."))
						Expect(testUI.Out).To(Say("Install and use plugins at your own risk\\."))
						Expect(testUI.Out).To(Say("Installing plugin some-path\\.\\.\\."))
					})
				})
			})

			Context("when the file does not exist", func() {
				BeforeEach(func() {
					fakeActor.FileExistsReturns(false)
				})

				It("does not print the plugin warning and returns a FileNotFound error", func() {
					Expect(executeErr).To(MatchError(shared.FileNotFoundError{
						Path: "some-path",
					}))

					Expect(testUI.Out).ToNot(Say("Attention: Plugins are binaries written by potentially untrusted authors\\."))
					Expect(testUI.Out).ToNot(Say("Installing plugin some-path\\.\\.\\."))
				})
			})
		})

		Context("when the -f argument is not given", func() {
			BeforeEach(func() {
				cmd.Force = false
				fakeActor.FileExistsReturns(true)
				fakeActor.InstallPluginFromPathReturns(pluginaction.Plugin{
					Name:    "some-plugin",
					Version: "1.0.0",
				}, nil)
			})

			Context("when the user chooses yes", func() {
				BeforeEach(func() {
					input.Write([]byte("y\n"))
				})

				It("prints the warnings installs the plugin", func() {
					Expect(executeErr).ToNot(HaveOccurred())

					Expect(testUI.Out).To(Say("Attention: Plugins are binaries written by potentially untrusted authors\\."))
					Expect(testUI.Out).To(Say("Install and use plugins at your own risk\\."))
					Expect(testUI.Out).To(Say("Do you want to install the plugin some-path\\? \\[yN\\]"))
					Expect(testUI.Out).To(Say("Installing plugin some-path\\.\\.\\."))
					Expect(testUI.Out).To(Say("OK"))
					Expect(testUI.Out).To(Say("Plugin some-plugin 1\\.0\\.0 successfully installed\\."))

					Expect(fakeActor.FileExistsCallCount()).To(Equal(1))
					Expect(fakeActor.FileExistsArgsForCall(0)).To(Equal("some-path"))
					Expect(fakeActor.InstallPluginFromPathCallCount()).To(Equal(1))
					Expect(fakeActor.InstallPluginFromPathArgsForCall(0)).To(Equal("some-path"))
				})
			})

			Context("when the user chooses no", func() {
				BeforeEach(func() {
					input.Write([]byte("n\n"))
				})

				It("cancels plugin installation", func() {
					Expect(executeErr).To(MatchError(shared.PluginInstallationCancelled{}))
				})
			})

			Context("when the user chooses the default", func() {
				BeforeEach(func() {
					input.Write([]byte("\n"))
				})

				It("cancels plugin installation", func() {
					Expect(executeErr).To(MatchError(shared.PluginInstallationCancelled{}))
				})
			})

			Context("when the user input is invalid", func() {
				BeforeEach(func() {
					input.Write([]byte("e\n"))
				})

				It("returns an error", func() {
					Expect(executeErr).To(HaveOccurred())

					Expect(fakeActor.InstallPluginFromPathCallCount()).To(Equal(0))
				})
			})
		})
	})
})
