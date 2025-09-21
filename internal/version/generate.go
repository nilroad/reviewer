//go:generate sh -c "echo package version > build_version.go"
//go:generate sh -c "echo 'const CommitHash = \"'`git rev-parse HEAD`'\"' >> build_version.go"
//go:generate sh -c "echo 'const Version = \"'`git describe --tags --abbrev=0`'\"' >> build_version.go"
//go:generate sh -c "echo 'const ReleaseDate = \"'`date`'\"' >> build_version.go"

package version
