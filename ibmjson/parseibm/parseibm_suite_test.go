package parseibm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestParseibm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parseibm Suite")
}
