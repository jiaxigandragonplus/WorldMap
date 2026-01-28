package worldmap

type MapPlayerManager struct {
	players map[int64]*MapPlayer
}

func NewMapPlayerManager() *MapPlayerManager {
	return &MapPlayerManager{
		players: make(map[int64]*MapPlayer),
	}
}
