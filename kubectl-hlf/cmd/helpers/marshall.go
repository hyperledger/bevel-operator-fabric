package helpers

import (
	"sigs.k8s.io/yaml"
)

func MarshallWithoutStatus(item interface{}) ([]byte, error) {
	ot, err := yaml.Marshal(item)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = yaml.Unmarshal(ot, &data)
	if err != nil {
		return nil, err
	}
	delete(data, "status")
	ot, err = yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return ot, nil
}
