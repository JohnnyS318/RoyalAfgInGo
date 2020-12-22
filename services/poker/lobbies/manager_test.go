package lobbies

import (
	"testing"
)

func TestAppendLobby(t *testing.T) {

	manager := NewManager(2, [][]int{{5, 10}, {11, 25}, {26, 50}})

	e, err := manager.AppendLobby(0)
	if err != nil {
		t.Error(err)
	}
	lobby, ok := manager.Lobbies[e]
	if !ok || lobby == nil {
		t.Errorf("No Lobby got created")
	}
}