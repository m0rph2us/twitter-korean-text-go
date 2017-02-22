package util

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strconv"
	"strings"
	"tktg/shared"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var resourcePath = os.Getenv("KRGO_DIC_RSRC")

var KoreanEntityFreq map[string]float64

var KoreanDictionary data_structure.Map

var SpamNouns data_structure.Set

var TypoDictionaryByLength map[int]map[string]string

var ProperNouns data_structure.Set

var nameDictionary map[string]data_structure.Set

var PredicateStems map[KoreanPos]map[string]string

func init() {
	KoreanEntityFreq = readWordFreqs(resourcePath + "freq/entity-freq.txt.gz")
	KoreanDictionary = getKoreanDictionary()
	SpamNouns = readWords([]string{resourcePath + "noun/spam.txt", resourcePath + "noun/profane.txt"})
	TypoDictionaryByLength = getTypoDictionaryByLength()
	ProperNouns = readWords(
		[]string{
			resourcePath + "noun/entities.txt",
			resourcePath + "noun/names.txt",
			resourcePath + "noun/twitter.txt",
			resourcePath + "noun/lol.txt",
			resourcePath + "noun/company_names.txt",
			resourcePath + "noun/foreign.txt",
			resourcePath + "noun/geolocations.txt",
			resourcePath + "substantives/given_names.txt",
			resourcePath + "noun/kpop.txt",
			resourcePath + "noun/bible.txt",
			resourcePath + "noun/pokemon.txt",
			resourcePath + "noun/congress.txt",
			resourcePath + "noun/wikipedia_title_nouns.txt",
		})
	nameDictionary = map[string]data_structure.Set{
		"family_name": readWords([]string{resourcePath + "substantives/family_names.txt"}),
		"given_name":  readWords([]string{resourcePath + "substantives/given_names.txt"}),
		"full_name": readWords([]string{
			resourcePath + "noun/kpop.txt",
			resourcePath + "noun/foreign.txt",
			resourcePath + "noun/names.txt",
		}),
	}
	PredicateStems = getPredicateStems()
}

func getPredicateStems() map[KoreanPos]map[string]string {
	getConjugationMap := func(words data_structure.Set, isAdjective bool) map[string]string {
		v := []interface{}{}

		for word := range words.Iter() {
			s := conjugatePredicated(data_structure.NewLinkedSet(word), isAdjective)

			for x := range s.Iter() {
				v = append(v, []string{x.(string), word.(string) + "다"})
			}
		}

		return scala.ToMap(v).(map[string]string)
	}

	return map[KoreanPos]map[string]string{
		Verb:      getConjugationMap(readWordsAsSet([]string{resourcePath + "verb/verb.txt"}), false),
		Adjective: getConjugationMap(readWordsAsSet([]string{resourcePath + "adjective/adjective.txt"}), true),
	}
}

func convertTypoDictionaryByLength(value map[interface{}]map[string]string) map[int]map[string]string {
	m := map[int]map[string]string{}
	for k, v := range value {
		m[k.(int)] = v
	}
	return m
}

func getTypoDictionaryByLength() map[int]map[string]string {
	x := scala.GroupBy(readWordMap(resourcePath+"typos/typos.txt"),
		func(key interface{}, value interface{}) interface{} {
			return utf8.RuneCountInString(key.(string))
		},
	).(map[interface{}]map[string]string)

	return convertTypoDictionaryByLength(x)
}

