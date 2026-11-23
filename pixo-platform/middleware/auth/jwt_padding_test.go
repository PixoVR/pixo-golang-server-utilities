package auth_test

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("ParseJWT accepts base64url padding in JSON JWS", func(general bool, paddedField string) {
	padded := map[string]string{
		"payload":   base64.URLEncoding.EncodeToString([]byte(`{"userId":42,"orgId":7}`)),
		"protected": base64.URLEncoding.EncodeToString([]byte(`{"alg":"HS256","kid":"a"}`)),
		"signature": base64.URLEncoding.EncodeToString([]byte("sig!")),
	}
	Expect(padded[paddedField]).To(HaveSuffix("="))
	encoded := func(field string) string {
		if field == paddedField {
			return padded[field]
		}
		return strings.TrimRight(padded[field], "=")
	}
	envelope := map[string]any{"payload": encoded("payload")}
	if general {
		envelope["signatures"] = []map[string]string{
			{"protected": strings.TrimRight(padded["protected"], "="), "signature": strings.TrimRight(padded["signature"], "=")},
			{"protected": encoded("protected"), "signature": encoded("signature")},
		}
	} else {
		envelope["protected"] = encoded("protected")
		envelope["signature"] = encoded("signature")
	}
	serialized, err := json.Marshal(envelope)
	Expect(err).NotTo(HaveOccurred())

	raw, err := auth.ParseJWT(string(serialized))
	Expect(err).NotTo(HaveOccurred())
	Expect(raw).To(Equal(auth.RawToken{UserID: 42, OrgID: 7}))
},
	Entry("flattened payload", false, "payload"),
	Entry("flattened protected header", false, "protected"),
	Entry("flattened signature", false, "signature"),
	Entry("general payload", true, "payload"),
	Entry("general protected header", true, "protected"),
	Entry("general signature", true, "signature"),
)
