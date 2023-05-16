package semverparams

import (
	"errors"
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v4/checksetter"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/psetter"
	"github.com/nickwells/semver.mod/v3/semver"
)

// SemverVals holds the semantic version ID values - the version ID and any
// pre-release or build IDs. If you want to have multiple SemverVals one can
// have an empty Prefix but each of the rest will need to have its own
// distinct Prefix. When you add the parameters for these they will all
// appear in the same parameter group.
type SemverVals struct {
	// Prefix is the optional prefix to apply to the parameter names. If it
	// is not empty it will be separated from the rest of the parameter name
	// with '-'.
	//
	// Note that it must be suitable to be part of a parameter name (it must
	// start with a letter and be followed with letters, digits or dashes
	// '-')
	Prefix string

	// Desc is optional text to appear at the start of any errors reported by
	// the checks on build IDs or pre-release IDs. If it is not empty it will
	// be separated from the rest of the error message with ': '
	Desc string

	// SemVer is a semantic version number that will be set by the parameter
	// parsing if it is passed to the program
	SemVer semver.SV

	// SemverAttrs gives the attributes to be applied to the parameter for
	// setting the SemVer
	SemverAttrs param.Attributes

	// PreRelIDs is a list of Pre-Release IDs that will be set by the parameter
	// parsing if the list is passed to the program
	PreRelIDs []string

	// PreRelIDAttrs gives the attributes to be applied to the parameter for
	// setting the pre-release IDs
	PreRelIDAttrs param.Attributes

	// BuildIDs is a list of Build IDs that will be set by the parameter parsing
	// if the list is passed to the program
	BuildIDs []string

	// BuildIDAttrs gives the attributes to be applied to the parameter for
	// setting the build IDs
	BuildIDAttrs param.Attributes
}

// SemverChecks holds the checks to be applied to the pre-release and build
// IDs. If you want to have multiple SemverChecks each will need its own
// distinct Name. Each set of parameters will appear in their own parameter
// group, only the group with an empty Name will have an associated group
// config file.
type SemverChecks struct {
	// Name, if not empty, will be applied as a prefix to the parameter
	// names, separated from the rest of the parameter name with '-'. The
	// Name is also used as a suffix to the group name so that all the
	// parameters with the same prefix appear in a group of their own with a
	// name reflecting their prefix.
	//
	// Note that it must be suitable to be part of a parameter name (it must
	// start with a letter and be followed with letters, digits or dashes
	// '-')
	//
	// Only the default group (with Name == "") has an associated group
	// config file
	Name string

	// Desc will be added to the description of the parameter group and to
	// each of the check-setting parameters' help text
	Desc string

	// PreRelIDChecks is a list of checks to be applied to the pre-release IDs
	PreRelIDChecks []check.ValCk[[]string]

	// BuildIDChecks is a list of checks to be applied to the build IDs
	BuildIDChecks []check.ValCk[[]string]
}

const (
	semverGroupName       = "semver"
	semverChecksGroupName = "semver-checks"
)

func AddSemverGroup(ps *param.PSet) error {
	ps.AddGroup(semverGroupName,
		"common parameters concerned with "+semver.Names)
	return nil
}

// AddSemverParam returns a function that will add a parameter for setting
// the semantic version number to the passed PSet
func (svv *SemverVals) AddSemverParam(svCks *SemverChecks) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		prefix := ""
		if svv.Prefix != "" {
			prefix = svv.Prefix + "-"
		}

		ps.Add(prefix+"semver", SVSetter{Value: &svv.SemVer},
			"specify the "+semver.Name+" to be used",
			param.AltNames(prefix+"svn"),
			param.GroupName(semverGroupName),
			param.Attrs(svv.SemverAttrs),
		)

		if svCks != nil {
			ps.AddFinalCheck(
				checkSemverIDs(svv, svCks))
		}

		return nil
	}
}

// IDListSetter will return a psetter.StrList correctly constructed for
// setting a list of semver IDs (either pre-release or build IDs). You should
// pass the appropriate semver.Check...ID function depending on the type of
// list of IDs you want to set. It will panic (as this is a coding error) if
// the idChk function is nil.
func IDListSetter(val *[]string, idChk check.ValCk[string]) psetter.StrList {
	if idChk == nil {
		panic(errors.New(
			"the function to check the parts of the ID list is nil"))
	}

	return psetter.StrList{
		Value:            val,
		StrListSeparator: psetter.StrListSeparator{Sep: "."},
		Checks: []check.ValCk[[]string]{
			check.SliceAll[[]string](idChk),
			check.SliceLength[[]string](check.ValGT(0)),
		},
	}
}

