package common

import "errors"

func CheckGram(gramStr string) error {
	for i := 0; i < len(gramStr); i++ {
		if string(gramStr[i]) == "." {
			if len(gramStr)-i+1 <= 3 {
				return nil
			}
			return errors.New("gram is not multiplication of 0.001")
		}
	}

	return nil
}
