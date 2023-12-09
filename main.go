package main

/*
*	
*	Launcher
*	
*/

import(
	"runtime"
	"io/ioutil"

	"github.com/buffermet/epoxy/log"
	"github.com/buffermet/epoxy/parser"
	"github.com/buffermet/epoxy/session"
)

func initiatePrint(s *session.SessionConfig) {
	log.Raw( string( parser.Parse(s).Body ) )
}

func initiateWrite(s *session.SessionConfig) {
	if s.Recurse > 0 {
		log.Info("parsing " + s.Source + " ...")

		*s = parser.Parse(s)

		log.Info("saving payload as " + log.BOLD + "epoxy-" + s.Source + log.RESET + " ...")

		ioutil.WriteFile("epoxy-" + s.Source, s.Body, 0600)
	} else {
		log.Info("encoding " + s.Source + " ...")

		*s = parser.Parse(s)

		log.Info("saving payload as " + log.BOLD + s.Source + ".url" + log.RESET + " ...")

		ioutil.WriteFile(s.Source + ".url", s.Body, 0600)
	}
}

func main() {
	log.Raw("")

	s := session.NewSession()

	runtime.GOMAXPROCS(session.Cores)

	if session.Print {
		initiatePrint(&s)
	} else {
		initiateWrite(&s)
	}

	log.Raw("")
}
