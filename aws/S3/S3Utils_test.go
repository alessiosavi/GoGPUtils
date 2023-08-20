package S3utils

import (
	csvutils "github.com/alessiosavi/GoGPUtils/csv"
	"github.com/alessiosavi/GoGPUtils/helper"
	"log"
	"os"
	"testing"
)

func Test1(t *testing.T) {
	objects, err := ListBucketObjects("thom-browne-images", "farfetch_remapped")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("farfetch_remapped.json", []byte(helper.Marshal(objects)), 0755)
	if err != nil {
		panic(err)
	}

	objects, err = ListBucketObjects("thom-browne-images", "farfetch_raw")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("farfetch_raw.json", []byte(helper.Marshal(objects)), 0755)
	if err != nil {
		panic(err)
	}
	object, err := GetObject("thom-browne-images", "farfetch_raw/valid_images.csv")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("valid_images.csv", object, 0755)
	if err != nil {
		panic(err)
	}

}

func Test2(t *testing.T) {
	f := "/mnt/c/Users/alessio.savi/AppData/Local/Temp/image/image_not_processed.csv"
	objects, err := ListBucketObjects("thom-browne-images", "farfetch_raw")
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}
	_, csvData, err := csvutils.ReadCSV(file, ',')
	var errs []string
	if err != nil {
		panic(err)
	}
	for i := range objects {
		for k := range csvData {
			if objects[i] == csvData[k][0] {
				errs = append(errs, objects[i])
				continue
			}
		}
	}
	log.Println(errs)
}
