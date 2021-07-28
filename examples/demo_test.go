package examples

import (
	"github.com/bjartek/go-with-the-flow/gwtf"
	"github.com/onflow/cadence"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoWithTheFlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go With The Flow Suite")
}

var _ = Describe("Go With The Flow", func() {

	Context("Script", func() {
			g := gwtf.NewTestingEmulator()
			It("Should run inline script", func() {

				value, err := g.Script(`
pub fun main(account: Address): String {
    return getAccount(account).address.toString()
}`).AccountArgument("second").RunReturns()

				Expect(err).Should(BeNil())

				Expect(value).Should(BeEquivalentTo(cadence.String("0x179b6b1cb6755e31")))
			})

			It("Should run script from file", func() {
				value, err :=g.ScriptFromFile("test").AccountArgument("second").RunReturns()
				Expect(err).Should(BeNil())
				Expect(value).Should(BeEquivalentTo(cadence.String("0x179b6b1cb6755e31")))
			})

			It("Should report error on wrong script path", func() {

				value, err := g.ScriptFromFile("foo").RunReturns()
				Expect(value).Should(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("./scripts/foo.cdc: no such file"))
			})
	})

	Context("Transaction", func () {
		g := gwtf.NewTestingEmulator()

		It("Should run inline transaction", func() {
			value, err :=g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(test:String) {
  prepare(acct: AuthAccount) {
	Debug.log(test)
 }
}`).SignProposeAndPayAs("first").StringArgument("foobar").RunFormatEvents()

		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(value)).Should(BeEquivalentTo(1))
		Expect(value[0].Name).Should(ContainSubstring("Debug.Log"))
		Expect(value[0].Fields["msg"]).Should(BeEquivalentTo("foobar"))

		})

		It("Should sign with two auth accounts", func() {

			g.TransactionFromFile("signWithMultipleAccounts").SignProposeAndPayAs("first").PayloadSigner("second").Run()

		})


	})

})
