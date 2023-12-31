package cache_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/cache"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/request"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProfileRepository", Ordered, func() {

	var (
		ctx              context.Context
		s                *miniredis.Miniredis
		c                *redis.Client
		gameProfileCache *cache.GameProfileCacheClient

		PartyCodeRequest = request.PartyMatchRequest{
			BaseTicketRequest: request.BaseTicketRequest{Capacity: 1},
			PartyCode:         "test",
		}

		PixoRequest = request.MatchRequest{
			BaseTicketRequest: request.BaseTicketRequest{Capacity: 1},
			OrgID:             1,
			ModuleID:          1,
			ServerVersion:     "1.00.00",
		}
	)

	BeforeEach(func() {
		var err error
		gameProfileCache, s, c, err = cache.NewMiniGameProfileCache()
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
	})

	AfterEach(func() {
		s.Close()
	})

	It("can save a valid basic matchmaking profile in the cache", func() {
		Expect(gameProfileCache).ToNot(BeNil())
		err := gameProfileCache.SaveProfile(ctx, &PartyCodeRequest)
		Expect(err).To(BeNil())

		formattedKey := PartyCodeRequest.GetLabel()
		keys := c.Keys(ctx, "*").Val()
		Expect(len(keys)).Should(BeNumerically(">", 0))
		Expect(keys[0]).To(ContainSubstring(formattedKey))
	})

	It("can save a valid pixo matchmaking profile in the cache", func() {
		Expect(gameProfileCache).ToNot(BeNil())
		err := gameProfileCache.SaveProfile(ctx, &PixoRequest)
		Expect(err).To(BeNil())

		formattedKey := PixoRequest.GetLabel()
		keys := c.Keys(ctx, "*").Val()
		Expect(len(keys)).Should(BeNumerically(">", 0))
		Expect(keys[0]).To(ContainSubstring(formattedKey))
	})

	It("returns all saved profiles", func() {
		Expect(gameProfileCache).ToNot(BeNil())
		err := gameProfileCache.SaveProfile(ctx, &PixoRequest)
		Expect(err).To(BeNil())

		profiles, err := gameProfileCache.GetAllProfiles(ctx)
		Expect(err).To(BeNil())
		Expect(len(profiles)).Should(BeNumerically(">", 0))

		Expect(profiles[0]).To(Equal(`{"capacity":1,"moduleId":1,"orgId":1,"serverVersion":"1.00.00"}`))
	})
})
