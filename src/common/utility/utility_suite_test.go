package utility_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUtility(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utility Suite")
}
