package lobby

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	allocationv1 "agones.dev/agones/pkg/apis/allocation/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/JohnnyS318/RoyalAfgInGo/pkg/poker/models"
)

//NewLobby allocates a new GameServer for a new RoundId
func (m *Manager) NewLobby(classIndex int) (*TicketRequestResult, error) {

	if classIndex < 0 || classIndex >= len(m.classes) {
		return nil, errors.New("the class index has to be a valid class registered at service start")
	}

	id := newID()
	//model := models.NewLobby(id, class)

	gsa := m.agonesClient.AllocationV1()

	class := m.classes[classIndex]
	serverLabels := make(map[string]string, 3)
	serverLabels["lobbyId"] = id
	serverLabels["min-buy-in"] = strconv.Itoa(class.Min)
	serverLabels["max-buy-in"] = strconv.Itoa(class.Max)
	serverLabels["blind"] = strconv.Itoa(class.Blind)
	serverLabels["class-index"] = strconv.Itoa(classIndex)

	alloc := &allocationv1.GameServerAllocation{
		ObjectMeta: v1.ObjectMeta{
			Name:      id,
			Namespace: "default",
		},
		Spec: allocationv1.GameServerAllocationSpec{
			Required: v1.LabelSelector{
				MatchLabels: map[string]string{
					"game": "poker",
				},
				MatchExpressions: nil,
			},
			Preferred: nil,
			MetaPatch: allocationv1.MetaPatch{
				Labels:      serverLabels,
				Annotations: nil,
			},
		}}

	allocationResponse, err := gsa.GameServerAllocations("default").Create(alloc)

	m.logger.Warnw("Allocation", "error", err, "lobbyId", id)

	if err != nil {
		return nil, err
	}

	if allocationResponse.Status.GameServerName == "" || len(allocationResponse.Status.Ports) <= 0 {
		return nil, errors.New("no new server can be allocated")
	}

	ip := allocationResponse.Status.Address
	port := allocationResponse.Status.Ports[0].Port
	addr := fmt.Sprintf("%s:%v", ip, port)

	err = m.rdg.Set(context.Background(), id, addr, 0).Err()

	if err != nil {
		return nil, err
	}

	m.lobbies[classIndex] = append(m.lobbies[classIndex], models.LobbyBase{
		LobbyID: id,
		Class:   &m.classes[classIndex],
		ClassIndex: classIndex,
		PlayerCount: 0,
	})

	return &TicketRequestResult{
		Address: addr,
		LobbyId: id,
	}, nil
}

const idLength = 7
const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//newID generates a new ID for a new RoundId. RoundId ID are composed of letters for easy share.
func newID() string {
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(idLength)
	for i := 0; i < idLength; i++ {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}
	return sb.String()
}
