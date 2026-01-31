package worldmap

// unit对象管理器
type UnitManager struct {
	units       map[int64]Unit                 // id -> Unit
	unitsByType map[MapUnitType]map[int64]Unit // 类型 -> id -> Unit
}

func NewUnitManager() *UnitManager {
	return &UnitManager{
		units:       make(map[int64]Unit),
		unitsByType: make(map[MapUnitType]map[int64]Unit),
	}
}

func (mgr *UnitManager) GetUnitById(id int64) Unit {
	return mgr.units[id]
}

func (mgr *UnitManager) GetUnitByType(unitType MapUnitType) map[int64]Unit {
	return mgr.unitsByType[unitType]
}

// 添加unit
func (mgr *UnitManager) AddUnit(unit Unit) {
	mgr.units[unit.GetId()] = unit
	if units, ok := mgr.unitsByType[unit.GetType()]; ok {
		units[unit.GetId()] = unit
	} else {
		mgr.unitsByType[unit.GetType()] = map[int64]Unit{unit.GetId(): unit}
	}
}

// 删除unit
func (mgr *UnitManager) RemoveUnit(unit Unit) {
	delete(mgr.units, unit.GetId())
	delete(mgr.unitsByType[unit.GetType()], unit.GetId())
}
