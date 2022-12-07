package main

import gt "github.com/mangenotwork/gathertool"

func main() {
	caseUrl := "http://127.0.0.1:14444/search/case1?word=%E7%A6%8F%E5%BD%A9&sort=o&out=full&pg=255&count=999"

	test := gt.NewTestUrl(caseUrl, "Get", 10, 1)
	test.Run()

}
