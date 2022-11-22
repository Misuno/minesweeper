package main

import (
	"bufio"
	"fmt"
	"minesweeper/minesweeper"
	consoleclient "minesweeper/minesweeper/console_client"
	"os"
	"os/exec"
)

const (
	vSize      int = 5
	hSize      int = 5
	bombsCount int = 2
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	//init()
	CallClear()
	board := minesweeper.MakeBoard(hSize, vSize, bombsCount)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(
			`Управление:
Команды вводятся латиницей в виде <команда> <x> <y> (без уголков)
Примеры:
o 1 2 откроет поле под номером (1, 2)
f 2 3 поставил флажок на поле (2, 3)`)
		consoleclient.DrawBoard(board)
		fmt.Print("Что делаем? (o - open, f - flag) ->")

		res, err := consoleclient.ProcessInput(reader, board)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if res == minesweeper.StateBomb {
			fmt.Println("BOOOM!")
			consoleclient.DrawBoardOpen(board)
			return
		}

		if board.IsAllOpened() {
			fmt.Println("YOU ARE WINNER!")
			return
		}
	}
}
