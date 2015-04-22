package main

import "fmt"
import "encoding/json"
//import "time"
//import "reflect"

type innerdata struct{
	Id int
	Sa int
}

type innerdata2 struct{
	BE int
}


type data struct{
	Typeof string
	Data innerdata
}


var unjson data

func main(){
	data2 := innerdata{2,6}

	test := data{"dette var kult",data2}

	fmt.Println(test)

	testj, _ := json.Marshal(test)

	fmt.Println(testj)
	

	//time.Sleep(2 * time.Second)

	err := json.Unmarshal(testj, &unjson)
	if err != nil{
		fmt.Println("ERROR!")
	}
	fmt.Println(unjson)

	fmt.Println(unjson.Data.Id)


}