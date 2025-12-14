package utils

var forbiddenMongoKeys = []string{
	"$where",
	"$expr",
	"$function",
	"$accumulator",
	"$lookup",
	"$graphLookup",
}

func SanitizeMongoMap(input map[string]interface{}) map[string]interface{} {
	clean := make(map[string]interface{})

	for k, v := range input {
		if len(k) > 0 && k[0] == '$' {
			continue
		}

		for _, forbidden := range forbiddenMongoKeys {
			if k == forbidden {
				goto skip
			}
		}

		clean[k] = v
	skip:
	}

	return clean
}
