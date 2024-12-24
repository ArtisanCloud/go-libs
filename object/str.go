package object

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	/**
	 * The cache of snake-cased words.
	 *
	 * @var array
	 */
	SnakeCache StringMap

	/**
	 * The cache of camel-cased words.
	 *
	 * @var array
	 */
	CamelCache StringMap

	/**
	 * The cache of studly-cased words.
	 *
	 * @var array
	 */
	StudlyCache StringMap

	/**
	 * The callback that should be used to generate UUIDs.
	 *
	 * @var callable
	 */
	UidFactory StringMap
)

func init() {

	SnakeCache = StringMap{}
	CamelCache = StringMap{}
	StudlyCache = StringMap{}
	UidFactory = StringMap{}
}

func UniqueID(prefix string) string {
	now := time.Now()
	sec := now.Unix()
	usec := now.UnixNano() % 0x100000
	return fmt.Sprintf("%s%08x%05x", prefix, sec, usec)
}

/**
 * Convert a value to camel case.
 *
 * @param  string  value
 * @return string
 */
func Camel(value string) string {

	// if cache has converted
	if _, ok := CamelCache[value]; ok {
		return CamelCache[value]
	}

	// low case first char and store into cache
	CamelCache[value] = LCFirst(Studly(value))

	return CamelCache[value]
}

/**
 * Convert a string to snake case.
 *
 * @param  string  $value
 * @param  string  $delimiter
 * @return string
 */
func Snake(value string, delimiter string) string {
	if delimiter == "" {
		delimiter = "_"
	}

	key := value + delimiter
	// if cache has converted
	if _, ok := CamelCache[key]; ok {
		return CamelCache[key]
	}

	if !IsUpper(value) {
		value = RegexpReplace("/\\s+/u", "", UCWords(value))
		//value = Lower(RegexpReplace("/(.)(?=[A-Z])/g", "$1"+delimiter, value))
		value = Lower(strings.Replace(value, "-", delimiter, -1))
	}
	//
	CamelCache[key] = value
	return CamelCache[key]
}

/**
 * Convert a value to studly caps case.
 *
 * @param  string  value
 * @return string
 */
func Studly(value string) string {

	// if cache has converted
	if _, ok := StudlyCache[value]; ok {
		return StudlyCache[value]
	}

	// replace "-" or "_" with " "
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.ReplaceAll(value, "_", " ")

	// Up Case words
	value = UCWords(value)

	// replace " " with "", and store into cache
	StudlyCache[value] = strings.ReplaceAll(value, " ", "")

	return StudlyCache[value]

}

/**
 * Make a string's first character lowercase
 * @param string str <p>
 * The input string.
 * </p>
 * @return string the resulting string.
 */
func LCFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

/**
 * Uppercase the first character of each word in a string
 * @param string str <p>
 * The input string.
 * </p>
 * @param string delimiters [optional] <p>
 * @return string the modified string.
 */
func UCWords(str string) string {
	return strings.Title(str)
}

/**
 * Check for uppercase character(s)
 * @param string str <p>
 * The input string.
 * </p>
 * @param string<p>
 * @return bool.
 */
func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

/**
 * Check for lowercase character(s)
 * @param string str <p>
 * The input string.
 * </p>
 * @param string<p>
 * @return bool.
 */
func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

/**
 * Convert the given string to lower-case.
 *
 * @param  string  $value
 * @return string
 */
func Lower(value string) string {
	return strings.ToLower(value)
}

/**
 * Convert the given string to upper-case.
 *
 * @param  string  $value
 * @return string
 */
func Upper(value string) string {
	return strings.ToUpper(value)
}

/**
 * Replace by regex
 *
 * @param  string  $value
 * @return string
 */
func RegexpReplace(pattern string, replacement string, subject string) string {
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllString(subject, replacement)

}

func Implode(glue string, arrayStrs []string) string {
	var strValues string
	if len(arrayStrs) > 0 {
		for _, v := range arrayStrs {
			strValues += v + "|"
		}
		strValues = strValues[:len(strValues)-1]
	}
	return strValues
}

func Shuffle(str string) string {

	randBias := rand.Int63n(100)
	timestamp := time.Now().Unix()
	randBias = timestamp + randBias
	//fmt2.Dump(randBias)
	rand.Seed(randBias)

	inRune := []rune(str)
	rand.Shuffle(len(inRune), func(i, j int) {
		//fmt.Printf("i:%d, j:%d ", i, j)
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func QuickRandom(length int) string {
	pool := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// repeat the string
	strRepeat := strings.Repeat(pool, length)

	// shuffle the string
	strShuffle := Shuffle(strRepeat)

	return strShuffle[0:length]
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMask(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func ContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
