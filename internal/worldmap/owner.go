package worldmap

type Owner struct {
	Id   int64
	Type OwnerType
}

func NewOwner(id int64, ownerType OwnerType) *Owner {
	return &Owner{
		Id:   id,
		Type: ownerType,
	}
}

func NewPlayerOwner(id int64) *Owner {
	return NewOwner(id, OwnerType_Player)
}

func NewUnionOwner(id int64) *Owner {
	return NewOwner(id, OwnerType_Union)
}

func NewNpcOwner(id int64) *Owner {
	return NewOwner(id, OwnerType_Npc)
}

func (w *Owner) GetRelation(other Owner) Relation {
	if other.Id == w.Id {
		return Relation_Self
	}

	if w.Type == OwnerType_Union && other.Type == OwnerType_Union &&
		other.Id == w.Id {
		return Relation_Union
	}

	return Relation_Enemy
}
