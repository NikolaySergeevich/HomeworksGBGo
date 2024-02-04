package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	file, err := os.Open("text.txt")
	if err != nil {
		fmt.Println("Не смог открыть файл:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	amountLet := getAmountLetStr(&text)

	mapa := make(map[rune]map[int]float64)
	
	for _, v := range text {
		if unicode.IsLetter(v){
			if inMap, exist := mapa[v]; exist{
				for kye := range mapa[v] {
					delete(inMap, kye)
					mapa[v][kye+1] = float64(kye+1)/float64(amountLet)
				}
			}else{
				mapa[v] = make(map[int]float64)
				mapa[v][1] = float64(1)/float64(amountLet)
			}
		}		
	}
	printMap(mapa, &text)

}
// вывод мапы согласно условиям
func printMap(mapa map[rune]map[int]float64, textc *string) {
	help := make(map[rune]string)//для невывода повторяющихся букв
	for _, v := range *textc {
		_, ok := help[v]
		if unicode.IsLetter(v) && !ok{
			for kye, val := range mapa[v] {
				// fmt.Println(string(v))
				fmt.Printf("%c - %d %0.2f\n", v, kye, val)
			}
			help[v] = ""
		}
		
	}
}
// возвращает кол-во букв в строке
func getAmountLetStr (str *string) (res int){
	for _, v := range *str {
		if unicode.IsLetter(v) {
			res++
		}
	}
	return
}
