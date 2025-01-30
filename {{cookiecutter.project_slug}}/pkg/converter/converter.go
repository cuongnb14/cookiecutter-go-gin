package converter

import (
	"log"

	"github.com/jinzhu/copier"
)

func MustConvert[T any](from any) *T {
	var t T
	err := copier.Copy(&t, from)
	if err != nil {
		log.Fatal(err)
	}
	return &t
}

func MustConvertList[F any, T any](from *[]F) *[]T {
	if from == nil {
		return nil
	}
	tList := make([]T, 0, len(*from))
	for _, elem := range *from {
		var t T
		err := copier.Copy(&t, elem)
		if err != nil {
			log.Fatal(err)
		}
		tList = append(tList, t)
	}
	return &tList
}

func MustConvertToListPointer[F any, T any](from *[]F) *[]*T {
	if from == nil {
		return nil
	}
	tList := make([]*T, 0, len(*from))
	for _, elem := range *from {
		var t T
		err := copier.Copy(&t, elem)
		if err != nil {
			log.Fatal(err.Error())
		}
		tList = append(tList, &t)
	}
	return &tList
}
