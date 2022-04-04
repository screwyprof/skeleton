package version

import "fmt"

var (
	// AppName provides application name used for metrics/tracing/etc... Set by linker.
	AppName = "skeleton"
	// AppVersion provides version of the application. Typically name of git branch. Set by linker.
	AppVersion = "develop"
	// GoVersion provides version of go the binary was compiled with. Set by linker.
	GoVersion string
	// BuildDate contains date of the build. Set by linker.
	BuildDate string
	// GitRev provides exact git revision of the source the binary was built from. Set by linker.
	GitRev string
	// GitLog provides exact git log. Set by linker.
	GitLog string
)

func PrintBuildInfo() {
	fmt.Printf("App Name: %s\n", AppName)
	fmt.Printf("App Version: %q, Rev: [%s]\n", AppVersion, GitLog)
	fmt.Printf("Built at %s with %q\n", BuildDate, "Go "+GoVersion)
}
