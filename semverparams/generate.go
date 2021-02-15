// +build generate

package semverparams

//go:generate mkparamfilefunc -private -group semver
//go:generate mkparamfilefunc -private -group semver-checks
