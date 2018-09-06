package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"encoding/json"
	"log"
)

type recievePack struct {
	BitMap [][]int
}

type returnPack struct {
	ReturnAns interface{}
}

type Loc struct {
	X int
	Y int
	TubeState int
}

var bitMap = make([][]int, 5)
var markedMap = make([][]int, 5)
var pathStack []Loc
var ansPath []Loc
var flag = false

func main(){
	fmt.Println("aaaaaaaaaaa")

	for i:= 0; i < len(bitMap); i++{
		bitMap[i] = make([]int, 5)
		markedMap[i] = make([]int, 5)
	}


	//GenerateMap()
	//bitMap[0][0] = 0
	//bitMap[0][1] = 0
	//bitMap[0][2] = 0
	//bitMap[0][3] = 0
	//bitMap[0][4] = 0
	//findSolution(bitMap)

	http.HandleFunc("/TunelGameProc",TunelGameProc)
	http.HandleFunc("/systemAns",systemAns)
	err := http.ListenAndServe(":4399", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func GenerateMap(){
	for i := 0; i < 5; i++{
		for j := 0; j < 5; j++{
			rand := rand.Int() % 6
			bitMap[i][j] = rand
		}
	}
	bitMap[3][4] = 1
}

func recieveData(w http.ResponseWriter,r *http.Request) recievePack {
	//decode package string->byte->struct
	strToByte := []byte(r.FormValue("first"))

	//convert to struct
	var calPackage recievePack;
	err :=json.Unmarshal(strToByte, &calPackage);
	if err != nil{
		fmt.Println("ERROR", err);
	}

	return calPackage;
}

func systemAnsCheck(w http.ResponseWriter,r *http.Request, flag bool){
	fmt.Println("------Start Print Response------")

	//encode to byte[]
	stringInfoInByte, err := json.Marshal(returnPack{flag})
	//convert byte[] to string
	strConverted := string(stringInfoInByte)

	json.NewEncoder(w).Encode(string(strConverted))
	//check error
	if err != nil{
		fmt.Println("ERROR")
	}
	fmt.Println(stringInfoInByte)
	fmt.Println(strConverted)
}

//when page loaded send the map first check whether has answer if no ans return false to reload
func systemAns(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//retrieve data
	recDate := recieveData(w, r)
	bitMap = recDate.BitMap
	findSolution()

	//return bool to server
	systemAnsCheck(w, r, flag);

	//reset
	flag = false
}

func TunelGameProc(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//retrieve data
	recDate := recieveData(w, r)
	bitMap = recDate.BitMap
	checkSolution()
	if flag == true{
		userAnsReturn(w, r, ansPath)
	}else{
		userAnsReturn(w, r, flag)
	}
	

	//reset
	flag = false
	ansPath = []Loc{}
}

func userAnsReturn(w http.ResponseWriter,r *http.Request, correctPath interface{}){
	fmt.Println("------Start Print Response------")

	//encode to byte[]
	stringInfoInByte, err := json.Marshal(returnPack{correctPath})
	//convert byte[] to string
	strConverted := string(stringInfoInByte)

	json.NewEncoder(w).Encode(string(strConverted))
	//check error
	if err != nil{
		fmt.Println("ERROR")
	}
	fmt.Println(stringInfoInByte)
	fmt.Println(strConverted)
}

func checkSolution(){
	checkdfs(0,0,1)
}

func checkdfs(x int, y int, waterEntry int){
	//check end
	if x == 4 && y == 5{
		flag = true
		return
	}

	//edge checking
	if (x < 0 || x > 4) || (y < 0 || y > 4) {
		return
	}

	//fmt.Println(x, y, waterEntry, bitMap[x][y])
	ansPath = append(ansPath, Loc{x, y, bitMap[x][y]})
	//straight tube
	if bitMap[x][y] == 0 {
		switch waterEntry {
		case 1:
			checkdfs(x, y + 1, 1)
			break
		case 2:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 3:
			checkdfs(x, y - 1, 3)
			break
		case 4:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		}
		//vertical tube
	}else if bitMap[x][y] == 1{
		switch waterEntry {
		case 1:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 2:
			checkdfs(x + 1, y, 2)
			break
		case 3:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 4:
			checkdfs(x - 1, y, 4)
			break
		}
		//right-bot
	}else if bitMap[x][y] == 2{
		switch waterEntry {
		case 1:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 2:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 3:
			checkdfs(x + 1, y, 2)
			break
		case 4:
			checkdfs(x, y + 1, 1)
			break
		}
		//left bot
	}else if bitMap[x][y] == 3{
		switch waterEntry {
		case 1:
			checkdfs(x + 1, y, 2)
			break
		case 2:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 3:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 4:
			checkdfs(x, y - 1, 3)
			break
		}
		//left-top
	}else if bitMap[x][y] == 4{
		switch waterEntry {
		case 1:
			checkdfs(x - 1, y, 4)
			break
		case 2:
			checkdfs(x, y - 1, 3)
			break
		case 3:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 4:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		}
		//right-top
	}else{
		switch waterEntry {
		case 1:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		case 2:
			checkdfs(x, y + 1, 1)
			break
		case 3:
			ansPath = ansPath[:len(ansPath) - 1]
			checkdfs(x - 1, y, 4)
			break
		case 4:
			ansPath = ansPath[:len(ansPath) - 1]
			return
		}
	}
}

/*
 *	here is the indication of water flow direction
 *	----------------------------------------------
 *				2
 *		 1 ->		<- 3
 *				4
 */
func dfs(x int, y int, waterEntry int){
	//fmt.Println(x, y, waterEntry)

	//check end
	if x == 4 && y == 5{
		flag = true
		return
	}

	//edge checking
	if (x < 0 || x > 4) || (y < 0 || y > 4) {
		return
	}
	//check whether the path is marked
	if markedMap[x][y] == 1{
		return
	}

	//fmt.Println(x, y, bitMap[x][y],waterEntry)
	//set the cur path to marked
	markedMap[x][y] = 1
	//push into path
	pathStack = append(pathStack, Loc{x, y, bitMap[x][y]})

	//tackle the straight tunel
	if bitMap[x][y] == 0 || bitMap[x][y] == 1{
		//entry at left, move to right
		if waterEntry == 1{
			dfs(x, y + 1, 1)
		//entry at top, move down
		}else if waterEntry == 2{
			dfs(x + 1, y, 2)
		//entry at right, move right
		}else if waterEntry == 3{
			dfs(x, y - 1, 3)
		}else{
			dfs(x - 1, y, 4)
		}
	}

	//tackle the conner
	if bitMap[x][y] >= 2 && bitMap[x][y] <= 5{
		//entry at left, move to top or down
		if waterEntry == 1{
			dfs(x + 1, y, 2)
			dfs(x - 1, y, 4)
			//entry at top, move left or right
		}else if waterEntry == 2{
			dfs(x, y - 1, 3)
			dfs(x, y + 1, 1)
			//entry at right, move top or bot
		}else if waterEntry == 3{
			dfs(x + 1, y, 2)
			dfs(x - 1, y, 4)
		}else{
			dfs(x, y - 1, 3)
			dfs(x, y + 1, 1)
		}
	}

	//finish to reset
	markedMap[x][y] = 0
	//pop the last
	pathStack = pathStack[:len(pathStack) - 1]
	return
}


func findSolution(){
	dfs(0, 0, 1)
}