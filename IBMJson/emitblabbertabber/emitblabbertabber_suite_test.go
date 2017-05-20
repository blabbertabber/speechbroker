package emitblabbertabber_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEmitblabbertabber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Emitblabbertabber Suite")
}
