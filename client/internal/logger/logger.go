package logger

import (
	"log"

	"github.com/fatih/color"
)

func Info(str string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	log.Printf("%s %s", yellow("[info]"), str)
}
func Warn(str string) {
	red := color.New(color.FgRed).SprintFunc()
	log.Printf("%s %s", red("[warn]"), str)
}
func Norm(str string) {
	green := color.New(color.FgGreen).SprintFunc()
	log.Printf("%s %s", green("[norm]"), str)
}
