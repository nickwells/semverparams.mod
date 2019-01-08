package semverparams

import (
	"errors"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/param"
	"github.com/nickwells/param.mod/param/psetter"
	"github.com/nickwells/semver.mod/semver"
)

// SemVer is a semantic version number that will be set by the parameter
// parsing if it is passed to the program
var SemVer *semver.SV

// PreRelIDs is a list of Pre-Release IDs that will be set by the parameter
// parsing if the list is passed to the program
var PreRelIDs []string

// BuildIDs is a list of Build IDs that will be set by the parameter parsing
// if the list is passed to the program
var BuildIDs []string

var semverParam *param.ByName

// SetAttrOnSVStringParam allows you to set any desired attributes on the
// parameter which is used to supply the semantic version number. a typical
// use might be to set the param.MustBeSet attribute to ensure that there is
// a semantic version number to work with.
func SetAttrOnSVStringParam(attrs param.Attributes) error {
	if semverParam == nil {
		return errors.New("the semver parameter has not been created yet. " +
			" Call semverparams.AddParamVersion before setting the attributes")
	}

	return param.Attrs(attrs)(semverParam)
}

// AddSVStringParam will add parameters for setting the semantic version
// number to the passed ParamSet
func AddSVStringParam(ps *param.ParamSet) error {
	semverParam = ps.Add("semver", &SVSetter{Value: &SemVer},
		"specify the semantic version number to be used",
		param.AltName("vsn"))

	return nil
}

// AddIDParams will add parameters for setting the pre-release and build IDs
// of a semantic version number to the passed ParamSet
func AddIDParams(ps *param.ParamSet) error {
	ps.Add("pre-rel-IDs",
		psetter.StrListSetter{
			Value:            &PreRelIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.StringSlice{
				semver.CheckAllPreRelIDs,
				check.StringSliceLenGT(0),
			},
		},
		"specify a non-empty list of pre-release IDs"+
			" suitable for setting on a semantic version number",
		param.AltName("prIDs"),
	)

	ps.Add("build-IDs",
		psetter.StrListSetter{
			Value:            &BuildIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.StringSlice{
				semver.CheckAllBuildIDs,
				check.StringSliceLenGT(0),
			},
		},
		"specify a non-empty list of build IDs"+
			" suitable for setting on a semantic version number",
		param.AltName("bldIDs"),
	)

	return nil
}
