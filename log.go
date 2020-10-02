package main

import (
	"fmt"
	"github.com/TwinProduction/go-color"
)

func logInfo(reference, data string) {
	fmt.Println(color.Ize(color.Green, "["+reference+"] --INF-- "+data))
}

func logError(reference, data string) {
	fmt.Println(color.Ize(color.Red, "["+reference+"] --ERR-- "+data))
}

func logWarning(reference, data string) {
	fmt.Println(color.Ize(color.Yellow, "["+reference+"] --WAR-- "+data))
}
