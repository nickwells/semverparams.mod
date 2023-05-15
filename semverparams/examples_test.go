package semverparams_test

import (
	"fmt"

	"github.com/nickwells/param.mod/v5/param/paramset"
	"github.com/nickwells/semver.mod/v3/semver"
	"github.com/nickwells/semverparams.mod/v5/semverparams"
)

// Example demonstrates how to use the package to add semantic versioning
// parameters to your program interface
func Example() {
	svp := semverparams.SemverVals{SemVer: &semver.SV{}}
	ps := paramset.NewOrDie(
		semverparams.AddSemverGroup,
		svp.AddSemverParam(nil),
	)

	// ps.Parse()
	//
	// Supplying arguments to Parse is only necessary for the example. In
	// production code you should pass nothing to Parse (as in the first line
	// of this comment) in which case it will use the program arguments.
	ps.Parse([]string{"-semver", "v1.2.3"})
	if svp.SemVer.HasBeenSet() {
		fmt.Println(svp.SemVer.String())
	}
	// Output:
	// v1.2.3
}

// Example_withChecks demonstrates how to use the package to add semantic
// versioning parameters to your program interface. This example shows how
// you can use user-supplied checks to validate the build and pre-release
// IDs. We pass a non-nil SemverChecks pointer to AddSemverParam and that
// generates a FinalChecks func which checks that the pre-release and build
// IDs pass the supplied checks.
func Example_withChecks() {
	svc := semverparams.SemverChecks{}
	svp := semverparams.SemverVals{SemVer: &semver.SV{}}

	ps := paramset.NewOrDie(
		semverparams.AddSemverGroup,
		svp.AddSemverParam(&svc),
		svc.AddCheckParams(),
	)

	// ps.Parse()
	//
	// Supplying arguments to Parse is only necessary for the example. In
	// production code you should call Parse with no parameters (as in the
	// first line of this comment) in which case it will use the program
	// arguments.
	ps.Parse([]string{
		"-semver", "v1.2.3-rc.1",
		"-pre-rel-ID-checks", `Or(Length(EQ(0)), Length(EQ(2)))`,
	})
	if svp.SemVer.HasBeenSet() {
		fmt.Println(svp.SemVer.String())
	}
	// Output:
	// v1.2.3-rc.1
}

// Example_withFailingChecks demonstrates how to use the package to add
// semantic versioning parameters to your program interface. This example
// shows the behaviour when user-supplied checks fail.
func Example_withFailingChecks() {
	svc := semverparams.SemverChecks{}
	svp := semverparams.SemverVals{SemVer: &semver.SV{}}

	// we use a testing-specific paramset generator to suppress the
	// exit-on-error and error reporting behaviour
	ps := paramset.NewNoHelpNoExitNoErrRptOrPanic(
		semverparams.AddSemverGroup,
		svp.AddSemverParam(&svc),
		svc.AddCheckParams(),
	)

	// ps.Parse()
	//
	// Supplying arguments to Parse is only necessary for the example. In
	// production code you should call Parse with no parameters (as in the
	// first line of this comment) in which case it will use the program
	// arguments.
	errs := ps.Parse([]string{
		"-semver", "v1.2.3-invalid",
		"-pre-rel-ID-checks", `Or(Length(EQ(0)), Length(EQ(2)))`,
	})
	if len(errs) != 0 {
		// We just print the start of the first error. The normal behaviour
		// of param.PSet is to print the whole error in a user-accessible
		// format and then exit
		fmt.Printf("Error: %-13.13s\n", errs["Final Checks"][0])
	}
	// Output:
	// Error: Bad PreRelIDs
}

// Example_unchecked demonstrates how to use the package to add semantic
// versioning parameters to your program interface. This example shows how
// you can have multiple parameters, some checked, some not.
func Example_unchecked() {
	svc := semverparams.SemverChecks{}
	svp1 := semverparams.SemverVals{SemVer: &semver.SV{}}
	svp2 := semverparams.SemverVals{Prefix: "sv2", SemVer: &semver.SV{}}

	ps := paramset.NewOrDie(
		semverparams.AddSemverGroup,
		svp1.AddSemverParam(&svc),
		svp2.AddSemverParam(nil),
		svc.AddCheckParams(),
	)

	// ps.Parse()
	//
	// Supplying arguments to Parse is only necessary for the example. In
	// production code you should call Parse with no parameters (as in the
	// first line of this comment) in which case it will use the program
	// arguments.
	ps.Parse([]string{
		"-semver", "v1.2.3-rc.1",
		"-sv2-semver", "v2.3.4-unchecked",
		"-pre-rel-ID-checks", `Or(Length(EQ(0)), Length(EQ(2)))`,
	})
	if svp1.SemVer.HasBeenSet() {
		fmt.Println(svp1.SemVer.String())
	}
	if svp2.SemVer.HasBeenSet() {
		fmt.Println(svp2.SemVer.String())
	}
	// Output:
	// v1.2.3-rc.1
	// v2.3.4-unchecked
}
