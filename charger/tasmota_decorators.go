package charger

// Code generated by github.com/evcc-io/evcc/cmd/tools/decorate.go. DO NOT EDIT.

import (
	"github.com/evcc-io/evcc/api"
)

func decorateTasmota(base *Tasmota, phaseVoltages func() (float64, float64, float64, error), phaseCurrents func() (float64, float64, float64, error)) api.Charger {
	switch {
	case phaseCurrents == nil && phaseVoltages == nil:
		return base
	}

	return nil
}

type decorateTasmotaPhaseCurrentsImpl struct {
	phaseCurrents func() (float64, float64, float64, error)
}

func (impl *decorateTasmotaPhaseCurrentsImpl) Currents() (float64, float64, float64, error) {
	return impl.phaseCurrents()
}

type decorateTasmotaPhaseVoltagesImpl struct {
	phaseVoltages func() (float64, float64, float64, error)
}

func (impl *decorateTasmotaPhaseVoltagesImpl) Voltages() (float64, float64, float64, error) {
	return impl.phaseVoltages()
}
