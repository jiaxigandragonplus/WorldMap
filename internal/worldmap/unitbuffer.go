package worldmap

type UnitBuffer struct {
	Units []Unit
}

func NewUnitBuffer() *UnitBuffer {
	return &UnitBuffer{
		Units: make([]Unit, 0),
	}
}

func (ub *UnitBuffer) AppendUnit(unit Unit) {
	ub.Units = append(ub.Units, unit)
}

func (ub *UnitBuffer) AppendUnits(units []Unit) {
	ub.Units = append(ub.Units, units...)
}

func (ub *UnitBuffer) Reset() {
	ub.Units = make([]Unit, 0)
}

func (ub *UnitBuffer) GetUnits() []Unit {
	return ub.Units
}
