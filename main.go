package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Taxes struct {
	FirstName   string `json:"first_Name"`
	LastName    string
	Email       string
	PhoneNumber string
	PaymentDue  int
}

type OverDue struct {
	FilePath string
	Data     []string
}

func main() {

	urgent := OverDue{FilePath: `\\qumulo\BLIS\ScottGetzFiles\NeverRoad_Processing\scott_urgent.txt`}
	nonUrgent := OverDue{FilePath: `\\qumulo\BLIS\ScottGetzFiles\NeverRoad_Processing\scott_nonurgent.txt`}
	superUrgent := OverDue{FilePath: `\\qumulo\BLIS\ScottGetzFiles\NeverRoad_Processing\scott_superurgent.txt`}
	superSuperUrgent := OverDue{FilePath: `\\qumulo\BLIS\ScottGetzFiles\NeverRoad_Processing\scott_supersuperurgent.txt`}

	miTekjson, err := filepath.Glob(`\\qumulo\BLIS\ScottGetzFiles\Mi-Tek\*.json`)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, mt := range miTekjson {
		content, err := os.ReadFile(mt)
		if err != nil {
			fmt.Println(err)
			continue
		}

		structTax := []Taxes{}

		err = json.Unmarshal(content, &structTax)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, st := range structTax {
			if st.PaymentDue < 2000 {
				continue
			}

			newPhoneNumber, err := generatePhone(st.PhoneNumber)
			if err != nil {
				fmt.Println(err)
				continue
			}

			overdueString := fmt.Sprintf("%v %v|%v|%s|%v", st.FirstName, st.LastName, st.Email, newPhoneNumber, fmt.Sprintf("Hello %s %s. You owe Cabarrus County $%d in tax fees", st.FirstName, st.LastName, st.PaymentDue))

			if st.PaymentDue > 20000 {
				superSuperUrgent.Data = append(superSuperUrgent.Data, overdueString)
				continue
			}

			if st.PaymentDue > 8000 {
				superUrgent.Data = append(superUrgent.Data, overdueString)
				continue
			}

			if st.PaymentDue > 4000 {
				urgent.Data = append(urgent.Data, overdueString)
				continue
			}

			nonUrgent.Data = append(nonUrgent.Data, overdueString)
		}
	}

	writeDataToFile(urgent, nonUrgent, superUrgent)
}

func generatePhone(s string) (string, error) {
	if !strings.Contains(s, "-") {
		return "", fmt.Errorf("%s Does not contain any dashes", s)
	}

	pSplit := strings.Split(s, "-")

	if len(pSplit) != 3 {
		return "", fmt.Errorf("%s Is not valid", s)
	}

	if (len(pSplit[0]) != 3) || (len(pSplit[1]) != 3) || (len(pSplit[2]) != 4) {
		return "", fmt.Errorf("%s Is not valid", s)
	}

	_, err := strconv.Atoi(pSplit[0])
	if err != nil {
		return "", err
	}

	_, err = strconv.Atoi(pSplit[1])
	if err != nil {
		return "", err
	}

	_, err = strconv.Atoi(pSplit[2])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s)-%s-%s", pSplit[0], pSplit[1], pSplit[2]), nil
}

func writeDataToFile(o ...OverDue) {
	for _, v := range o {
		if len(v.Data) == 0 {
			continue
		}

		err := writeData(v.Data, v.FilePath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func writeData(s []string, fp string) error {
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	var dataString string

	for _, ds := range s {
		dataString = fmt.Sprintf("%s%s\n", dataString, ds)
	}
	_, err = file.WriteString(dataString)
	if err != nil {
		return err
	}
	return nil
}
