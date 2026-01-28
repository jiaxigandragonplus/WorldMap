package worldmap

// unit对象管理器
type UnitManager struct {
	units map[int64]Unit // id -> Unit
}

func NewUnitManager() *UnitManager {
	return &UnitManager{
		units: make(map[int64]Unit),
	}
}
