package matchmaking_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"open-match.dev/open-match/pkg/pb"
	"time"

	"github.com/golang/protobuf/ptypes/any"
)

var _ = Describe("Matchmaking Profile", func() {

	var (
		matchmakingProfile *matchmaking.MultiplayerMatchmakingProfile
		maxNumberOfPlayers int
	)

	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
		maxNumberOfPlayers = rand.Intn(25) + 1

		extensions := make(map[string]*any.Any)
		val, err := ptypes.MarshalAny(&wrappers.Int32Value{Value: int32(maxNumberOfPlayers)})
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}

		extensions[matchmaking.MaxPlayersExtensionKey] = val
		matchmakingProfile = matchmaking.NewMatchmakingProfile(&pb.MatchProfile{
			Extensions: extensions,
		})

	})

	It("can get the max number of players per match for the pixo matchmaking profile", func() {
		maxPlayers := matchmakingProfile.GetMaxPlayersPerMatch()
		Expect(maxPlayers).To(Equal(maxNumberOfPlayers))
	})
})
