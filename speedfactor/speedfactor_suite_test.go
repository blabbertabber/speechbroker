package speedfactor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSpeedfactor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Speedfactor Suite")
}