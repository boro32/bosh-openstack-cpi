package volumetype_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVolumeTypeService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Volume Type Service Suite")
}
