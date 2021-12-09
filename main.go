package smspdu

import (
	"math"
	"strconv"
	"unicode/utf8"
)

func decimalToAny(num int, n int, minLength int) string {
	num2char := "0123456789ABCDEF"
	newNumStr := ""
	var remainder int
	var remainderString string
	for num != 0 {
		remainder = num % n
		remainderString = string(num2char[remainder])
		newNumStr = remainderString + newNumStr //注意顺序
		num = num / n
	}
	if minLength != 0 {
		length := len(newNumStr)
		if length < minLength { //如果小于8位
			for i := 1; i <= minLength-length; i++ {
				newNumStr = "0" + newNumStr
			}
		}
	}

	return newNumStr
}

func toHex(i int) string {
	var sHex = "0123456789ABCDEF"
	var Out = ""
	Out = string(sHex[i&0xf])
	i >>= 4
	Out = string(sHex[i&0xf]) + Out
	return Out
}

func semiOctetToString(str string) string {
	text := []rune(str)
	length := len(text)
	out := ""
	for i := 0; i < length; i = i + 2 {
		out = out + string(text[i+1]) + string(text[i])
	}
	return out
}

func binToInt(x string) int {
	total := 0
	length := len(x)
	var power = length - 1
	for i := 0; i < length; i++ {
		if string(x[i]) == "1" {
			total = total + int(math.Pow(2, float64(power)))
		}
		power--
	}
	return total
}

func intToHex(i int) string {
	sHex := "0123456789ABCDEF"
	h := ""
	for j := 0; j <= 3; j++ {
		h += string(sHex[(i>>(j*8+4))&0x0F]) + string(sHex[(i>>(j*8))&0x0F])
	}
	return h[0:2]
}

func getSevenBit(r rune) int {
	sevenbitdefault := []rune{'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å', '\u0394', '_', '\u03a6', '\u0393', '\u039b', '\u03a9', '\u03a0', '\u03a8', '\u03a3', '\u0398', '\u039e', '€', 'Æ', 'æ', 'ß', 'É', ' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?', '¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§', '¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à'}
	for _, c := range sevenbitdefault {
		if c == r {
			return 1
		}
	}
	return 0
}

// StringToPDU /**
func StringToPDU(smsText, phoneNumber, smscNumber string, size, mclass, valid int, receipt bool) string {

	octetFirst := ""
	octetSecond := ""
	output := ""
	//Make header
	SmscInfoLengthNum := 0
	SmscInfoLength := "00"
	SmscLength := 0
	SmscNumberFormat := ""
	Smsc := ""

	if smscNumber != "" {
		SmscNumberFormat = "81" // national
		if smscNumber[0:1] == "+" {
			SmscNumberFormat = "91" // international
			smscNumber = smscNumber[1:]
		} else if smscNumber[0:1] != "0" {
			SmscNumberFormat = "91" // international
		}
		if len(smscNumber)%2 != 0 {
			// add trailing F
			smscNumber += "F"
		}

		Smsc = semiOctetToString(smscNumber)
		SmscInfoLengthNum = len(SmscNumberFormat+""+Smsc) / 2
		SmscLength = SmscInfoLengthNum
	}

	if SmscInfoLengthNum < 10 {
		SmscInfoLength = "0" + strconv.Itoa(SmscInfoLengthNum)
	}
	firstOctet := "1100"
	if receipt {
		if valid != 0 {
			firstOctet = "3100" // 18 is mask for validity period // 10 indicates relative
		} else {
			firstOctet = "2100"
		}
	} else {
		if valid != 0 {
			firstOctet = "1100"
		} else {
			firstOctet = "0100"
		}
	}

	ReiverNumberFormat := "81" // national
	if phoneNumber[0:1] == "+" {
		ReiverNumberFormat = "91"     // international
		phoneNumber = phoneNumber[1:] //,phoneNumber.length-1);
	} else if phoneNumber[0:1] != "0" {
		ReiverNumberFormat = "91" // international
	}
	ReiverNumberLength := decimalToAny(len(phoneNumber), 16, 2)
	if len(phoneNumber)%2 != 0 {
		// add trailing F
		phoneNumber += "F"
	}
	ReiverNumber := semiOctetToString(phoneNumber)
	ProtoId := "00"
	DCS := 0
	if mclass != -1 { // AJA
		DCS = mclass | 0x10
	}
	if size == 8 {
		DCS = DCS | 4
	} else if size == 16 {
		DCS = DCS | 8
	}
	DataEncoding := decimalToAny(DCS, 16, 2)
	ValidPeriod := "0"
	if valid != 0 {
		ValidPeriod = decimalToAny(valid, 16, 2) // AA
	}

	smsTextRune := []rune(smsText)
	smsTextRuneLength := len(smsTextRune)
	userDataSize := "00"

	if size == 7 {
		userDataSize = decimalToAny(smsTextRuneLength, 16, 0)

		for i := 0; i <= smsTextRuneLength; i++ {
			if i == smsTextRuneLength {
				if octetSecond != "" { // AJA Fix overshoot
					output = output + "" + intToHex(binToInt(octetSecond))
				}
				break
			}

			current := decimalToAny(getSevenBit(smsTextRune[i]), 2, 7)

			currentOctet := ""
			if i != 0 && i%8 != 0 {
				octetFirst = current[7-(i)%8:]
				currentOctet = octetFirst + octetSecond //put octet parts together

				output = output + "" + (intToHex(binToInt(currentOctet)))
				octetSecond = current[0 : 7-(i)%8] //set net second octet
			} else {
				octetSecond = current[0 : 7-(i)%8]
			}
		}

	} else if size == 8 {
		userDataSize = decimalToAny(smsTextRuneLength, 16, 0)

		var CurrentByte = 0
		for i := 0; i < smsTextRuneLength; i++ {
			char, _ := utf8.DecodeRuneInString(string(smsTextRune[i]))
			CurrentByte = int(char)
			output = output + "" + toHex(CurrentByte)
		}
	} else if size == 16 {
		userDataSize = decimalToAny(smsTextRuneLength*2, 16, 0)

		myChar := 0
		for i := 0; i < smsTextRuneLength; i++ {
			char, _ := utf8.DecodeRuneInString(string(smsTextRune[i]))
			myChar = int(char)
			output = output + toHex((myChar&0xff00)>>8) + toHex(myChar&0xff)
		}
	}

	header := SmscInfoLength + SmscNumberFormat + Smsc + firstOctet + ReiverNumberLength + ReiverNumberFormat + ReiverNumber + ProtoId + DataEncoding + ValidPeriod + userDataSize
	PDU := header + output
	AT := "AT+CMGS=" + strconv.Itoa(len(PDU)/2-SmscLength-1)

	//CMGW
	return AT + "\n" + PDU
}
