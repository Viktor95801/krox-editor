package title_gen

import (
	"math/rand"
	"strings"
)

const (
	titleEndings = "Also try 'vscode'!,Also try 'lite-xl'!,Did you mean 'krox-editor'?,How can I eat cheese without a knife?,I'm a teapot,We're a teapot,Welcome to the Matrix,What's the difference between a dog and a cat?,Quick brown fox jumps over the lazy dog,How to center a div?,A mango yeets a orange (?!?!),Manuals are a teapot,Mans are a teapot,I love cheese,I like cats,Spaghetti is a teapot,Spaghetti <3,I kinda like spaghetti tbh,programming is a teapot,unix is a teapot,posix is a teapot,windows is not a teapot,minix 3 is not safe from being a teapot,universe is a teapot (except for windows),how to make a teapot,how to make a teapot out of a teapot,how to make a teapot out of a teapot out of a teapot,i promise no more teapots,i promise no more teapots out of teapots,i promise no more teapots out of teapots out of teapots"
)

func GetRandomTitleEnding() string {
	endings := strings.Split(titleEndings, ",")
	return endings[rand.Intn(len(endings))]
}

func GetRandomTitle() string {
	return "EGG: " + GetRandomTitleEnding()
}
