package ibmservicecreds_test

import (
	. "github.com/blabbertabber/speechbroker/ibmservicecreds"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strings"
)

var _ = Describe("IBMServiceCreds", func() {
	Context(".ReadCredsFromReader", func() {
		Context("When a reader is passed an empty JSON", func() {
			It("returns an empty struct", func() {
				sourceReader := strings.NewReader("{}")
				expectation := IBMServiceCreds{}
				readCreds, err := ReadCredsFromReader(sourceReader)
				Expect(err).To(BeNil())
				Expect(readCreds).To(Equal(expectation))
			})
		})
		Context("When a reader is passed a populated JSON", func() {
			It("returns a populated struct", func() {
				source, err := ioutil.ReadFile("../assets/test/ibm_service_creds.json")
				Expect(err).To(BeNil())
				sourceReader := strings.NewReader(string(source))
				expectation := IBMServiceCreds{
					Url:      "https://stream.watsonplatform.net/speech-to-text/api",
					Username: "9f6cdead-d9d3-49db-96e4-deadbeefdead",
					Password: "8rgJeCunnie8",
				}
				readCreds, err := ReadCredsFromReader(sourceReader)
				Expect(err).To(BeNil())
				Expect(readCreds).To(Equal(expectation))
			})
		})
	})
	Context(".ReadCredsFromPath", func() {
		Context("When path is exists and is valid JSON", func() {
			It("returns the expected struct", func() {
				readCreds, err := ReadCredsFromPath("../assets/test/ibm_service_creds.json")
				Expect(err).To(BeNil())
				expectation := IBMServiceCreds{
					Url:      "https://stream.watsonplatform.net/speech-to-text/api",
					Username: "9f6cdead-d9d3-49db-96e4-deadbeefdead",
					Password: "8rgJeCunnie8",
				}
				Expect(readCreds).To(Equal(expectation))
			})
		})
		Context("When path is non-existent", func() {
			It("panics", func() {
				Expect(func() {
					ReadCredsFromPath("/non/existent/path")
				}).Should(Panic())
			})
		})
	})

})
