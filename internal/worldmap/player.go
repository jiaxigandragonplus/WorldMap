package worldmap

type MapPlayer struct {
	PlayerId int64 // 玩家ID
	CityId   int64 // 玩家主城id
}

func NewMapPlayer(playerId int64, cityId int64) *MapPlayer {
	return &MapPlayer{
		PlayerId: playerId,
		CityId:   cityId,
	}
}
