package semverparams

// Code generated by mkparamfilefunc DO NOT EDIT.
// with parameters set at:
//	command line:3: -group semver
//	command line:1: -private
import (
	"path/filepath"

	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/xdg.mod/xdg"
)

/*
setConfigFileForGroupSemver adds a config file to the set which the param parser
will process before checking the command line parameters.

This function is one of a pair which add the global and personal config files.
It is generally best practice to add the global config file before adding the
personal one. This allows any system-wide defaults to be overridden by personal
choices. Also any parameters which can only be set once can be set in the global
config file, thereby enforcing a global policy.
*/
func setConfigFileForGroupSemver(ps *param.PSet) error {
	baseDir := xdg.ConfigHome()

	ps.AddGroupConfigFile("semver",
		filepath.Join(baseDir,
			"github.com",
			"nickwells",
			"semverparams.mod",
			"v4",
			"semverparams",
			"group-semver.cfg"),
		filecheck.Optional)
	return nil
}

/*
setGlobalConfigFileForGroupSemver adds a config file to the set which the param
parser will process before checking the command line parameters.

This function is one of a pair which add the global and personal config files.
It is generally best practice to add the global config file before adding the
personal one. This allows any system-wide defaults to be overridden by personal
choices. Also any parameters which can only be set once can be set in the global
config file, thereby enforcing a global policy.
*/
func setGlobalConfigFileForGroupSemver(ps *param.PSet) error {
	dirs := xdg.ConfigDirs()
	if len(dirs) == 0 {
		return nil
	}
	baseDir := dirs[0]

	ps.AddGroupConfigFile("semver",
		filepath.Join(baseDir,
			"github.com",
			"nickwells",
			"semverparams.mod",
			"v4",
			"semverparams",
			"group-semver.cfg"),
		filecheck.Optional)
	return nil
}
