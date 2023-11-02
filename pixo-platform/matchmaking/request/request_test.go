package request_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/matchmaking/request"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {

	var (
		matchRequest = request.MatchRequest{
			BaseTicketRequest: request.BaseTicketRequest{
				Capacity: 1,
			},
			OrgID:         1,
			ModuleID:      1,
			ServerVersion: "1.00.00",
		}
		matchRequestString = `{"capacity":1,"moduleId":1,"orgId":1,"serverVersion":"1.00.00"}`

		partyCodeMatchRequest = request.PartyMatchRequest{
			BaseTicketRequest: request.BaseTicketRequest{
				Capacity: 1,
			},
			PartyCode: "1234",
		}
		partyCodeMatchRequestString = `{"capacity":1,"partyCode":"1234"}`
	)

	It("can format the label for a match request", func() {
		Expect(matchRequest.GetLabel()).To(Equal("o-1-m-1-v-1.00.00"))
	})

	It("can marshal a match request", func() {
		data, err := matchRequest.MarshalJSON()
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(string(data)).To(BeEquivalentTo(matchRequestString))
	})

	It("can unmarshal a match request", func() {
		req := request.MatchRequest{}
		err := req.UnmarshalJSON([]byte(matchRequestString))
		Expect(err).To(BeNil())
		Expect(req).To(Equal(matchRequest))
	})

	It("can marshal a party code match request", func() {
		data, err := partyCodeMatchRequest.MarshalJSON()
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(string(data)).To(BeEquivalentTo(partyCodeMatchRequestString))
	})

	It("can unmarshal a party code match request", func() {
		req := request.PartyMatchRequest{}
		err := req.UnmarshalJSON([]byte(partyCodeMatchRequestString))
		Expect(err).To(BeNil())
		Expect(req).To(Equal(partyCodeMatchRequest))
	})

})