func getKoreanDictionary() data_structure.Map /*map[KoreanPos]data_structure.Set*/ {
	m := data_structure.NewLinkedMap()

	m.Put(Noun, readWords([]string{
		resourcePath + "noun/nouns.txt",
		resourcePath + "noun/entities.txt",
		resourcePath + "noun/spam.txt",
		resourcePath + "noun/names.txt",
		resourcePath + "noun/twitter.txt",
		resourcePath + "noun/lol.txt",
		resourcePath + "noun/slangs.txt",
		resourcePath + "noun/company_names.txt",
		resourcePath + "noun/foreign.txt",
		resourcePath + "noun/geolocations.txt",
		resourcePath + "noun/profane.txt",
		resourcePath + "substantives/given_names.txt",
		resourcePath + "noun/kpop.txt",
		resourcePath + "noun/bible.txt",
		resourcePath + "noun/pokemon.txt",
		resourcePath + "noun/congress.txt",
		resourcePath + "noun/wikipedia_title_nouns.txt",
	}))

	m.Put(Verb, conjugatePredicatesToSet(readWordsAsSet([]string{resourcePath + "verb/verb.txt"}), false))
	m.Put(Adjective, conjugatePredicatesToSet(readWordsAsSet([]string{resourcePath + "adjective/adjective.txt"}), true))
	m.Put(Adverb, readWordsAsSet([]string{resourcePath + "adverb/adverb.txt"}))
	m.Put(Determiner, readWordsAsSet([]string{resourcePath + "auxiliary/determiner.txt"}))
	m.Put(Exclamation, readWordsAsSet([]string{resourcePath + "auxiliary/exclamation.txt"}))
	m.Put(Josa, readWordsAsSet([]string{resourcePath + "josa/josa.txt"}))
	m.Put(Eomi, readWordsAsSet([]string{resourcePath + "verb/eomi.txt"}))
	m.Put(PreEomi, readWordsAsSet([]string{resourcePath + "verb/pre_eomi.txt"}))
	m.Put(Conjunction, readWordsAsSet([]string{resourcePath + "auxiliary/conjunctions.txt"}))
	m.Put(NounPrefix, readWordsAsSet([]string{resourcePath + "substantives/noun_prefix.txt"}))
	m.Put(VerbPrefix, readWordsAsSet([]string{resourcePath + "verb/verb_prefix.txt"}))
	m.Put(Suffix, readWordsAsSet([]string{resourcePath + "substantives/suffix.txt"}))

	return m
}

func readStreamByLine(reader io.Reader) []string {
	scanner := bufio.NewScanner(reader)

	lineSlice := []string{}
	for scanner.Scan() {
		lineSlice = append(lineSlice, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	x := scala.Map(lineSlice, func(key interface{}, value interface{}) interface{} {
		return strings.TrimSpace(value.(string))
	}).([]interface{})

	y := scala.Filter(x, func(value interface{}) bool {
		return len(value.(string)) > 0
	}).([]interface{})

	return shared.ConvertSliceOfInterfaceToSliceOfString(y)
}

func readWordFreqs(filename string) map[string]float64 {
	freqMap := map[string]float64{}

	for _, v := range readFileByLineFromResources(filename) {
		if strings.Contains(v, "\t") {
			sp := strings.Split(v, "\t")
			if f, err := strconv.ParseFloat(string(sp[1][0:6]), 64); err == nil {
				freqMap[sp[0]] = f
			} else {
				panic(err.Error())
			}
		}
	}

	return freqMap
}

func readWordMap(filename string) map[string]string {
	d := readFileByLineFromResources(filename)
	x := scala.Filter(d, func(value interface{}) bool {
		return strings.Contains(value.(string), " ")
	}).([]interface{})

	y := scala.Map(x, func(key interface{}, value interface{}) interface{} {
		v := strings.Split(value.(string), " ")
		return []string{v[0], v[1]}
	}).([]interface{})

	return scala.ToMap(y).(map[string]string)
}

func readWordAsSlice(filename string) []string {
	return readFileByLineFromResources(filename)
}

func readWordsAsSet(filenames []string) data_structure.Set {
	s := scala.FoldLeft(filenames, data_structure.NewLinkedSet(), func(folded interface{}, key interface{}, value interface{}) interface{} {
		for _, v := range readFileByLineFromResources(value.(string)) {
			folded.(data_structure.Set).Add(v)
		}
		return folded
	}).(interface{})

	return s.(data_structure.Set)
}

func readWords(filenames []string) data_structure.Set {
	s := data_structure.NewLinkedSet()

	for _, v := range filenames {
		for _, x := range readFileByLineFromResources(v) {
			s.Add(x)
		}
	}

	return s
}

func getStream(filename string) io.Reader {
	if file, err := os.Open(filename); err == nil {
		if strings.HasSuffix(filename, ".gz") {
			if r, e := gzip.NewReader(file); e == nil {
				return r
			} else {
				panic(e.Error())
			}
		} else {
			return file
		}
	} else {
		panic(err.Error())
	}
}

func readFileByLineFromResources(filename string) []string {
	return readStreamByLine(getStream(filename))
}

func AddWordsToDictionary(pos KoreanPos, words []string) {
	for _, v := range words {
		KoreanDictionary.Get(pos).(data_structure.Set).Add(v)
	}
}
