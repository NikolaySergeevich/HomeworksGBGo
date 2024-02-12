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
	Tags string
	Link string
}

//замена "/" на ", " в тегах
func(t *Item) exchangeLet(){
	t.Tags = strings.ReplaceAll(t.Tags, "/", ", ")
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
			adres := Item{
				Name: args[1],
				Date: time.Now(),
				Tags: args[len(args)-1],
				Link: args[0],
			}
			adres.exchangeLet()
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
				fmt.Printf("\nИмя: %s\nURL: %s\nТеги: %v\nДата: %d %v %d\n\n", v.Name, v.Link, v.Tags, year, month, day)
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