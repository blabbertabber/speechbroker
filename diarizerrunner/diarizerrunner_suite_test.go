package diarizerrunner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDiarizerrunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diarizerrunner Suite")
}
