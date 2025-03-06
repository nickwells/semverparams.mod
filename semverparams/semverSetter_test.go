package semverparams_test

import (
	"testing"

	"github.com/nickwells/semver.mod/v3/semver"
	"github.com/nickwells/semverparams.mod/v6/semverparams"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSetter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		paramVal string
		expVal   string
	}{
		{
			ID:       testhelper.MkID("good"),
			paramVal: "v1.2.3",
			expVal:   "v1.2.3",
		},
		{
			ID: testhelper.MkID("bad"),
			ExpErr: testhelper.MkExpErr(
				"bad semantic version ID - it does not start with a 'v'"),
			paramVal: "x1.2.3",
		},
	}

	for _, tc := range testCases {
		var sv semver.SV
		svs := semverparams.SVSetter{Value: &sv}

		testhelper.DiffString(t, tc.IDStr(), "semantic version number",
			svs.CurrentValue(), "")

		err := svs.SetWithVal("", tc.paramVal)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffString(t, tc.IDStr(), "semantic version number",
				svs.CurrentValue(), tc.expVal)
		}
	}

	var sv semver.SV
	checkSetterTests := []struct {
		testhelper.ID
		testhelper.ExpPanic
		svs semverparams.SVSetter
	}{
		{
			ID:  testhelper.MkID("good"),
			svs: semverparams.SVSetter{Value: &sv},
		},
		{
			ID: testhelper.MkID("panic on nil Value pointer"),
			ExpPanic: testhelper.MkExpPanic(
				": SVSetter Check failed: the Value to be set is nil"),
			svs: semverparams.SVSetter{},
		},
	}

	for _, tc := range checkSetterTests {
		panicked, panicVal := testhelper.PanicSafe(func() {
			tc.svs.CheckSetter("test")
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}

	svs := semverparams.SVSetter{Value: &sv}
	testhelper.DiffString(t, "AllowedValues", "",
		svs.AllowedValues(),
		"a semantic version number such as v1.2.3"+
			" optionally followed by non-empty lists of dot-separated"+
			" pre-release and build IDs."+
			" For instance, 'v1.2.3-a.b.c+x.y.z'."+
			" See the Semantic Versioning spec for full details.")
}
