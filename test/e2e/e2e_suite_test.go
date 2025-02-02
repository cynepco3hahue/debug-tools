package e2eknit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jaypipes/ghw/pkg/snapshot"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	envVarKniSnapshotKeep string = "KNI_SNAPSHOT_KEEP"
)

var (
	knitBaseDir  string
	binariesPath string
	snapshotKeep bool
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "E2E Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		ginkgo.Fail("Cannot retrieve tests directory")
	}
	basedir := filepath.Dir(file)
	knitBaseDir = filepath.Clean(filepath.Join(basedir, "..", ".."))
	binariesPath = filepath.Clean(filepath.Join(knitBaseDir, "_output"))
	fmt.Fprintf(ginkgo.GinkgoWriter, "using binaries at %q\n", binariesPath)

	if _, ok = os.LookupEnv(envVarKniSnapshotKeep); ok {
		snapshotKeep = true
	}
})

func snapshotBeforeEach(fixtureName, snapshotName string) string {
	path := filepath.Join(dataDirFor(fixtureName), snapshotName)

	unpackedPath, err := snapshot.Unpack(path)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Failed to unpack the snapshot %q: %v", path, err))
	}

	fmt.Fprintf(ginkgo.GinkgoWriter, "unpacked snapshot %q at %q\n", path, unpackedPath)
	return unpackedPath
}

func snapshotAfterEach(snapshotRoot string) {
	if snapshotKeep {
		return
	}
	if err := snapshot.Cleanup(snapshotRoot); err != nil {
		ginkgo.Fail(fmt.Sprintf("Failed to cleanup the snapshot at %q: %v", snapshotRoot, err))
	}
}

func areJSONBlobsEqual(b1, b2 []byte) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	if err := json.Unmarshal(b1, &o1); err != nil {
		return false, fmt.Errorf("Error unmarshalling string 1: %v", err)
	}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false, fmt.Errorf("Error unmarshalling string 2: %v", err)
	}

	return reflect.DeepEqual(o1, o2), nil
}

func dataDirFor(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		ginkgo.Fail("Cannot retrieve tests directory")
	}
	basedir := filepath.Dir(file)
	return filepath.Clean(filepath.Join(basedir, "..", "data", name))

}
