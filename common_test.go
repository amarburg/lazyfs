package lazyfs

//import "net/url"
import "reflect"
import "github.com/amarburg/go-lazyfs-testfiles"

var BufSize = 10

var test_pairs = []struct {
	offset int64
	length int
}{
	{0, BufSize},
	{10, BufSize},
	{20, 6},
}

var AlphabetSize int64 = 26

var TestCacheDir = "test_cache/"

var LocalFilesRoot = lazyfs_testfiles.RepoRoot() + "/"

var LocalStoreRoot = TestCacheDir + "localsource_localstore/"
var SparseStoreRoot = TestCacheDir + "localsource_sparsestore/"

var HttpSourceSparseStore = TestCacheDir + "httpsource_sparsestore/"

var BadPath = TestCacheDir + "a/y/foo.fs"

var OOIRawDataRootURL = "https://rawdata.oceanobservatories.org/files/"

func CheckTestFile(buf []byte, off int64) bool {
	// I'm sure there's a beautiful idiomatic Go way to do this
	testString := [26]byte{65, 66, 67, 68, 69, 70, 71, 72, 73, 74,
		75, 76, 77, 78, 79, 80, 81, 82, 83, 84,
		85, 86, 87, 88, 89, 90}

	l := int(off) + len(buf)
	return reflect.DeepEqual(buf, testString[int(off):l])
}
