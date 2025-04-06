package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Simple struct {
	ID       int
	Username string
	Active   bool
}

type IDBlock struct {
	ID int
}

type Complex struct {
	SubSimple  Simple
	ManySimple []Simple
	Blocks     []IDBlock
}

func i2s(data interface{}, out interface{}) error {
	dataType := reflect.TypeOf(out)
	result := make(map[string]interface{})
	if dataMap, ok := data.(map[string]interface{}); ok {
		for key, value := range dataMap {
			result[key] = value
			log.Println(key, value)
		}
	}
	log.Println(dataType)
	res, err := json.Marshal(result)
	if err != nil {
		return err
	}
	log.Println(string(res))
	err = json.Unmarshal(res, &out)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	fmt.Println("simple:", smpl)
	// expected := &Complex{
	// 	SubSimple:  smpl,
	// 	ManySimple: []Simple{smpl, smpl},
	// 	Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	// }

	jsonRaw, _ := json.Marshal(smpl)
	// fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)
	fmt.Println("tmpData:", tmpData)

	result := new(Simple)
	err := i2s(tmpData, result)
	if err != nil {
		panic(err)
	}

	fmt.Println(*result)
}
