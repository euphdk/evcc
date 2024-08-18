package charger

// Code generated by github.com/evcc-io/evcc/cmd/tools/decorate.go. DO NOT EDIT.

import (
	"github.com/evcc-io/evcc/api"
)

func decorateEm2Go(base *Em2Go, chargerEx func(float64) error, phaseSwitcher func(int) error, phaseGetter func() (int, error)) api.Charger {
	switch {
	case chargerEx == nil && phaseGetter == nil && phaseSwitcher == nil:
		return base

	case chargerEx != nil && phaseGetter == nil && phaseSwitcher == nil:
		return &struct {
			*Em2Go
			api.ChargerEx
		}{
			Em2Go: base,
			ChargerEx: &decorateEm2GoChargerExImpl{
				chargerEx: chargerEx,
			},
		}

	case chargerEx == nil && phaseGetter == nil && phaseSwitcher != nil:
		return &struct {
			*Em2Go
			api.PhaseSwitcher
		}{
			Em2Go: base,
			PhaseSwitcher: &decorateEm2GoPhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}

	case chargerEx != nil && phaseGetter == nil && phaseSwitcher != nil:
		return &struct {
			*Em2Go
			api.ChargerEx
			api.PhaseSwitcher
		}{
			Em2Go: base,
			ChargerEx: &decorateEm2GoChargerExImpl{
				chargerEx: chargerEx,
			},
			PhaseSwitcher: &decorateEm2GoPhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}

	case chargerEx == nil && phaseGetter != nil && phaseSwitcher != nil:
		return &struct {
			*Em2Go
			api.PhaseGetter
			api.PhaseSwitcher
		}{
			Em2Go: base,
			PhaseGetter: &decorateEm2GoPhaseGetterImpl{
				phaseGetter: phaseGetter,
			},
			PhaseSwitcher: &decorateEm2GoPhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}

	case chargerEx != nil && phaseGetter != nil && phaseSwitcher != nil:
		return &struct {
			*Em2Go
			api.ChargerEx
			api.PhaseGetter
			api.PhaseSwitcher
		}{
			Em2Go: base,
			ChargerEx: &decorateEm2GoChargerExImpl{
				chargerEx: chargerEx,
			},
			PhaseGetter: &decorateEm2GoPhaseGetterImpl{
				phaseGetter: phaseGetter,
			},
			PhaseSwitcher: &decorateEm2GoPhaseSwitcherImpl{
				phaseSwitcher: phaseSwitcher,
			},
		}
	}

	return nil
}

type decorateEm2GoChargerExImpl struct {
	chargerEx func(float64) error
}

func (impl *decorateEm2GoChargerExImpl) MaxCurrentMillis(p0 float64) error {
	return impl.chargerEx(p0)
}

type decorateEm2GoPhaseGetterImpl struct {
	phaseGetter func() (int, error)
}

func (impl *decorateEm2GoPhaseGetterImpl) GetPhases() (int, error) {
	return impl.phaseGetter()
}

type decorateEm2GoPhaseSwitcherImpl struct {
	phaseSwitcher func(int) error
}

func (impl *decorateEm2GoPhaseSwitcherImpl) Phases1p3p(p0 int) error {
	return impl.phaseSwitcher(p0)
}
