package errorhandler

import "log"

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
func Basic(err error) {
	if err != nil {
		log.Println("an error occured:", err)
	}
}
