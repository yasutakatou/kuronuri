package main

import (
	"crypto/aes"
	"crypto/cipher"
	crt "crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type encDecData struct {
	Type int    `json:"Type"`
	From string `json:"From"`
	To   string `json:"To"`
}

var encDec = []encDecData{}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	batFile := ""
	batCommand := ""

	_noRun := flag.Bool("noRun", false, "[-noRun=no run command or script.]")
	_Dry := flag.Bool("dry", false, "[-dry=no run command and not create script.]")
	_dstFile := flag.String("dst", "", "[-dst=create file name]")
	_wrapper := flag.String("wrap", "busybox.exe", "[-wrap=wrapper command]")
	_wrapOpt := flag.String("opt", "bash", "[-wrap=wrapper command arg option]")

	flag.Parse()

	eFlag := 0
	stra := ""
	strb := ""

	for i := 0; i < flag.NArg(); i++ {
		eFlag, stra, strb = encDecString(flag.Arg(i))
		if eFlag != 0 {
			if len(strb) > 16 || len(strb) <= 3 {
				fmt.Println("invalid key size > 16 or <= 3: ", stra)
			} else {
				strb = addSpace(strb)
				addEncDecStruct(eFlag, stra, strb)
			}
		} else {
			if Exists(flag.Arg(i)) == true {
				batFile = flag.Arg(i)
			} else {
				batCommand += flag.Arg(i) + " "
			}
		}
	}

	if len(batFile) == 0 && flag.NArg() < 1 {
		fmt.Println("  ubat: unix command enabler on windows.")
		fmt.Println(" usage) ubat [args..] (bat command or script file name)*require ")
		fmt.Println("option) -dst  : output file name.")
		fmt.Println("option) -noRun: no Run and no Delete script.")
		fmt.Println("option) -dry  : no run command and not create script.")
		fmt.Println("option) -wrap : wrapper command.")
		fmt.Println("option) -opt  : wrapper command arg option.")
		os.Exit(1)
	}

	do(*_wrapper, *_wrapOpt, batFile, batCommand, *_dstFile, *_noRun, *_Dry)

	os.Exit(0)
}

func addEncDecStruct(eFlag int, stra, strb string) {
	strs := ""
	if len(stra) == 0 {
		return
	}

	for i := 0; i < len(stra); i++ {
		switch int(stra[i]) {
		case 59: // ; 59
			encDec = append(encDec, encDecData{Type: eFlag, From: strs, To: strb})
			strs = ""
		default:
			strs += string(stra[i])
		}
	}

	encDec = append(encDec, encDecData{Type: eFlag, From: strs, To: strb})
}

func addSpace(strs string) string {
	for i := 0; len(strs) < 16; i++ {
		strs += "0"
	}
	return strs
}

func getFilename(fileName string) string {
	strs := ""

	if len(fileName) == 0 {
		strs = RandStr(8) + ".bat"
	} else {
		strs = fileName
	}
	return strs
}

func encodeOrDecode(tmpStr string) string {
	for r := 0; r < len(encDec); r++ {
		if strings.Index(tmpStr, encDec[r].From) != -1 {
			encDecStr := splitWord(tmpStr, strings.Index(tmpStr, encDec[r].From)+len(encDec[r].From))
			if len(encDecStr) < 256 {
				edStr := switchEncDec(encDec[r].Type, encDecStr, encDec[r].To)
				tmpStr = tmpStr[:strings.Index(tmpStr, encDec[r].From)+len(encDec[r].From)] + edStr + tmpStr[strings.Index(tmpStr, encDec[r].From)+len(encDec[r].From)+len(encDecStr):]
			} else {
				fmt.Println("invalid target word size > 256: ", encDecStr)
			}
		}
	}

	return tmpStr
}

func splitWord(strs string, wCnt int) string {
	var tmpStr = ""

	for i := wCnt; i < len(strs); i++ {
		switch strs[i] {
		case 9:
			return tmpStr
		case 10:
			return tmpStr
		case 13:
			return tmpStr
		case 32:
			return tmpStr
		default:
			tmpStr += string(strs[i])
		}
	}
	return tmpStr
}

