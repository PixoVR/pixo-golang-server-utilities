package matchmaking_test

import (
	"fmt"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking"
	"github.com/alicebob/miniredis/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ValidRequest = matchmaking.TicketRequest{
		MatchRequest: matchmaking.MatchRequest{
			OrgID:         1,
			ModuleID:      1,
			ServerVersion: "1.00.00",
		},
		Capacity: 1,
	}
)

var _ = Describe("ProfileRepository", Ordered, func() {
	var redis *miniredis.Miniredis
	var profileRepository *matchmaking.GameProfileRepository

	BeforeAll(func() {
		redisMock, err := miniredis.Run()
		Expect(err).ToNot(HaveOccurred())
		redis = redisMock
	})

	AfterAll(func() {
		redis.Close()
	})

	BeforeEach(func() {
		profileRepository = matchmaking.NewGameProfileRepository(redis.Addr(), "")
	})

	It("should connect to redis", func() {
		profileRepository := matchmaking.NewGameProfileRepository(redis.Addr(), "")
		Expect(profileRepository).ToNot(BeNil())
	})

	It("returns no error when saving a profile to redis", func() {
		Expect(profileRepository).ToNot(BeNil())
		err := profileRepository.SaveProfile(ValidRequest)
		Expect(err).To(BeNil())
	})

	It("saves org id, module id, server version", func() {
		Expect(profileRepository).ToNot(BeNil())
		err := profileRepository.SaveProfile(ValidRequest)
		Expect(err).To(BeNil())
		Expect(len(redis.Keys())).Should(BeNumerically(">", 0))
	})

	It("uses org id, module id, server version in the key", func() {
		Expect(profileRepository).ToNot(BeNil())
		err := profileRepository.SaveProfile(ValidRequest)
		Expect(err).To(BeNil())

		formattedKey := fmt.Sprintf("profile:%d%d%s",
			ValidRequest.OrgID,
			ValidRequest.ModuleID,
			ValidRequest.ServerVersion,
		)
		key := redis.Keys()[0]
		Expect(key).To(ContainSubstring(formattedKey))
	})

	It("returns all saved profiles", func() {
		Expect(profileRepository).ToNot(BeNil())
		err := profileRepository.SaveProfile(ValidRequest)
		Expect(err).To(BeNil())

		profiles, err := profileRepository.GetAllProfiles()
		Expect(err).To(BeNil())
		Expect(len(profiles)).Should(BeNumerically(">", 0))

		Expect(profiles[0]).To(Equal(ValidRequest))
	})
})
