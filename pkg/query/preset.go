package query

import (
	"fmt"
)

// PresetPair is a struct used by presetting / prechecking given model.
type PresetPair struct {
	Scope  *Scope
	Filter *FilterField
	Key    interface{}

	rawPreset string
	rawFilter string
	Error     error
}

func (p *PresetPair) String() string {
	return fmt.Sprintf("Preset: %s, Filter: %s", p.rawPreset, p.rawFilter)
}

func (p *PresetPair) GetPair() (*Scope, *FilterField) {
	scope := p.Scope.copy(true, nil)
	filter := p.Filter.copy()
	return scope, filter
}

// WithKey sets the key value for the given preset pair.
// Used to get value from contexts that matches given pair.
func (p *PresetPair) WithKey(key interface{}) *PresetPair {
	p.Key = key
	return p
}

func (p *PresetPair) ErrOnFail(err error) *PresetPair {
	p.Error = err
	return p
}

func (p *PresetPair) GetStringKey() string {
	if p.Key == nil {
		return ""
	}
	strKey, ok := p.Key.(string)
	if !ok {
		return ""
	}
	return strKey
}

// PresetFilter is used to point and filter the target scope with preset values.
// It is a FilterField with a Key field.
// The Key field is used to retrieve the value for given filter from the context.
type PresetFilter struct {
	*FilterField
	Key interface{}
}

// WithKey sets the key value for the given preset filter.
// Used to get value from the context.
func (p *PresetFilter) WithKey(key interface{}) *PresetFilter {
	p.Key = key
	return p
}

// GetStringKey gets the key for given preset filter in the string form.
func (p *PresetFilter) GetStringKey() string {
	if p.Key == nil {
		return ""
	}
	strKey, ok := p.Key.(string)
	if !ok {
		return ""
	}
	return strKey
}
