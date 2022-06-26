package wisdom

import (
	"math/rand"
	"time"
)

func Wisdom() string {
	WordsOfWisdom := []string{
		"There are no excuses. There are no good reasons",
		"More than cleverness, we need kindness and gentleness. Charlie Chaplin kindness",
		"Do not waste the precious gift of your life on grudges, petty differences, or even big offenses.",
		"Find yourself and be yourself: Remember, there is no one else on earth like you.",
		"My words are powerful. I use them to uplift and inspire. I speak blessings, favor, and peace.",
		"Education is something we have to keep pursuing day after day.",
		"Being a leader is about making decisions that put you in the position to influence and inspire.",
	}
	rand.Seed(time.Now().Unix())
	return WordsOfWisdom[rand.Intn(len(WordsOfWisdom))]
}
