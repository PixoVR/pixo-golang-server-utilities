package common_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("DateRange", func() {
	Context("SetDefaultDateRange", func() {
		It("can return a default date range for the last 30 days if neither start nor end date is provided", func() {
			dateRangeStart, dateRangeEnd := common.SetDefaultDateRange(nil, nil)
			Expect(dateRangeStart).To(BeTemporally("~", dateRangeEnd.AddDate(0, 0, -30), 1*time.Second))
			Expect(dateRangeEnd).To(BeTemporally("~", time.Now(), 1*time.Second))
		})

		It("can return a default date range for the last 30 days if only the end date is provided", func() {
			endDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			dateRangeStart, dateRangeEnd := common.SetDefaultDateRange(nil, &endDate)
			Expect(dateRangeStart).To(BeTemporally("~", endDate.AddDate(0, 0, -30), 1*time.Second))
			Expect(dateRangeEnd).To(Equal(endDate))
		})

		It("can return a default date range from the start date to the current date if only the start date is provided", func() {
			startDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			dateRangeStart, dateRangeEnd := common.SetDefaultDateRange(&startDate, nil)
			Expect(dateRangeStart).To(Equal(startDate))
			Expect(dateRangeEnd).To(BeTemporally("~", time.Now(), 1*time.Second))
		})
	})

	DescribeTable("can get date 13 months from next month", func(startDate time.Time) {
		newDate := common.GetDateThirteenMonthsFromNextMonth(startDate)
		Expect(newDate.Month()).To(Equal(startDate.AddDate(0, 2, 0).Month()))
		Expect(newDate.Year()).To(Equal(startDate.AddDate(1, 0, 0).Year()))
		Expect(newDate.Day()).To(Equal(1))
	},
		Entry("2021-01-01", time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
		Entry("2025-02-01", time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)),
		Entry("1999-03-01", time.Date(1999, 3, 2, 0, 0, 0, 0, time.UTC)),
	)

})
