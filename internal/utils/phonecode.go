package utils

import (
	"fmt"
	"math"
	"ngevent/internal/dto"
	"sort"

	"github.com/nyaruka/phonenumbers"
)

func ListAllPhoneCodes(page, limit int) *dto.PaginationResponse[dto.PhoneCodeRespone] {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 15
	}

	regions := phonenumbers.GetSupportedRegions()
	result := make([]dto.PhoneCodeRespone, 0, len(regions))

	for region := range regions {
		code := phonenumbers.GetCountryCodeForRegion(region)
		if code == 0 {
			continue
		}

		result = append(result, dto.PhoneCodeRespone{
			ISO:      region,
			DialCode: fmt.Sprintf("+%d", code),
		})
	}

	// Sort by ISO
	sort.Slice(result, func(i, j int) bool {
		return result[i].ISO < result[j].ISO
	})

	totalData := len(result)
	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	offset := (page - 1) * limit
	if offset > totalData {
		offset = totalData
	}

	end := offset + limit
	if offset > totalData {
		end = totalData
	}

	paginatedData := result[offset:end]

	return &dto.PaginationResponse[dto.PhoneCodeRespone]{
		Data:       paginatedData,
		Page:       page,
		Limit:      limit,
		TotalData:  totalData,
		TotalPages: totalPages,
	}
}
