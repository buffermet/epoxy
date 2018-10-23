package main

/*
*	
*	Launcher
*	
*/

import(
	"io/ioutil"

	"github.com/yungtravla/epoxy/log"
	"github.com/yungtravla/epoxy/parser"
	"github.com/yungtravla/epoxy/session"
)

func main() {
	log.Raw("")

	s := session.NewSession()

	log.Info("parsing " + s.Source + " ...")

	s = parser.Parse(&s)

	log.Info("saving payload as epoxy-" + s.Source + " ...")

	ioutil.WriteFile("epoxy-" + s.Source, s.Body, 0600)

	log.Raw("")
}
