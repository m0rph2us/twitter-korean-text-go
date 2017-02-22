package util

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

func TestReadLineByLine(t *testing.T) {
	if file, err := os.Open("hangul.go"); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		lineSlice := []string{}
		for scanner.Scan() {
			lineSlice = append(lineSlice, scanner.Text())
		}

		if err = scanner.Err(); err != nil {
			panic(err.Error())
		}

		if len(lineSlice) == 0 {
			t.Errorf("Expected result should match.")
		}
	} else {
		panic(err.Error())
	}
}

func TestReadLineByLineMapFilter(t *testing.T) {
	if file, err := os.Open("hangul.go"); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		lineSlice := []string{}
		for scanner.Scan() {
			lineSlice = append(lineSlice, scanner.Text())
		}

		if err = scanner.Err(); err != nil {
			panic(err.Error())
		}

		x := scala.Map(lineSlice, func(key interface{}, value interface{}) interface{} {
			return strings.TrimSpace(value.(string))
		}).([]interface{})

		y := scala.Filter(x, func(p interface{}) bool {
			return len(p.(string)) > 0
		}).([]interface{})

		if len(y) == 0 {
			t.Errorf("Expected result should match.")
		}
	} else {
		panic(err.Error())
	}
}

func TestReadStreamByLine(t *testing.T) {
	if file, err := os.Open("hangul.go"); err == nil {
		defer file.Close()

		v := readStreamByLine(file)

		if len(v) == 0 {
			t.Errorf("Expected result should match.")
		}
	} else {
		panic(err.Error())
	}
}

func TestReadWordMap(t *testing.T) {
	readWordMap("hangul.go")
}

func TestReadWordFreqs(t *testing.T) {
	v := readWordFreqs("../resources/freq/entity-freq.txt.gz")

	if len(v) == 0 {
		t.Errorf("Expected result should match.")
	}
}

func TestAddWordsToDictionary(t *testing.T) {
	nonExistentWord := "없는명사다"
	if KoreanDictionary.Get(Noun).(data_structure.Set).Contains(nonExistentWord) {
		t.Error("Expected result should match.")
	}

	AddWordsToDictionary(Noun, []string{nonExistentWord})

	if !KoreanDictionary.Get(Noun).(data_structure.Set).Contains(nonExistentWord) {
		t.Error("Expected result should match.")
	}
}
