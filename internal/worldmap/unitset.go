package worldmap

import "sort"

type UnitSet struct {
	units []Unit
}

func NewUnitSet(initUnits []Unit) *UnitSet {
	set := &UnitSet{
		units: initUnits,
	}
	return set
}

func (us *UnitSet) cmp(a, b Unit) int {
	return int(a.GetId() - b.GetId())
}

func (us *UnitSet) less(a, b Unit) bool {
	return us.cmp(a, b) < 0
}

func (us *UnitSet) equal(a, b Unit) bool {
	return us.cmp(a, b) == 0
}

func (us *UnitSet) sort() {
	if len(us.units) <= 1 {
		return
	}

	sort.Slice(us.units, func(i, j int) bool {
		return us.less(us.units[i], us.units[j])
	})

	// 去重
	tmp := us.units[:1]
	last := tmp[0]
	for i := 1; i < len(us.units); i++ {
		v := (us.units)[i]
		if !us.equal(v, last) {
			tmp = append(tmp, v)
			last = v
		}
	}
	us.units = tmp
}

func (us *UnitSet) Len() int {
	return len(us.units)
}

func (us *UnitSet) Reset() {
	us.units = us.units[:0]
}

func (us *UnitSet) Equal(other *UnitSet) bool {
	if us.Len() != other.Len() {
		return false
	}
	for i := 0; i < us.Len(); i++ {
		if !us.equal(us.units[i], other.units[i]) {
			return false
		}
	}
	return true
}

// 弹出末尾元素
func (us *UnitSet) Pop() Unit {
	if us.Len() == 0 {
		return nil
	}
	last := us.units[us.Len()-1]
	us.units = us.units[:us.Len()-1]
	return last
}

func (us *UnitSet) GetUnits() []Unit {
	return us.units
}

// 获取指定索引的元素，不存在则返回-1
func (us *UnitSet) getIndex(unit Unit) (int, bool) {
	if len(us.units) == 0 {
		return -1, false
	}

	// 二分查找
	i, j := 0, len(us.units)
	for i < j {
		h := int(uint(i+j) >> 1)
		if us.less(us.units[h], unit) {
			i = h + 1
		} else {
			j = h
		}
	}

	return i, i < len(us.units) && us.equal(us.units[i], unit)
}

func (us *UnitSet) Search(unit Unit) int {
	index, exist := us.getIndex(unit)
	if !exist {
		return -1
	}
	return index
}

func (us *UnitSet) IsExist(unit Unit) bool {
	_, exist := us.getIndex(unit)
	return exist
}

func (us *UnitSet) Insert(unit Unit) bool {
	index, exist := us.getIndex(unit)
	if exist {
		us.units[index] = unit
		return false
	}

	us.units = append(us.units, unit)
	return true
}

func (us *UnitSet) Delete(unit Unit) bool {
	index, exist := us.getIndex(unit)
	if !exist {
		return false
	}

	us.units = append(us.units[:index], us.units[index+1:]...)
	return true

}
