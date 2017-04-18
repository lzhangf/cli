package pushaction_test

import (
	"errors"

	. "code.cloudfoundry.org/cli/actor/pushaction"
	"code.cloudfoundry.org/cli/actor/pushaction/manifest"
	"code.cloudfoundry.org/cli/actor/pushaction/pushactionfakes"
	"code.cloudfoundry.org/cli/actor/v2action"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application Config", func() {
	var (
		actor       *Actor
		fakeV2Actor *pushactionfakes.FakeV2Actor
	)

	BeforeEach(func() {
		fakeV2Actor = new(pushactionfakes.FakeV2Actor)
		actor = NewActor(fakeV2Actor)
	})

	Describe("Apply", func() {
		var (
			eventStream    <-chan Event
			warningsStream <-chan Warnings
			errorStream    <-chan error

			config ApplicationConfig
		)
		BeforeEach(func() {
			config = ApplicationConfig{
				DesiredApplication: v2action.Application{
					Name:      "some-app-name",
					SpaceGUID: "some-space-guid",
				},
			}
		})

		JustBeforeEach(func() {
			eventStream, warningsStream, errorStream = actor.Apply(config)
		})

		AfterEach(func() {
			Eventually(warningsStream).Should(BeClosed())
			Eventually(eventStream).Should(BeClosed())
			Eventually(errorStream).Should(BeClosed())
		})

		Context("when the app exists", func() {
			BeforeEach(func() {
				config.CurrentApplication = v2action.Application{
					Name:      "some-app-name",
					GUID:      "some-app-guid",
					SpaceGUID: "some-space-guid",
					Buildpack: "java",
				}
				config.DesiredApplication = v2action.Application{
					Name:      "some-app-name",
					GUID:      "some-app-guid",
					SpaceGUID: "some-space-guid",
					Buildpack: "ruby",
				}
			})

			Context("when the update is successful", func() {
				BeforeEach(func() {
					fakeV2Actor.UpdateApplicationReturns(v2action.Application{}, v2action.Warnings{"update-warning"}, nil)
				})

				It("updates the application", func() {
					Eventually(warningsStream).Should(Receive(ConsistOf("update-warning")))
					Eventually(eventStream).Should(Receive(Equal(ApplicationUpdated)))
					Eventually(eventStream).Should(Receive(Equal(Complete)))

					Expect(fakeV2Actor.UpdateApplicationCallCount()).To(Equal(1))
					Expect(fakeV2Actor.UpdateApplicationArgsForCall(0)).To(Equal(v2action.Application{
						Name:      "some-app-name",
						GUID:      "some-app-guid",
						SpaceGUID: "some-space-guid",
						Buildpack: "ruby",
					}))
				})
			})

			Context("when the update errors", func() {
				var expectedErr error
				BeforeEach(func() {
					expectedErr = errors.New("oh my")
					fakeV2Actor.UpdateApplicationReturns(v2action.Application{}, v2action.Warnings{"update-warning"}, expectedErr)
				})

				It("returns warnings and error and stops", func() {
					Eventually(warningsStream).Should(Receive(ConsistOf("update-warning")))
					Eventually(errorStream).Should(Receive(MatchError(expectedErr)))
					Consistently(eventStream).ShouldNot(Receive(Equal(ApplicationUpdated)))
				})
			})
		})

		Context("when the app does not exist", func() {
			Context("when the creation is successful", func() {
				BeforeEach(func() {
					fakeV2Actor.CreateApplicationReturns(v2action.Application{}, v2action.Warnings{"create-warning"}, nil)
				})

				It("creates the application", func() {
					Eventually(warningsStream).Should(Receive(ConsistOf("create-warning")))
					Eventually(eventStream).Should(Receive(Equal(ApplicationCreated)))
					Eventually(eventStream).Should(Receive(Equal(Complete)))

					Expect(fakeV2Actor.CreateApplicationCallCount()).To(Equal(1))
					Expect(fakeV2Actor.CreateApplicationArgsForCall(0)).To(Equal(v2action.Application{
						Name:      "some-app-name",
						SpaceGUID: "some-space-guid",
					}))
				})
			})

			Context("when the creation errors", func() {
				var expectedErr error
				BeforeEach(func() {
					expectedErr = errors.New("oh my")
					fakeV2Actor.CreateApplicationReturns(v2action.Application{}, v2action.Warnings{"create-warning"}, expectedErr)
				})

				It("returns warnings and error and stops", func() {
					Eventually(warningsStream).Should(Receive(ConsistOf("create-warning")))
					Eventually(errorStream).Should(Receive(MatchError(expectedErr)))
					Consistently(eventStream).ShouldNot(Receive(Equal(ApplicationCreated)))
				})
			})
		})
	})

	Describe("ConvertToApplicationConfig", func() {
		var (
			orgGUID      string
			spaceGUID    string
			manifestApps []manifest.Application

			configs    []ApplicationConfig
			warnings   Warnings
			executeErr error
		)

		BeforeEach(func() {
			orgGUID = "some-org-guid"
			spaceGUID = "some-space-guid"
			manifestApps = []manifest.Application{{
				Name: "some-app",
				Path: "some-path",
			}}
		})

		JustBeforeEach(func() {
			configs, warnings, executeErr = actor.ConvertToApplicationConfig(orgGUID, spaceGUID, manifestApps)
		})

		Context("when the application exists", func() {
			var app v2action.Application

			BeforeEach(func() {
				app = v2action.Application{
					Name:      "some-app",
					GUID:      "some-app-guid",
					SpaceGUID: spaceGUID,
				}

				fakeV2Actor.GetApplicationByNameAndSpaceReturns(app, v2action.Warnings{"some-app-warning-1", "some-app-warning-2"}, nil)
			})

			It("sets the current and desired application to the current", func() {
				Expect(executeErr).ToNot(HaveOccurred())
				Expect(warnings).To(ConsistOf("some-app-warning-1", "some-app-warning-2"))
				Expect(configs).To(Equal([]ApplicationConfig{{
					CurrentApplication: app,
					DesiredApplication: app,
					TargetedSpaceGUID:  spaceGUID,
					Path:               "some-path",
				}}))
			})
		})

		Describe("when the application does not exist", func() {
			BeforeEach(func() {
				fakeV2Actor.GetApplicationByNameAndSpaceReturns(v2action.Application{}, v2action.Warnings{"some-app-warning-1", "some-app-warning-2"}, v2action.ApplicationNotFoundError{})
			})

			It("creates a new application and sets it to the desired application", func() {
				Expect(executeErr).ToNot(HaveOccurred())
				Expect(warnings).To(ConsistOf("some-app-warning-1", "some-app-warning-2"))
				Expect(configs).To(Equal([]ApplicationConfig{{
					DesiredApplication: v2action.Application{
						Name:      "some-app",
						SpaceGUID: spaceGUID,
					},
					TargetedSpaceGUID: spaceGUID,
					Path:              "some-path",
				}}))
			})
		})
	})
})
