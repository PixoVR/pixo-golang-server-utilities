package matchmaking_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"open-match.dev/open-match/pkg/pb"
)

var _ = Describe("PixoTicket", func() {

	var (
		testTicket *pb.Ticket
	)

	BeforeEach(func() {
		testExtenions := make(map[string]*any.Any)
		val, _ := ptypes.MarshalAny(&wrappers.Int32Value{Value: int32(1)})
		testExtenions[matchmaking.TicketMatchAttemptExtensionKey] = val

		testTicket = &pb.Ticket{
			SearchFields: &pb.SearchFields{
				StringArgs: map[string]string{
					"attributes.moduleId": "moduleID",
					"attributes.orgId":    "orgID",
				},
			},
			Extensions: testExtenions,
		}
	})

	It("can set the matchmaking attempt count on the ticket", func() {
		pixoTicket := matchmaking.NewPixoTicket(testTicket)
		pixoTicket.SetMatchmakingAttemptCount(2)

		var val wrappers.Int32Value
		err := ptypes.UnmarshalAny(pixoTicket.PersistentField[matchmaking.TicketMatchAttemptExtensionKey], &val)
		Expect(err).NotTo(HaveOccurred())

		Expect(int(val.Value)).To(Equal(2))
	})

	It("can get the matchmaking attempt count from the ticket", func() {
		pixoTicket := matchmaking.NewPixoTicket(testTicket)

		pixoTicket.SetMatchmakingAttemptCount(44)
		count, err := pixoTicket.GetMatchmakingAttemptCount()
		Expect(err).NotTo(HaveOccurred())

		Expect(count).To(Equal(int32(44)))
	})

})
