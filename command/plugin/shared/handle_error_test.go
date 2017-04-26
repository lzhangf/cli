package shared_test

import (
	"errors"

	"code.cloudfoundry.org/cli/actor/pluginaction"
	. "code.cloudfoundry.org/cli/command/plugin/shared"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("HandleError", func() {
	err := errors.New("some-error")

	DescribeTable("error translations",
		func(passedInErr error, expectedErr error) {
			actualErr := HandleError(passedInErr)
			Expect(actualErr).To(MatchError(expectedErr))
		},

		Entry("pluginaction.PluginNotFoundError -> PluginNotFoundError",
			pluginaction.PluginNotFoundError{Name: "some-plugin"},
			PluginNotFoundError{Name: "some-plugin"}),
		Entry("pluginaction.GettingPluginRepositoryError -> GettingPluginRepositoryError",
			pluginaction.GettingPluginRepositoryError{Name: "some-repo", Message: "404"},
			GettingPluginRepositoryError{Name: "some-repo", Message: "404"}),
		Entry("pluginaction.PluginInvalidError -> PluginInvalidError",
			pluginaction.PluginInvalidError{Path: "some-path"},
			PluginInvalidError{
				Path:              "some-path",
				WrappedErrMessage: pluginaction.PluginInvalidError{}.Error(),
			}),
		Entry("pluginaction.PluginCommandConflictError -> PluginCommandConflictError",
			pluginaction.PluginCommandConflictError{
				PluginName:    "some-plugin",
				PluginVersion: "1.1.1",
				CommandName:   "some-command",
			},
			PluginCommandConflictError{
				PluginName:        "some-plugin",
				PluginVersion:     "1.1.1",
				CommandName:       "some-command",
				WrappedErrMessage: pluginaction.PluginCommandConflictError{}.Error(),
			}),
		Entry("pluginaction.PluginAliasConflictError -> PluginAliasConflictError",
			pluginaction.PluginAliasConflictError{
				PluginName:    "some-plugin",
				PluginVersion: "1.1.1",
				CommandAlias:  "aa",
			},
			PluginAliasConflictError{
				PluginName:        "some-plugin",
				PluginVersion:     "1.1.1",
				CommandAlias:      "aa",
				WrappedErrMessage: pluginaction.PluginAliasConflictError{}.Error(),
			}),

		Entry("default case -> original error",
			err,
			err),
	)

	It("returns nil for a nil error", func() {
		nilErr := HandleError(nil)
		Expect(nilErr).To(BeNil())
	})
})
