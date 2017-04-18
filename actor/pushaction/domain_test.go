package pushaction_test

import (
	. "code.cloudfoundry.org/cli/actor/pushaction"
	"code.cloudfoundry.org/cli/actor/pushaction/pushactionfakes"
	"code.cloudfoundry.org/cli/actor/v2action"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Domains", func() {
	var (
		actor       *Actor
		fakeV2Actor *pushactionfakes.FakeV2Actor
	)

	BeforeEach(func() {
		fakeV2Actor = new(pushactionfakes.FakeV2Actor)
		actor = NewActor(fakeV2Actor)
	})

	Describe("DefaultDomain", func() {
		var (
			orgGUID       string
			defaultDomain v2action.Domain
			warnings      Warnings
			executeErr    error
		)

		BeforeEach(func() {
			orgGUID = "some-org-guid"
		})

		JustBeforeEach(func() {
			defaultDomain, warnings, executeErr = actor.DefaultDomain(orgGUID)
		})

		Context("when private domains exist", func() {
			It("returns the first private domain and warnings", func() {

			})
		})

		Context("when only shared domains exist", func() {
			It("returns the first shared domain and warnings", func() {

			})
		})

		Context("no domains exist", func() {
			It("returns the first shared domain and warnings", func() {

			})
		})
	})
})
