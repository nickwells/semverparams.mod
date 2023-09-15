package semverparams_test

import (
	"errors"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/errutil.mod/errutil"
	"github.com/nickwells/param.mod/v6/param"
	"github.com/nickwells/param.mod/v6/paramset"
	"github.com/nickwells/param.mod/v6/paramtest"
	"github.com/nickwells/semver.mod/v3/semver"
	"github.com/nickwells/semverparams.mod/v6/semverparams"
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

	return testhelper.DiffVals(val, expVal,
		[]string{"svCks"},
		[]string{"svv", "semverParam"},
		[]string{"svv", "preRelIDsParam"},
		[]string{"svv", "buildIDsParam"})
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
		svvInit := semverparams.SemverVals{}
		svvExp := semverparams.SemverVals{
			SemVer: *semver.NewSVOrPanic(1, 2, 3, nil, nil),
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
		}
		svvExp := semverparams.SemverVals{
			Prefix: "a",
			SemVer: *semver.NewSVOrPanic(1, 2, 3, nil, nil),
		}

		testCases = append(testCases,
			mkTestParser(errutil.ErrMap{},
				testhelper.MkID("good semver, prefix: a"),
				semverPair{svv: &svvInit},
				semverPair{svv: &svvExp},
				"-a-semver", "v1.2.3"))
	}
	{
		svvInit := semverparams.SemverVals{}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			SemVer: *semver.NewSVOrPanic(1, 2, 3,
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
		svvInit := semverparams.SemverVals{}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			SemVer: *semver.NewSVOrPanic(1, 2, 3,
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
				" either [the length of the list (3) is incorrect:"+
				" the value (3) must equal 0]"+
				" or"+
				" [the length of the list (3) is incorrect:"+
				" the value (3) must equal 2]"))

		svvInit := semverparams.SemverVals{}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			// The Final Checks don't prevent the value being set
			SemVer: *semver.NewSVOrPanic(1, 2, 3,
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
				" either [the length of the list (3) is incorrect:"+
				" the value (3) must equal 0]"+
				" or"+
				" [the length of the list (3) is incorrect:"+
				" the value (3) must equal 2]"))

		svvInit := semverparams.SemverVals{}
		svCksInit := semverparams.SemverChecks{}
		svvExp := semverparams.SemverVals{
			// The Final Checks don't prevent the value being set
			SemVer: *semver.NewSVOrPanic(1, 2, 3,
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
				" either [the length of the list (3) is incorrect:"+
				" the value (3) must equal 0]"+
				" or"+
				" [the length of the list (3) is incorrect:"+
				" the value (3) must equal 2]"))

		svvInit := semverparams.SemverVals{
			Prefix: "test-pfx",
			Desc:   desc,
		}
		svCksInit := semverparams.SemverChecks{Name: "test-name"}
		svvExp := semverparams.SemverVals{
			Prefix:    "test-pfx",
			Desc:      desc,
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
				" either [the length of the list (3) is incorrect:"+
				" the value (3) must equal 0]"+
				" or"+
				" [the length of the list (3) is incorrect:"+
				" the value (3) must equal 2]"))

		svvInit := semverparams.SemverVals{
			Desc: desc,
		}
		svCksInit := semverparams.SemverChecks{Name: "test-name"}
		svvExp := semverparams.SemverVals{
			Desc:     desc,
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

func TestIDListSetter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		testhelper.ExpErr
		idChk  check.ValCk[string]
		args   []string
		expVal []string
	}{
		{
			ID:    testhelper.MkID("good - pre-rel-ids, no args"),
			idChk: semver.CheckPreRelID,
		},
		{
			ID:     testhelper.MkID("good - pre-rel-ids, with args"),
			idChk:  semver.CheckPreRelID,
			args:   []string{"-set", "a.b.42"},
			expVal: []string{"a", "b", "42"},
		},
		{
			ID: testhelper.MkID("good - pre-rel-ids, with bad args"),
			ExpErr: testhelper.MkExpErr(
				"list entry: 0 (bad,val) does not pass the test:" +
					" the Pre-Rel ID: 'bad,val' must be" +
					" a non-empty string of letters, digits or hyphens"),
			idChk: semver.CheckPreRelID,
			args:  []string{"-set", "bad,val"},
		},
		{
			ID: testhelper.MkID("good - pre-rel-ids, with bad args: leading 0"),
			ExpErr: testhelper.MkExpErr(
				"list entry: 0 (01) does not pass the test:" +
					" the Pre-Rel ID: '01' must" +
					" have no leading zero if it's all numeric"),
			idChk: semver.CheckPreRelID,
			args:  []string{"-set", "01"},
		},
		{
			ID:    testhelper.MkID("good - build-ids, no args"),
			idChk: semver.CheckBuildID,
		},
		{
			ID:     testhelper.MkID("good - build-ids, with args"),
			idChk:  semver.CheckBuildID,
			args:   []string{"-set", "a.b.42.01"},
			expVal: []string{"a", "b", "42", "01"},
		},
		{
			ID: testhelper.MkID("good - build-ids, with bad args"),
			ExpErr: testhelper.MkExpErr(
				"list entry: 0 (bad,val) does not pass the test:" +
					" the Build ID: 'bad,val' must be" +
					" a non-empty string of letters, digits or hyphens"),
			idChk: semver.CheckBuildID,
			args:  []string{"-set", "bad,val"},
		},
		{
			ID: testhelper.MkID("bad - nil check"),
			ExpPanic: testhelper.MkExpPanic(
				"the function to check the parts of the ID list is nil"),
		},
	}

	for _, tc := range testCases {
		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatalf("Unexpected error creating the paramset: %v", err)
		}
		val := []string{}
		panicked, panicVal := testhelper.PanicSafe(func() {
			ps.Add("set", semverparams.IDListSetter(&val, tc.idChk), "help")
		})
		if !testhelper.CheckExpPanicError(t, panicked, panicVal, tc) &&
			!panicked {
			errMap := ps.Parse(tc.args)
			if len(errMap) != 0 {
				for _, v := range errMap {
					if len(v) > 0 {
						err = v[0]
					}
				}
			}
			if testhelper.CheckExpErr(t, err, tc) && err == nil {
				testhelper.DiffStringSlice(t,
					tc.IDStr(), "value",
					val, tc.expVal)
			}
		}
	}
}
