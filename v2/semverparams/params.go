package semverparams

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
	"github.com/nickwells/filecheck.mod/filecheck"
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
	"github.com/nickwells/semver.mod/semver"
	"github.com/nickwells/xdg.mod/xdg"
)

// SemVer is a semantic version number that will be set by the parameter
// parsing if it is passed to the program
var SemVer *semver.SV

// PreRelIDs is a list of Pre-Release IDs that will be set by the parameter
// parsing if the list is passed to the program
var PreRelIDs []string

// PreRelIDChecks is a list of checks to be applied to the pre-release IDs
var PreRelIDChecks []check.StringSlice

// BuildIDs is a list of Build IDs that will be set by the parameter parsing
// if the list is passed to the program
var BuildIDs []string

// BuildIDChecks is a list of checks to be applied to the build IDs
var BuildIDChecks []check.StringSlice

var semverParam *param.ByName

const semverGroupName = "semver"

var setGroupOnce sync.Once

// addSVGroup will add the semver group to the set of parameter groups
func addSVGroup(ps *param.PSet) {
	setGroupOnce.Do(func() {
		ps.AddGroup(semverGroupName,
			"common parameters concerned with "+semver.Names)
		ps.AddGroupConfigFile(semverGroupName,
			filepath.Join(xdg.ConfigHome(), "semver.config"),
			filecheck.Optional)
	})
}

// SetAttrOnSVStringParam allows you to set any desired attributes on the
// parameter which is used to supply the semantic version number. a typical
// use might be to set the param.MustBeSet attribute to ensure that there is
// a semantic version number to work with.
func SetAttrOnSVStringParam(attrs param.Attributes) error {
	if semverParam == nil {
		return errors.New("the semver parameter has not been created yet. " +
			" Call semverparams.AddSVStringParam before setting the attributes")
	}

	return param.Attrs(attrs)(semverParam)
}

// AddSVStringParam will add parameters for setting the semantic version
// number to the passed PSet
func AddSVStringParam(ps *param.PSet) error {
	addSVGroup(ps)
	semverParam = ps.Add("semver", &SVSetter{Value: &SemVer},
		"specify the "+semver.Name+" to be used",
		param.AltName("svn"),
		param.GroupName(semverGroupName),
	)

	return nil
}

// AddIDParams will add parameters for setting the pre-release and build IDs
// of a semantic version number to the passed PSet
func AddIDParams(ps *param.PSet) error {
	addSVGroup(ps)
	ps.Add("pre-rel-IDs",
		psetter.StrList{
			Value:            &PreRelIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.StringSlice{
				semver.CheckAllPreRelIDs,
				check.StringSliceLenGT(0),
			},
		},
		"specify a non-empty list of pre-release IDs"+
			" suitable for setting on a "+semver.Name,
		param.AltName("prIDs"),
		param.GroupName(semverGroupName),
	)

	ps.Add("build-IDs",
		psetter.StrList{
			Value:            &BuildIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.StringSlice{
				semver.CheckAllBuildIDs,
				check.StringSliceLenGT(0),
			},
		},
		"specify a non-empty list of build IDs"+
			" suitable for setting on a "+semver.Name,
		param.AltName("bldIDs"),
		param.GroupName(semverGroupName),
	)

	return nil
}

// AddIDCheckerParams will add parameters for setting the checks to be
// applied to any pre-release and build IDs of a semantic version number.
func AddIDCheckerParams(ps *param.PSet) error {
	addSVGroup(ps)
	const paramDescIntro = "specify a non-empty list of check functions to apply"

	ps.Add("pre-rel-ID-checks",
		checksetter.StringSlice{
			Value: &PreRelIDChecks,
		},
		paramDescIntro+" to the pre-release IDs for the "+semver.Name,
		param.AltName("prID-checks"),
		param.GroupName(semverGroupName),
	)

	ps.Add("build-ID-checks",
		checksetter.StringSlice{
			Value: &BuildIDChecks,
		},
		paramDescIntro+" to the build IDs for the "+semver.Name,
		param.AltName("bldID-checks"),
		param.GroupName(semverGroupName),
	)

	return nil
}
