package helpers

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

func YamlToMap(y []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := yaml.Unmarshal(y, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func JsonBytesToYamlBytes(b []byte) ([]byte, error) {
	return yaml.JSONToYAML(b)
}

func YamlBytesToJSONBytes(b []byte) ([]byte, error) {
	return yaml.YAMLToJSON(b)
}

func MarshalViaJSONToYAML(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return JsonBytesToYamlBytes(b)
}
