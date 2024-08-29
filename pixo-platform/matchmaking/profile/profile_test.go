package profile_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/profile"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math/rand"
	"open-match.dev/open-match/pkg/pb"
)

var _ = Describe("Matchmaking Profile", func() {

	var (
		matchmakingProfile *profile.OpenMatchProfile
		maxNumberOfPlayers int
	)

	BeforeEach(func() {
		maxNumberOfPlayers = rand.Intn(25) + 1

		val, err := anypb.New(&wrapperspb.Int32Value{Value: int32(maxNumberOfPlayers)})
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}

		matchmakingProfile = profile.NewOpenMatchProfile(&pb.MatchProfile{
			Extensions: map[string]*anypb.Any{
				profile.MaxPlayersExtensionKey: val,
			},
		})

	})

	It("can get the max number of players per match for the pixo matchmaking profile", func() {
		maxPlayers := matchmakingProfile.GetMaxPlayersPerMatch()
		Expect(maxPlayers).To(Equal(maxNumberOfPlayers))
	})

})
