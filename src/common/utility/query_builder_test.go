package utility_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "xtremepay.com/backoffice/utility"
)

var _ = Describe("QueryBuilder", func() {
	var (
		queryBuilder DBQueryBuilder
	)

	BeforeEach(func() {
		queryBuilder = DBQueryBuilder{}
	})

	Describe("Extract Query parameter details from q in HTTP QueryString", func() {
		Context("Confirm param came back as a slice of QueryStruct", func() {
			It("should be a novel", func() {
				Expect(queryBuilder.ParseQueryFilter(`{"filters":[{"name":"age","op":"ge","val":10}, {"name":"age","op":"ge","val":10}, {"name":"age","op":"ge","val":10}]}`)).To(Equal("NOVEL"))
			})
		})

	})

})
