package main

import (
	"encoding/json"
	"fmt"
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
	dataType := reflect.TypeOf(data).Kind()

	fmt.Println(dataType)
	// middleMap := make(map[string]any)
	// dataMap, ok := data.(map[string]any)
	// if !ok {
	// 	return fmt.Errorf("failed type")
	// }
	// // for k, v := range dataMap {
	// // 	fmt.Println(k, v)
	// // }
	// dataInter := dataMap.(reflect.Value)
	// bytesData := dataInter.Bytes()
	// err := json.Unmarshal(bytesData, &middleMap)
	// if err != nil {
	// 	return err
	// }
	// out = middleMap

	return nil
}

func main() {
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	// expected := &Complex{
	// 	SubSimple:  smpl,
	// 	ManySimple: []Simple{smpl, smpl},
	// 	Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	// }

	jsonRaw, _ := json.Marshal(smpl)
	fmt.Println(string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)
	fmt.Println(tmpData)

	result := new(Complex)
	err := i2s(tmpData, result)
	if err != nil {
		panic(err)
	}

	fmt.Println(*result)
}