func do(wrapCommand, wrapOpt, batFile, batCommand, fileName string, noDelete, dryRun bool) {
	scripts := ""

	fileName = getFilename(fileName)

	if len(batFile) == 0 {
		runOrNoRun(wrapCommand, wrapOpt, fileName, batCommand, noDelete, dryRun)

		return
	}

	scriptDataTmp := readFileToString(batFile)
	masterEnterCode := detectReturnCode(scriptDataTmp)
	scriptData := stringToArray(scriptDataTmp)

	for i := 0; i < len(scriptData); i++ {
		tmpStr := scriptData[i]
		tmpStr = encodeOrDecode(tmpStr)
		scripts += tmpStr + masterEnterCode
	}

	runOrNoRun(wrapCommand, wrapOpt, fileName, scripts, noDelete, dryRun)
}

func runOrNoRun(wrapCommand, wrapOpt, fileName, scripts string, noDelete, dryRun bool) {
	if dryRun == false {
		writeFile(fileName, scripts)
		if noDelete == false {
			fmt.Println(Execmd(wrapCommand, wrapOpt, fileName))
			tmpDelete(noDelete, fileName)
			return
		}
	}
	fmt.Printf(" - - - %s - - - \n", fileName)
	fmt.Println(scripts)
	fmt.Println(" - - - - - - - ")
}

func switchEncDec(encDec int, strs, key string) string {
	edStr := ""
	var err error

	if encDec == 1 {
		edStr, err = encrypt(strs, []byte(key))
		if err != nil {
			fmt.Sprintf("unable to encrypt the data: %v\n", err)
			edStr = "[Can't Encode]"
		}
	} else {
		edStr, err = decrypt(strs, []byte(key))
		if err != nil {
			fmt.Sprintf("unable to encrypt the data: %v\n", err)
			edStr = "[Can't Decode]"
		}

	}
	return edStr
}

func tmpDelete(noDelete bool, fileName string) {
	if noDelete == false {
		if err := os.Remove(fileName); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("ouput filename: ", fileName)
	}
}

func readFileToString(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func stringToArray(str string) []string {
	var strs []string

	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(str, -1) {
		strs = append(strs, v)
	}
	return strs
}

func readFileToStringArray(filename string) []string {
	var strs []string

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(bytes), -1) {
		strs = append(strs, v)
	}
	return strs
}

func writeFile(filename, strs string) bool {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(strs)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func RandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func Execmd(wrapCommand, wrapOpt, command string) string {
	var cmd *exec.Cmd
	var out string
	var err error

	cmd = exec.Command(wrapCommand, wrapOpt, command)
	outs, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	out = string(outs)

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest([]byte(out))
	if err == nil {
		if result.Charset == "Shift_JIS" {
			out, _ = sjis_to_utf8(out)
		}
	}

	return string(out)
}

//FYI: https://qiita.com/uchiko/items/1810ddacd23fd4d3c934
// ShiftJIS から UTF-8
func sjis_to_utf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

// 0: False 1: Encode 2: Decode
// http://www9.plala.or.jp/sgwr-t/c_sub/ascii.html
func encDecString(strs string) (int, string, string) {
	eFlag := 0

	if len(strs) > 5 { //40: ( 41: )
		if strs[0] == 40 && strs[len(strs)-1] == 41 {
			eFlag = 1
		} else if strs[0] == 41 && strs[len(strs)-1] == 40 {
			eFlag = 2
		} else {
			return eFlag, "", ""
		}

		if strings.Index(strs, ":") != -1 {
			out := strings.Split(strs, ":")
			return eFlag, out[0][1:], out[1][:len(out[1])-1]
		} else {
			eFlag = 0
		}
	}
	return eFlag, "", ""
}

func detectReturnCode(strs string) string {
	r := regexp.MustCompile("\r\n")
	if r.MatchString(strs) == true {
		return "\r\n"
	}

	r = regexp.MustCompile("\n\r")
	if r.MatchString(strs) == true {
		return "\n\r"
	}

	r = regexp.MustCompile("\n")
	if r.MatchString(strs) == true {
		return "\n"
	}

	return "\r"
}

// FYI: http://www.inanzzz.com/index.php/post/f3pe/data-encryption-and-decryption-with-a-secret-key-in-golang
// encrypt encrypts plain string with a secret key and returns encrypt string.
func encrypt(plainData string, secret []byte) (string, error) {
	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(crt.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(plainData), nil)), nil
}

// decrypt decrypts encrypt string with a secret key and returns plain string.
func decrypt(encodedData string, secret []byte) (string, error) {
	encryptData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return "", err
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainData), nil
}
