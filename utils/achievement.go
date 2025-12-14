package utils

var allowedDetails = map[string][]string{
	"competition": {
		"competition_name",
		"competition_level",
		"rank",
		"medalType",
		"eventDate",
		"location",
		"organizer",
	},
	"publication": {
		"publication_type",
		"publication_title",
		"authors",
		"publisher",
		"issn",
	},
	"organization": {
		"organization_name",
		"position",
		"period",
	},
	"certification": {
		"certification_name",
		"issuedBy",
		"certification_number",
		"validUntil",
	},
}

func FilterDetails(
	achievementType string,
	input map[string]interface{},
) map[string]interface{} {

	if input == nil {
		return map[string]interface{}{}
	}

	allowed, ok := allowedDetails[achievementType]
	if !ok {
		return map[string]interface{}{}
	}

	set := make(map[string]bool)
	for _, f := range allowed {
		set[f] = true
	}

	clean := make(map[string]interface{})
	for k, v := range input {
		if set[k] {
			clean[k] = v
		}
	}

	return clean
}
