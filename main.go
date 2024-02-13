package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type Item struct {
	Name string
	Date time.Time
	Tags []string
	Link string
}


func main() {
	defer func() {
		// Завершаем работу с клавиатурой при выходе из функции
		_ = keyboard.Close()
	}()
	listItem := []*Item{}
	fmt.Println("Программа для добавления url в список")
	fmt.Println("Для выхода и приложения нажмите Esc")
	help()

OuterLoop:
	for {
		// Подключаем отслеживание нажатия клавиш
		fmt.Println("Введите команду")
		if err := keyboard.Open(); err != nil {
			log.Fatal(err)
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch char {
		case 'a':
			if err := keyboard.Close(); err != nil {
				log.Fatal(err)
			}

			// Добавление нового url в список хранения
			fmt.Println("Введите новую запись в формате <url описание теги>")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			args := strings.Fields(text)
			if len(args) < 3 {
				fmt.Println("Введите правильный аргументы в формате url описание теги")
				continue OuterLoop
			}
			if err := chekTags(args[2:]); !err {
				fmt.Println("Теги должны начинатся со знака '#' и иметь после знача каку-либо информацию!")
				continue OuterLoop
			}
			adres := Item{
				Name: strings.ToLower(args[1]),
				Date: time.Now(),
				Tags: args[2:],
				Link: args[0],
			}

			listItem = append(listItem, &adres)
			continue OuterLoop
		case 'l':
			// Вывод списка добавленных url. Выведите количество добавленных url и список с данными url
			// Вывод в формате
			// Имя: <Описание>
			// URL: <url>
			// Теги: <Теги>
			// Дата: <дата>
			for _, v := range listItem {
				year, month, day := v.Date.Date()
				fmt.Printf("\nИмя: %s\nURL: %s\nТеги: ", v.Name, v.Link)
				for i := 0; i < len(v.Tags)-1; i++ {
					fmt.Printf("%s, ", v.Tags[i])
				}
				fmt.Printf("%s", v.Tags[len(v.Tags)-1])
				fmt.Printf("\nДата: %d %v %d\n\n", year, month, day)
			}
		case 'r':
			if err := keyboard.Close(); err != nil {
				log.Fatal(err)
			}
			// Удаление url из списка хранения
			fmt.Println("Введите имя ссылки, которое нужно удалить")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)//убрал \n в кенце строки
			text = strings.ToLower(text)
			for ind, v := range listItem {
				if v.Name == text {
					listItem[ind] = listItem[len(listItem)-1]
					listItem = listItem[:len(listItem)-1]
					fmt.Println("Запись успешно удалена")
					continue OuterLoop
				}
			}
			fmt.Println("С таким названием записи нет")
		default:
			// Если нажата Esc выходим из приложения
			if key == keyboard.KeyEsc {
				return
			}
		}
	}
}

func help() {
	fmt.Println("\nкомнады для выполнения операций:\n" +
		"a - Добавление нового url в список хранения\n" +
		"l - Просмотр существующих url\n" +
		"r - Удаление url из списка хранения")
}

func chekTags(tags []string) bool{
	for _, v := range tags {
		if (v[0] != '#' || len(v)<2){
			return false
		}
	}
	return true
}