// AddIDParams returns a function that will add parameters for setting the
// pre-release and build IDs of a semantic version number to the passed
// PSet. If a non-nil SemverChecks is passed then a final check is added of
// the pre-release and build IDs against the checks, if any, given by the
// SemverChecks
func (svv *SemverVals) AddIDParams(svCks *SemverChecks) param.PSetOptFunc {
	return func(ps *param.PSet) error {
		prefix := ""
		if svv.Prefix != "" {
			prefix = svv.Prefix + "-"
		}

		ps.Add(prefix+"pre-rel-IDs",
			IDListSetter(&svv.PreRelIDs, semver.CheckPreRelID),
			"specify a non-empty list of pre-release IDs"+
				" suitable for setting on a "+semver.Name,
			param.AltNames(prefix+"prIDs"),
			param.GroupName(semverGroupName),
			param.Attrs(svv.PreRelIDAttrs),
		)

		ps.Add(prefix+"build-IDs",
			IDListSetter(&svv.BuildIDs, semver.CheckBuildID),
			"specify a non-empty list of build IDs"+
				" suitable for setting on a "+semver.Name,
			param.AltNames(prefix+"bldIDs"),
			param.GroupName(semverGroupName),
			param.Attrs(svv.BuildIDAttrs),
		)

		if svCks != nil {
			ps.AddFinalCheck(
				checkIDs(svv, svCks))
		}

		return nil
	}
}

// checkIDs checks that any supplied IDs conform to the specified checks
func checkIDs(svv *SemverVals, svCks *SemverChecks) param.FinalCheckFunc {
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

// checkSemverIDs checks that the pre-release and build IDs on the semver
// conform to the specified checks
func checkSemverIDs(svv *SemverVals, svCks *SemverChecks) param.FinalCheckFunc {
	return func() error {
		errPfx := ""
		if svv.Desc != "" {
			errPfx = svv.Desc + ": "
		}
		for _, chk := range svCks.PreRelIDChecks {
			err := chk(svv.SemVer.PreRelIDs())
			if err != nil {
				return fmt.Errorf("%sBad PreRelIDs: %s", errPfx, err)
			}
		}

		for _, chk := range svCks.BuildIDChecks {
			err := chk(svv.SemVer.BuildIDs())
			if err != nil {
				return fmt.Errorf("%sBad BuildIDs: %s", errPfx, err)
			}
		}

		return nil
	}
}

// AddCheckParams will add parameters for setting the checks to be
// applied to any pre-release and build IDs of a semantic version number.
func (svCks *SemverChecks) AddCheckParams() param.PSetOptFunc {
	return func(ps *param.PSet) error {
		prefix := ""
		groupName := semverChecksGroupName
		if svCks.Name != "" {
			prefix = svCks.Name + "-"
			groupName += "-" + svCks.Name
		}

		ps.AddGroup(groupName,
			"common parameters for specifying checks on "+
				semver.Name+
				" (pre-release and build IDs)"+svCks.Desc)

		if svCks.Name == "" {
			_ = setGlobalConfigFileForGroupSemverChecks(ps)
			_ = setConfigFileForGroupSemverChecks(ps)
		}

		helpText := func(part string) string {
			return fmt.Sprintf(
				"give a non-empty list of check functions"+
					" to apply to the %s for the %s%s",
				part, semver.Name, svCks.Desc)
		}

		ps.Add(prefix+"pre-rel-ID-checks",
			&checksetter.Setter[[]string]{
				Value: &svCks.PreRelIDChecks,
				Parser: checksetter.FindParserOrPanic[[]string](
					checksetter.StringSliceCheckerName),
			},
			helpText("pre-release IDs"),
			param.AltNames(prefix+"prID-checks"),
			param.GroupName(groupName),
		)

		ps.Add(prefix+"build-ID-checks",
			&checksetter.Setter[[]string]{
				Value: &svCks.BuildIDChecks,
				Parser: checksetter.FindParserOrPanic[[]string](
					checksetter.StringSliceCheckerName),
			},
			helpText("build IDs"),
			param.AltNames(prefix+"bldID-checks"),
			param.GroupName(groupName),
		)

		return nil
	}
}
