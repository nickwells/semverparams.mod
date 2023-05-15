package semverparams_test

import (
	"errors"
	"testing"

	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v5/param"
	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/param.mod/v5/paramtest"
	"github.com/nickwells/semver.mod/v3/semver"
	"github.com/nickwells/semverparams.mod/v5/semverparams"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type semverPair struct {
	svv   *semverparams.SemverVals
	svCks *semverparams.SemverChecks

	addIDs bool
}

func makePSet(svp *semverPair) *param.PSet {
	ps := paramset.NewNoHelpNoExitNoErrRptOrPanic()
	if svp.svv != nil {
		_ = semverparams.AddSemverGroup(ps)
		_ = svp.svv.AddSemverParam(svp.svCks)(ps)
		if svp.addIDs {
			_ = svp.svv.AddIDParams(svp.svCks)(ps)
		}
	}
	if svp.svCks != nil {
		_ = svp.svCks.AddCheckParams()(ps)
	}

	return ps
}

// cmpSemverPairs compares the value with the expected value and returns
// an error if they differ
func cmpSemverPairs(iVal, iExpVal any) error {
	val, ok := iVal.(semverPair)
	if !ok {
		return errors.New("Bad value: not a semverPair struct")
	}
	expVal, ok := iExpVal.(semverPair)
	if !ok {
		return errors.New("Bad expected value: not a semverPair struct")
	}

	return testhelper.DiffVals(val, expVal, []string{"svCks"})
}

// mkTestParser populates and returns a paramtest.Parser ready to be added to
// the testcases.
func mkTestParser(
	errs errutil.ErrMap,
	id testhelper.ID,
	initVals semverPair,
	expVals semverPair,
	args ...string,
) paramtest.Parser {
	return paramtest.Parser{
		ID:             id,
		ExpParseErrors: errs,
		Val:            initVals,
		Ps:             makePSet(&initVals),
		ExpVal:         expVals,
		Args:           args,
		CheckFunc:      cmpSemverPairs,
	}
}

func TestParse(t *testing.T) {
	testCases := []paramtest.Parser{}

	{
		svvInit := semverparams.SemverVals{
			SemVer: &semver.SV{},
		}
		svvExp := semverparams.SemverVals{
			SemVer: semver.NewSVOrPanic(1, 2, 3, nil, nil),
		}

		testCases = append(testCases,
			mkTestParser(errutil.ErrMap{},
				testhelper.MkID("good semver, no prefix"),
				semverPair{svv: &svvInit},
				semverPair{svv: &svvExp},
				"-semver", "v1.2.3"))
	}
	{
		svvInit := semverparams.SemverVals{
			Prefix: "a",
			SemVer: &semver.SV{},
		}
		svvExp := semverparams.SemverVals{
			Prefix: "a",
			SemVer: semver.NewSVOrPanic(1, 2, 3, nil, nil),
		}

		testCases = append(testCases,
			mkTestParser(errutil.ErrMap{},
				testhelper.MkID("good semver, prefix: a"),
				semverPair{svv: &svvInit},
				semverPair{svv: &svvExp},
				"-a-semver", "v1.2.3"))
	}
	{
		svvInit := semverparams.SemverVals{
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			SemVer: semver.NewSVOrPanic(1, 2, 3,
				[]string{"rc", "1"}, nil),
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(errutil.ErrMap{},
				testhelper.MkID("good semver, with prID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-semver", "v1.2.3-rc.1",
				"-pre-rel-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}
	{
		svvInit := semverparams.SemverVals{
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			SemVer: semver.NewSVOrPanic(1, 2, 3,
				nil, []string{"rc", "1"}),
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(errutil.ErrMap{},
				testhelper.MkID("good semver, with buildID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-semver", "v1.2.3+rc.1",
				"-build-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}
	{
		parseErrs := errutil.ErrMap{}
		parseErrs.AddError(
			"Final Checks",
			errors.New("Bad PreRelIDs:"+
				" (the length of the list (3) is incorrect:"+
				" the value (3) must equal 0"+
				" or"+
				" the length of the list (3) is incorrect:"+
				" the value (3) must equal 2)"))

		svvInit := semverparams.SemverVals{
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			// The Final Checks don't prevent the value being set
			SemVer: semver.NewSVOrPanic(1, 2, 3,
				[]string{"rc", "1", "x"}, nil),
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(parseErrs,
				testhelper.MkID("bad semver, with prID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-semver", "v1.2.3-rc.1.x",
				"-pre-rel-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}
	{
		parseErrs := errutil.ErrMap{}
		parseErrs.AddError(
			"Final Checks",
			errors.New("Bad BuildIDs:"+
				" (the length of the list (3) is incorrect:"+
				" the value (3) must equal 0"+
				" or"+
				" the length of the list (3) is incorrect:"+
				" the value (3) must equal 2)"))

		svvInit := semverparams.SemverVals{
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			// The Final Checks don't prevent the value being set
			SemVer: semver.NewSVOrPanic(1, 2, 3,
				nil, []string{"rc", "1", "x"}),
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(parseErrs,
				testhelper.MkID("bad semver, with buildID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-semver", "v1.2.3+rc.1.x",
				"-build-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}
	{
		const desc = "test-desc"
		parseErrs := errutil.ErrMap{}
		parseErrs.AddError(
			"Final Checks",
			errors.New(desc+": Bad PreRelIDs:"+
				" (the length of the list (3) is incorrect:"+
				" the value (3) must equal 0"+
				" or"+
				" the length of the list (3) is incorrect:"+
				" the value (3) must equal 2)"))

		svvInit := semverparams.SemverVals{
			Prefix: "test-pfx",
			Desc:   desc,
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{Name: "test-name"}
		svvExp := semverparams.SemverVals{
			Prefix:    "test-pfx",
			Desc:      desc,
			SemVer:    &semver.SV{},
			PreRelIDs: []string{"rc", "1", "x"},
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(parseErrs,
				testhelper.MkID("bad pre-release IDs, with prID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-test-pfx-pre-rel-IDs", "rc.1.x",
				"-test-name-pre-rel-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}
	{
		const desc = "test-desc"
		parseErrs := errutil.ErrMap{}
		parseErrs.AddError(
			"Final Checks",
			errors.New(desc+": Bad BuildIDs:"+
				" (the length of the list (3) is incorrect:"+
				" the value (3) must equal 0"+
				" or"+
				" the length of the list (3) is incorrect:"+
				" the value (3) must equal 2)"))

		svvInit := semverparams.SemverVals{
			Desc:   desc,
			SemVer: &semver.SV{},
		}
		svCksInit := semverparams.SemverChecks{Name: "test-name"}
		svvExp := semverparams.SemverVals{
			Desc:     desc,
			SemVer:   &semver.SV{},
			BuildIDs: []string{"rc", "1", "x"},
		}
		svCksExp := semverparams.SemverChecks{}

		testCases = append(testCases,
			mkTestParser(parseErrs,
				testhelper.MkID("bad build IDs, with buildID checks"),
				semverPair{svv: &svvInit, svCks: &svCksInit, addIDs: true},
				semverPair{svv: &svvExp, svCks: &svCksExp, addIDs: true},
				"-build-IDs", "rc.1.x",
				"-test-name-build-ID-checks",
				`Or(Length(EQ(0)),`+
					`And(`+
					`  Length(EQ(2)),`+
					`  SliceByPos(`+
					`    EQ("rc"),`+
					`    MatchesPattern("[1-9][0-9]*", "numeric")`+
					`  )`+
					`))`))
	}

	for _, tc := range testCases {
		_ = tc.Test(t)
	}
}
