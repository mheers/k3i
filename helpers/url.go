package helpers

import "net/url"

func MapToUrlValues(data map[string]string) url.Values {
	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}
	return values
}
