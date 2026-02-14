package autofix

import "log"

func changeFirstLetter() {
	log.Println("Hello World") // want "log message should start with lowercase letter"
}

func dontChangeFirstLetter() {
	log.Printf("valid string") // ok
}

func removeSpecialCharacters() {
	log.Println("warning: something went wrong...") // want "log message should not contain special characters"
}

func dontRemoveFormat() {
	log.Printf("warning: something went wrong: %v", "some val") // want "log message should not contain special characters"
}
