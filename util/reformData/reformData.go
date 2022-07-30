package ktMicro

import "github.com/mitchellh/mapstructure"

func ReformJsonToModel(jsonData interface{}, value interface{}) error {
	if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: value, WeaklyTypedInput: true}); err == nil {
		return decoder.Decode(jsonData)
	} else {
		return err
	}
}
