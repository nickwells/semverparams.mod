package semverparams

import (
	"fmt"
	"sync"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v4/checksetter"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/semver.mod/v2/semver"
)

// SemverVals holds the semantic version ID values - the version ID and any
// pre-release or build IDs
type SemverVals struct {
	// Prefix is the optional prefix to apply to the parameter names. If it
	// is not empty it will be separated from the rest of the parameter name
	// with '-'.
	Prefix string

	// Desc is optional text to appear at the start of any errors reported by
	// the checks on build IDs or pre-release IDs. If it is not empty it will
	// be separated from the rest of the error message with ': '
	Desc string

	// SemVer is a semantic version number that will be set by the parameter
	// parsing if it is passed to the program
	SemVer *semver.SV

	// PreRelIDs is a list of Pre-Release IDs that will be set by the parameter
	// parsing if the list is passed to the program
	PreRelIDs []string

	// BuildIDs is a list of Build IDs that will be set by the parameter parsing
	// if the list is passed to the program
	BuildIDs []string
}

// SemverChecks holds the checks to be applied to the pre-release and build
// IDs
type SemverChecks struct {
	// Prefix is the optional prefix to apply to the parameter names. If it
	// is not empty it will be separated from the rest of the parameter name
	// with '-'.
	Prefix string

	// PreRelIDChecks is a list of checks to be applied to the pre-release IDs
	PreRelIDChecks []check.ValCk[[]string]

	// BuildIDChecks is a list of checks to be applied to the build IDs
	BuildIDChecks []check.ValCk[[]string]
}

const (
	semverGroupName       = "semver"
	semverChecksGroupName = "semver-checks"
)

var setGroupOnce sync.Once

// addSVGroup will add the semver group to the set of parameter groups
func addSVGroup(ps *param.PSet) {
	setGroupOnce.Do(func() {
		ps.AddGroup(semverGroupName,
			"common parameters concerned with "+semver.Names)
		_ = setGlobalConfigFileForGroupSemver(ps)
		_ = setConfigFileForGroupSemver(ps)
	})
}

// AddSVStringParam will add a parameter for setting the semantic version
// number to the passed PSet
func (svv *SemverVals) AddSVStringParam(
	ps *param.PSet, attrs param.Attributes) error {
	prefix := ""
	if svv.Prefix != "" {
		prefix = svv.Prefix + "-"
	}

	addSVGroup(ps)

	ps.Add(prefix+"semver", &SVSetter{Value: &svv.SemVer},
		"specify the "+semver.Name+" to be used",
		param.AltNames(prefix+"svn"),
		param.GroupName(semverGroupName),
		param.Attrs(attrs),
	)

	return nil
}

// AddIDParams will add parameters for setting the pre-release and build IDs
// of a semantic version number to the passed PSet
func (svv *SemverVals) AddIDParams(ps *param.PSet, svCks *SemverChecks) error {
	prefix := ""
	if svv.Prefix != "" {
		prefix = svv.Prefix + "-"
	}

	addSVGroup(ps)
	ps.Add(prefix+"pre-rel-IDs",
		psetter.StrList{
			Value:            &svv.PreRelIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.ValCk[[]string]{
				check.SliceAll[[]string](semver.CheckPreRelID),
				check.SliceLength[[]string](check.ValGT(0)),
			},
		},
		"specify a non-empty list of pre-release IDs"+
			" suitable for setting on a "+semver.Name,
		param.AltNames(prefix+"prIDs"),
		param.GroupName(semverGroupName),
	)

	ps.Add(prefix+"build-IDs",
		psetter.StrList{
			Value:            &svv.BuildIDs,
			StrListSeparator: psetter.StrListSeparator{Sep: "."},
			Checks: []check.ValCk[[]string]{
				check.SliceAll[[]string](semver.CheckBuildID),
				check.SliceLength[[]string](check.ValGT(0)),
			},
		},
		"specify a non-empty list of build IDs"+
			" suitable for setting on a "+semver.Name,
		param.AltNames(prefix+"bldIDs"),
		param.GroupName(semverGroupName),
	)

	if svCks != nil {
		ps.AddFinalCheck(checkIDs(svv, svCks))
	}

	return nil
}

// AddIDCheckerParams will add parameters for setting the checks to be
// applied to any pre-release and build IDs of a semantic version number.
func (svCks *SemverChecks) AddIDCheckerParams(ps *param.PSet) error {
	prefix := ""
	if svCks.Prefix != "" {
		prefix = svCks.Prefix + "-"
	}

	ps.AddGroup(semverChecksGroupName,
		"common parameters for specifying checks on "+
			semver.Name+
			" (pre-release and build IDs)")
	_ = setGlobalConfigFileForGroupSemverChecks(ps)
	_ = setConfigFileForGroupSemverChecks(ps)

	const paramDescIntro = "give a non-empty list of check functions to apply"

	ps.Add(prefix+"pre-rel-ID-checks",
		&checksetter.Setter[[]string]{
			Value: &svCks.PreRelIDChecks,
			Parser: checksetter.FindParserOrPanic[[]string](
				checksetter.StringSliceCheckerName),
		},
		paramDescIntro+" to the pre-release IDs for the "+semver.Name,
		param.AltNames(prefix+"prID-checks"),
		param.GroupName(semverChecksGroupName),
	)

	ps.Add(prefix+"build-ID-checks",
		&checksetter.Setter[[]string]{
			Value: &svCks.BuildIDChecks,
			Parser: checksetter.FindParserOrPanic[[]string](
				checksetter.StringSliceCheckerName),
		},
		paramDescIntro+" to the build IDs for the "+semver.Name,
		param.AltNames(prefix+"bldID-checks"),
		param.GroupName(semverChecksGroupName),
	)

	return nil
}

// checkIDs checks that any supplied IDs conform to the specified checks
func checkIDs(svv *SemverVals, svCks *SemverChecks) func() error {
	return func() error {
		errPfx := ""
		if svv.Desc != "" {
			errPfx = svv.Desc + ": "
		}
		for _, chk := range svCks.PreRelIDChecks {
			err := chk(svv.PreRelIDs)
			if err != nil {
				return fmt.Errorf("%sBad PreRelIDs: %s", errPfx, err)
			}
		}

		for _, chk := range svCks.BuildIDChecks {
			err := chk(svv.BuildIDs)
			if err != nil {
				return fmt.Errorf("%sBad BuildIDs: %s", errPfx, err)
			}
		}

		return nil
	}
}
