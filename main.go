package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text    string
	Options []string
	Answer  int
}

type GameState struct {
	Name      string
	Points    int
	Questions []Question
}

var Blue = "\033[34m"
var Gray = "\033[37m"
var White = "\033[97m"
var Green = "\033[32m"
var Red = "\033[31m"
var Yellow = "\033[33m"

func (g *GameState) Init() {
	fmt.Println("👾 Seja bem-vindo ao quizDEV.👾 ")
	fmt.Println("Digite seu nome:")
	reader := bufio.NewReader(os.Stdin)

	name, _ := reader.ReadString('\n')
	g.Name = strings.TrimSpace(name)

	fmt.Printf("Vamos ao jogo %s\n", g.Name)
}

func (g *GameState) ProcessCSV() {
	fmt.Println("Escolha um tema para iniciar o quiz:")
	fmt.Println("1. Linguagem de programação Golang")
	fmt.Println("2. Linguagem de programação JavaScript")

	var tema int
	fmt.Scanln(&tema)

	var arquivo string

	switch tema {
	case 1:
		arquivo = "./Files/quizgo.csv"
	case 2:
		arquivo = "./Files/quizjs.csv"
	default:
		fmt.Println("Tema inválido. Por favor escolha um tema valido.")
		return
	}

	f, err := os.Open(arquivo)

	if err != nil {
		panic("Erro ao abrir o arquivo")
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		panic("Erro ao ler o CSV")
	}

	for index, record := range records {
		//fmt.Println(index, record)

		if index > 0 {
			correctAnswer, _ := toInt(record[5])
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}
			g.Questions = append(g.Questions, question)
		}

	}

}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("não e permitido o uso de caracteres diferentes de números, por favor digite um número")
	}
	return i, nil
}

func (g *GameState) run() {
	//exibir pergunta pro usuário
	for index, question := range g.Questions {
		fmt.Printf(Green+"%d. %s\n"+Gray, index+1, question.Text)

		for j, option := range question.Options {
			fmt.Printf("[%d]. %s\n", j+1, option)
		}

		fmt.Println("Digite uma alternativa:")

		var answer int
		var err error

		//Tempo limite para responder....

		timeout := time.After(10 * time.Second)

		//cria um canal para receber a resposta do usuário.

		reponse := make(chan string)

		go func() {
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n')
			reponse <- strings.TrimSpace(read)
		}()

		select {
		case <-timeout:
			fmt.Println(Yellow + "Tempo esgotado. ⏰" + Yellow)
			return

		case resp := <-reponse:
			answer, err = toInt(resp)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if answer == question.Answer {
				fmt.Println("Parabéns, resposta correta ✔️")
				g.Points += 10
			} else {
				fmt.Println(Red + "Ops! Resposta errada 👻" + Red)
				fmt.Println("--------------------------------------")
			}

		}
	}
}

func main() {
	game := &GameState{Points: 0}
	game.ProcessCSV()
	game.Init()

	game.run()

	fmt.Printf(White+"Fim do jogo, voce fez %d pontos\n"+White, game.Points)
}
