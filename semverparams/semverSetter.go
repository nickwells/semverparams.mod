package semverparams

import (
	"github.com/nickwells/param.mod/v6/psetter"
	"github.com/nickwells/semver.mod/v3/semver"
)

// SVSetter is a parameter setter which will set a semantic version
// number. It satisfies the param.Setter interface and so can be used when
// specifying a command line argument using the param package.
type SVSetter struct {
	psetter.ValueReqMandatory

	Value *semver.SV
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (svs SVSetter) SetWithVal(_ string, paramVal string) error {
	v, err := semver.ParseSV(paramVal)
	if err != nil {
		return err
	}

	v.CopyInto(svs.Value)

	return nil
}

// AllowedValues returns a description of the allowed values
func (svs SVSetter) AllowedValues() string {
	return "a semantic version number such as v1.2.3" +
		" optionally followed by non-empty lists of dot-separated" +
		" pre-release and build IDs." +
		" For instance, 'v1.2.3-a.b.c+x.y.z'." +
		" See the Semantic Versioning spec for full details."
}

// CurrentValue returns the current setting of the parameter value
func (svs SVSetter) CurrentValue() string {
	return svs.Value.String()
}

// CheckSetter panics if the setter has not been properly created
func (svs SVSetter) CheckSetter(name string) {
	if svs.Value == nil {
		panic(name + ": SVSetter Check failed: the Value to be set is nil")
	}
}
