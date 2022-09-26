package initapp

//初始化的源头
import (
	_ "blog/resource/system"
	"log"
)

func init() {
	log.Println("initapp server mysql init")
	// do nothing,only import source package so that inits can be registered
}
