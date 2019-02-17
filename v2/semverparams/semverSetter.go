package semverparams

import (
	"errors"

	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/semver.mod/semver"
)

// SVSetter is a parameter setter which will set a semantic version number
type SVSetter struct {
	Value **semver.SV
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (svs SVSetter) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (svs SVSetter) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal checks that the parameter value meets the checks if any. It
// returns an error if the check is not satisfied. Only if the check
// is not violated is the Value set.
func (svs *SVSetter) SetWithVal(_ string, paramVal string) error {
	var err error
	*svs.Value, err = semver.ParseSV(paramVal)
	return err
}

// AllowedValues simply returns "any string" since StringSetter
// does not check its value
func (svs SVSetter) AllowedValues() string {
	return "a semantic version number such as v1.2.3" +
		" optionally followed by non-empty lists of dot-separated" +
		" pre-release and build IDs." +
		" For instance, 'v1.2.3-a.b.c+x.y.z'." +
		" See the Semantic Versioning spec for full details."
}

// CurrentValue returns the current setting of the parameter value
func (svs SVSetter) CurrentValue() string {
	if (*svs.Value) == nil {
		return "none"
	}
	return (*svs.Value).String()
}

// CheckSetter panics if the setter has not been properly created
func (svs SVSetter) CheckSetter(name string) {
	if svs.Value == nil {
		panic(name + ": SVSetter Check failed: the Value to be set is nil")
	}
}
