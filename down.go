package main

import (
	"os"
	"io"

	"net/http"
	"fmt"
	"path"
	"sync"
)

var (
	pathPrefix = "https://img.infinitynewtab.com/wallpaper/"
	pathSuffix = ".jpg"
)
var wg sync.WaitGroup

func main() {

	num := 4050

	photoNumCh := make(chan int, 4)

	go produce(num, photoNumCh)

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go downloadFile(photoNumCh)
	}

	wg.Wait()
}
func produce(num int, photoNumCh chan int) {
	for i := 1; i <= num; i++ {
		photoNumCh <- i
	}
	close(photoNumCh)
}

func downloadFile(photoNumCh chan int) {
	defer wg.Done()

	for photoNum, ok := <-photoNumCh; ok; photoNum, ok = <-photoNumCh {

		strFinal := fmt.Sprintf(pathPrefix+"%d"+pathSuffix, photoNum)
		res, err := http.Get(strFinal)
		if err != nil {
			continue
		}
		fileName := path.Base(strFinal)
		f, err := os.Create("E:/photo/" + fileName)
		if err != nil {
			continue
		}
		fileSize, writeErr := io.Copy(f, res.Body)

		fmt.Println(fileName+" download done,", "file size(byte)=", fileSize)
		if err != nil {
			fmt.Println(fileName+" download failed ", "errorInfo=", writeErr.Error())
			continue
		}

		continue
	}
}
