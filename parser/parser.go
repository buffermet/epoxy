package parser

/*
*	
*	File parser
*	
*	Embeds resources in HTML, CSS and SVG files using data URLs.
*	
*/

import(
	"sync"
	"mime"
	"regexp"
	"strings"
	"strconv"
	"runtime"
	"encoding/base64"

	"github.com/h2non/filetype"
	"github.com/yungtravla/epoxy/log"
	"github.com/yungtravla/epoxy/net"
	"github.com/yungtravla/epoxy/session"
)

func containsString(slice *[]string, str string) bool {
	for i := 0; i < len(*slice); i++ {
		if strings.EqualFold( (*slice)[i], str ) {
			return true
		}
	}
	return false
}

func pathToURL(path, origin string) string {
	r := regexp.MustCompile(`(?i)^http[s]?://[^/]+`)
	origin_host := r.FindString(origin)
	origin_path := r.ReplaceAllString(origin, "")
	origin_path = regexp.MustCompile(`[^/]*$`).ReplaceAllString(origin_path, "")

	url := ""

	if regexp.MustCompile(`^/`).FindString(path) == "" {
		if regexp.MustCompile(`(?i)^http[s]?://[^/]+`).FindString(path) == "" {
			if regexp.MustCompile(`^[.]`).FindString(path) == "" {
				url = origin_host + origin_path + path
			} else {
				count := len( regexp.MustCompile(`[.][.]/`).FindAllString(path, -1) )

				if count == 0 {
					log.Error("invalid path detected: " + path)
				} else if len( regexp.MustCompile(`[^/]+/`).FindAllString(origin_path, -1) ) < count {
					log.Error("unable to move out of directory: " + path + " (origin path is not long enough)")
				} else {
					stripped_path := origin_path
					for i := 0; i < count; i++ {
						stripped_path = regexp.MustCompile(`[^/]+/$`).ReplaceAllString(stripped_path, "")
					}

					url = origin_host + stripped_path + regexp.MustCompile(`^(?:[.][.]/)*`).ReplaceAllString(path, "")
				}
			}
		} else {
			url = path
		}
	} else {
		url = origin_host + path
	}

	return url
}

func findResources(s *session.SessionConfig) []string {
	var resources []string

	matches_src := regexp.MustCompile(`(?i)src=["'](.*?)["']`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_src); i++ {
		matches_src[i] = regexp.MustCompile(`(?i)src=["'](.*?)["']`).ReplaceAllString(matches_src[i], "${1}")
		if !strings.HasPrefix(matches_src[i], "data:") {
			resources = append(resources, matches_src[i])
		}
	}

	matches_href := regexp.MustCompile(`(?i)href=["'](.*?)["']`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_href); i++ {
		matches_href[i] = regexp.MustCompile(`(?i)href=["'](.*?)["']`).ReplaceAllString(matches_href[i], "${1}")
		if !strings.HasPrefix(matches_href[i], "data:") {
			resources = append(resources, matches_href[i])
		}
	}

	matches_url := regexp.MustCompile(`(?i)url[(]["']?(.*?)["']?[)]`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_url); i++ {
		matches_url[i] = regexp.MustCompile(`(?i)url[(]["']?(.*?)["']?[)]`).ReplaceAllString(matches_url[i], "${1}")
		if !strings.HasPrefix(matches_url[i], "data:") {
			resources = append(resources, matches_url[i])
		}
	}

	if len(resources) > 1 {
		log.Success("found " + strconv.Itoa( len(resources) ) + " embeddable resources in " + log.BOLD + s.Source + log.RESET + ".")
	} else if len(resources) == 1 {
		log.Success("found 1 link to a resource in " + log.BOLD + s.Source + log.RESET + ".")
	} else {
		log.Info("no resources found in " + log.BOLD + s.Source + log.RESET + ".")
	}

	return resources
}

func createDataURL(mimetype string, payload *[]byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(*payload)

	return []byte("data:" + mimetype + ";base64," + encoded)
}

