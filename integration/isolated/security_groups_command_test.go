package isolated

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("security-groups command", func() {
	Describe("help", func() {
		Context("when --help flag is provided", func() {
			It("displays command usage to output", func() {
				session := helpers.CF("security-groups", "--help")
				Eventually(session.Out).Should(Say("NAME:"))
				Eventually(session.Out).Should(Say("security-groups - List all security groups"))
				Eventually(session.Out).Should(Say("USAGE:"))
				Eventually(session.Out).Should(Say("cf security-groups"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	FDescribe("Unrefactored command", func() {
		var (
			session *Session
		)

		JustBeforeEach(func() {
			session = helpers.CF("security-groups", "fooo")
		})

		Context("when no API endpoint is set", func() {
			BeforeEach(func() {
				helpers.UnsetAPI()
			})

			It("fails with no API endpoint set message", func() {
				session := helpers.CF("security-groups")
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Out).Should(Say("No API endpoint set\\. Use 'cf login' or 'cf api' to target an endpoint\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when not logged in", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
			})

			It("fails with not logged in message", func() {
				session := helpers.CF("security-groups")
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Out).Should(Say("Not logged in\\. Use 'cf login' to log in\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when too many arguments are provided", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
			})

			It("fails with too many arguments message", func() {
				session := helpers.CF("security-groups", "foooo")
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Out).Should(Say("Incorrect Usage\\. No argument required"))
				Eventually(session).Should(Exit(1))
			})
		})

		FContext("when the environment is set-up correctly", func() {
			var (
				username string
			)

			BeforeEach(func() {
				username = helpers.LoginCF()
			})

			Context("when there are no security groups", func() {
				// _Can_ there be no security groups?
				It("lists no security groups", func() {
					session = helpers.CF("security-groups")
					Eventually(session.Out).Should(Say("Getting security groups as admin"))
					//Eventually(session.Out).Should(Say("No security groups"))
					Eventually(session.Out).Should(Say("OK"))
					Eventually(session).Should(Exit(0))
				})
			})

			FContext("when there are security groups", func() {
				var (
					securityGroup1 helpers.SecurityGroup
					securityGroup2 helpers.SecurityGroup
				)

				BeforeEach(func() {
					securityGroup1 = helpers.NewSecurityGroup(helpers.NewSecGroupName(), "tcp", "11.1.1.0/24", "80,443", "SG1")
					securityGroup1.Create()
					securityGroup2 = helpers.NewSecurityGroup(helpers.NewSecGroupName(), "tcp", "125.5.1.0/24", "25555", "SG2")
					securityGroup2.Create()
				})

				AfterEach(func() {
					securityGroup1.Delete()
					securityGroup2.Delete()
				})

				It("lists the security groups", func() {
					Eventually(session.Out).Should(Say("Getting security groups as admin"))
					Eventually(session.Out).Should(Say("\\s+Name\\s+Organization\\s+Space"))
					Eventually(session.Out).Should(Say("#\\d{1,2}%s\\s+%s\\s+%s", securityGroup1.Name, "", ""))
					Eventually(session.Out).Should(Say("#\\d{1,2}%s\\s+%s\\s+%s", securityGroup2.Name, "", ""))
					Eventually(session.Out).Should(Say("OK"))
					Eventually(session).Should(Exit(0))
				})
			})
		})
	})

	Describe("Refactored command", func() {
		var (
			session *Session
		)

		BeforeEach(func() {
			helpers.RunIfExperimental("skipping until approved")
		})

		JustBeforeEach(func() {
			session = helpers.CF("security-groups")
		})

		Context("when no API endpoint is set", func() {
			BeforeEach(func() {
				helpers.UnsetAPI()
			})

			It("fails with no API endpoint set message", func() {
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Err).Should(Say("No API endpoint set. Use 'cf login' or 'cf api' to target an endpoint."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when not logged in", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
			})

			It("fails with not logged in message", func() {
				session := helpers.CF("security-groups")
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Out).Should(Say("Not logged in\\. Use 'cf login' to log in\\."))
				Eventually(session).Should(Exit(1))
			})
		})

		Context("when too many arguments are provided", func() {
			BeforeEach(func() {
				helpers.LogoutCF()
			})

			It("fails with too many arguments message", func() {
				session := helpers.CF("security-groups", "foooo")
				Eventually(session.Out).Should(Say("FAILED"))
				Eventually(session.Out).Should(Say("Incorrect Usage\\. No argument required\\."))
				Eventually(session).Should(Exit(1))
			})
		})
	})
})
