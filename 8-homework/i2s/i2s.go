package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	inType := reflect.TypeOf(out)
	if inType.Kind() != reflect.Ptr {
		return fmt.Errorf("Type not pointer")
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for key, value := range dataMap {
			result[key] = value
		}

		res, err := json.Marshal(result)
		if err != nil {
			log.Println("failed marshaling map:", err)
			return err
		}
		// log.Println(string(res))

		err = json.Unmarshal(res, &out)
		if err != nil {
			log.Println("failed unmarshaling map:", err)
			return err
		}
	}

	if dataSlice, ok := data.([]interface{}); ok {
		result := make([]map[string]interface{}, 0)

		for _, item := range dataSlice {
			if dataMap, ok := item.(map[string]interface{}); ok {
				middleMap := make(map[string]interface{})
				for key, value := range dataMap {
					middleMap[key] = value
				}
				result = append(result, middleMap)
			}

			res, err := json.Marshal(result)
			if err != nil {
				log.Println("failed marshaling slice:", err)
				return err
			}
			// log.Println(string(res))

			err = json.Unmarshal(res, &out)
			if err != nil {
				log.Println("failed unmarshaling slice:", err)
				return err
			}
		}
	}
	return nil
}
