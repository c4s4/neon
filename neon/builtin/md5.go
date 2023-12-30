package builtin

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/c4s4/neon/neon/build"
)

func init() {
	build.AddBuiltin(build.BuiltinDesc{
		Name: "md5",
		Func: md5Sum,
		Help: `Return MD5 sum of given file.

Arguments:

- The file name to get MD5 sum for.

Returns:

- The MD5 sum of given file.

Examples:

    # get MD5 sum of file README
    md5("README")
    # returns: MD5 sum of file "README"`,
	})
}

func md5Sum(file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}
