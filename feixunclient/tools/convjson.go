package tools

import "encoding/json"

//结构体转字符串
func StructToJson(stut interface{})([]byte,error){
	data,err:=json.Marshal(stut)
	if err!=nil{
		return nil,err
	}

	return data,nil
}

func JsonToStruct(data []byte,tp interface{}){
	json.Unmarshal(data,tp)
}