func embedResources(s *session.SessionConfig) session.SessionConfig {
	matches_src := regexp.MustCompile(`(?i)src=["'][^"']+["']`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_src); i++ {
		path := regexp.MustCompile(`(?i)src=["']([^"']+)["']`).ReplaceAllString( string(matches_src[i]), "${1}" )

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			path = strings.Replace(path, "?", "\\?", -1)
			path = strings.Replace(path, "-", "\\-", -1)
			path = strings.Replace(path, ".", "\\.", -1)
			path = strings.Replace(path, "+", "\\+", -1)

			new_source := regexp.MustCompile(`(?i)src=("|')` + path + `("|')`).ReplaceAllString( string(s.Body), "src=${1}" + string(body) + "${2}" )
			s.Body = []byte(new_source)
		}
	}

	matches_href := regexp.MustCompile(`(?i)href=["'][^"']+["']?`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_href); i++ {
		path := regexp.MustCompile(`(?i)href=["']([^"']+)["']?`).ReplaceAllString( string(matches_href[i]), "$1" )

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			path = strings.Replace(path, "?", "\\?", -1)
			path = strings.Replace(path, "-", "\\-", -1)
			path = strings.Replace(path, ".", "\\.", -1)
			path = strings.Replace(path, "+", "\\+", -1)

			new_source := regexp.MustCompile(`(?i)href=("|')` + path + `("|')`).ReplaceAllString( string(s.Body), "href=${1}" + string(body) + "${2}" )
			s.Body = []byte(new_source)			
		}
	}

	matches_url := regexp.MustCompile(`(?i)url[(]["']?[^"')]+["']?[)]`).FindAllString( string(s.Body), -1 )

	for i := 0; i < len(matches_url); i++ {
		path := regexp.MustCompile(`(?i)url[(]["']?([^"')]+)["']?[)]`).ReplaceAllString(matches_url[i], "$1")

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			path = strings.Replace(path, "?", "\\?", -1)
			path = strings.Replace(path, "-", "\\-", -1)
			path = strings.Replace(path, ".", "\\.", -1)
			path = strings.Replace(path, "+", "\\+", -1)

			new_source := regexp.MustCompile(`(?i)url[(]("|'|)` + path + `("|'|)[)]`).ReplaceAllString( string(s.Body), "url(${1}" + string(body) + "${2})" )
			s.Body = []byte(new_source)
		}
	}

	return *s
}

func Parse(s *session.SessionConfig) session.SessionConfig {
	origin_path := regexp.MustCompile(`(?i)http[s]?://[^/]+`).ReplaceAllString(s.Origin, "")

	if origin_path == "" {
		origin_path = "/"
	} else if !strings.HasSuffix(origin_path, "/") {
		origin_path = regexp.MustCompile(`[^/]*$`).ReplaceAllString(origin_path, "")
	}

	resources := findResources(s)

	answer := ""
	if !s.Recursive {
		answer = log.Prompt( "fetch " + strconv.Itoa( len(resources) ) + " resource(s)? Y/n" )
	}

	if !s.Recursive || answer != "n" && answer != "N" {
		for i := 0; i < len(resources); i++ {
			if resources[i] != "" && !strings.HasPrefix(resources[i], "#") {
				var resource session.Resource

				address := pathToURL(resources[i], s.Origin)
				resource.Address = address

				s.RequestQueue.Add(1)

				runtime.GOMAXPROCS(8)
				go func(path string) {
					stripped_address := regexp.MustCompile(`(?:\?|#).*$`).ReplaceAllString(address, "")

					extension := regexp.MustCompile(`\.([a-zA-Z0-9)]+)$`).FindString(stripped_address)

					extension_mimetype := strings.Replace( mime.TypeByExtension(extension) , " ", "", -1 )
					if extension_mimetype == "" { extension_mimetype = "unknown" }

					if containsString(&s.Accept, extension_mimetype) {
						body, content_type := net.SendRequest(address, s)

						content_type = strings.Replace(content_type, " ", "", -1)

						sniffed_mimetype, err := filetype.Match(body)
						if err != nil {
							log.Info( "could not determine filetype, using Content-Type header value: " + content_type + "(" + err.Error() + ")" )
						}

						sniffed_mimetype.MIME.Value = strings.Replace(sniffed_mimetype.MIME.Value, " ", "", -1)

						if sniffed_mimetype.MIME.Value != "" {
							content_type = sniffed_mimetype.MIME.Value
						}

						resource.Type = content_type

						if containsString(&s.Accept, content_type) {
							log.Success( strconv.Itoa( len(body) ) + " B " + log.BOLD + "[" + content_type + "]" + log.RESET + " " + address )

							if regexp.MustCompile(`(?:text/(?:css|html)|image/svg\+xml)`).FindString(content_type) != "" {
								_s := session.SessionConfig{ resource.Address, resource.Address, body, s.Accept, []session.Resource{}, sync.WaitGroup{}, true }

								_s = Parse(&_s)

								resource.Body = _s.Body
							} else {
								resource.Body = body
							}

							s.Resources = append(s.Resources, resource)
						} else {
							log.Info( "skipping response: " + strconv.Itoa( len(body) ) + " B " + log.BOLD + "[" + content_type + "]" + log.RESET + " " + address )
						}
					} else {
						log.Info( "skipping request: " + log.BOLD + "[" + extension_mimetype + "]" + log.RESET + " " + address )
					}

					s.RequestQueue.Done()
				}(resources[i])
			}
		}

		s.RequestQueue.Wait()

		if len(s.Resources) > 0 {
			log.Info("generating base64 encoded data URLs ...")

			for i := 0; i < len(s.Resources); i++ {
				data_url := createDataURL(s.Resources[i].Type, &s.Resources[i].Body)

				s.Resources[i].Body = data_url
			}

			if s.Recursive {
				log.Info("embedding resources in " + log.BOLD + s.Source + log.RESET + " ...")
			} else {
				log.Info("embedding resources in source file ...")
			}

			*s = embedResources(s)
		}

		return *s
	} else {
		log.Raw("")
	}

	return *s
}
