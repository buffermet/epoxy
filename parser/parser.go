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
	"encoding/base64"

	"github.com/h2non/filetype"
	"github.com/buffermet/epoxy/log"
	"github.com/buffermet/epoxy/net"
	"github.com/buffermet/epoxy/session"
)

func containsString(slice *[]string, str string) bool {
	for i := 0; i < len(*slice); i++ {
		if strings.EqualFold((*slice)[i], str) {
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
	origin_scheme := regexp.MustCompile(`(?i)^[a-z]+:`).FindString(origin)

	url := ""

	if regexp.MustCompile(`^//`).FindString(path) == "" {
		if regexp.MustCompile(`^/`).FindString(path) == "" {
			if regexp.MustCompile(`(?i)^http[s]?://[^/]+`).FindString(path) == "" {
				if regexp.MustCompile(`^[.]`).FindString(path) == "" {
					url = origin_host + origin_path + path
				} else {
					count := len(regexp.MustCompile(`[.][.]/`).FindAllString(path, -1))

					if count == 0 {
						log.Error("invalid path detected: " + path)
					} else if len(regexp.MustCompile(`[^/]+/`).FindAllString(origin_path, -1)) < count {
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
	} else {
		url = origin_scheme + path
	}

	return url
}

func findResources(s *session.SessionConfig) []string {
	var resources []string

	matches_src := regexp.MustCompile(`(?i)src=["'](.*?)["']`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_src); i++ {
		matches_src[i] = regexp.MustCompile(`(?i)src=["'](.*?)["']`).ReplaceAllString(matches_src[i], "${1}")
		if !strings.HasPrefix(matches_src[i], "data:") {
			resources = append(resources, matches_src[i])
		}
	}

	matches_content := regexp.MustCompile(`(?i)content=["'](.*?)["']`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_content); i++ {
		matches_content[i] = regexp.MustCompile(`(?i)content=["'](.*?)["']`).ReplaceAllString(matches_content[i], "${1}")
		if !strings.HasPrefix(matches_content[i], "data:") {
			resources = append(resources, matches_content[i])
		}
	}

	matches_href := regexp.MustCompile(`(?i)href=["'](.*?)["']`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_href); i++ {
		matches_href[i] = regexp.MustCompile(`(?i)href=["'](.*?)["']`).ReplaceAllString(matches_href[i], "${1}")
		if !strings.HasPrefix(matches_href[i], "data:") {
			resources = append(resources, matches_href[i])
		}
	}

	matches_url := regexp.MustCompile(`(?i)url[(]["']?(.*?)["']?[)]`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_url); i++ {
		matches_url[i] = regexp.MustCompile(`(?i)url[(]["']?(.*?)["']?[)]`).ReplaceAllString(matches_url[i], "${1}")
		if !strings.HasPrefix(matches_url[i], "data:") {
			resources = append(resources, matches_url[i])
		}
	}

	unique_resources := []string{}
	m := make(map[string]bool)
	for _, entry := range resources {
		if _, value := m[entry]; !value {
			m[entry] = true
			unique_resources = append(unique_resources, entry)
		}
	}

	resources = unique_resources

	if len(resources) > 1 {
		log.Success("found " + strconv.Itoa(len(resources)) + " embeddable resources in " + log.BOLD + s.Source + log.RESET + ".")
	} else if len(resources) == 1 {
		log.Success("found 1 link to a resource in " + log.BOLD + s.Source + log.RESET + ".")
	} else {
		log.Info("no resources found in " + log.BOLD + s.Source + log.RESET + ".")
	}

	return resources
}

func createDataURL(mimetype string, payload *[]byte) []byte {
	encoded_body := base64.StdEncoding.EncodeToString(*payload)

	return []byte("data:" + mimetype + ";base64," + encoded_body)
}

func embedResources(s *session.SessionConfig) session.SessionConfig {
	matches_src := regexp.MustCompile(`(?i)src=["'][^"']+["']`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_src); i++ {
		path := regexp.MustCompile(`(?i)src=["']([^"']+)["']`).ReplaceAllString(string(matches_src[i]), "${1}")

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			if len(body) > 0 {
				path = strings.Replace(path, "?", "\\?", -1)
				path = strings.Replace(path, "-", "\\-", -1)
				path = strings.Replace(path, ".", "\\.", -1)
				path = strings.Replace(path, "+", "\\+", -1)

				new_source := regexp.MustCompile(`(?i)src=("|')` + path + `("|')`).ReplaceAllString(string(s.Body), "src=${1}" + string(body) + "${2}")
				s.Body = []byte(new_source)
			}
		}
	}

	matches_content := regexp.MustCompile(`(?i)content=["'][^"']+["']`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_content); i++ {
		path := regexp.MustCompile(`(?i)content=["']([^"']+)["']`).ReplaceAllString(string(matches_content[i]), "${1}")

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			if len(body) > 0 {
				path = strings.Replace(path, "?", "\\?", -1)
				path = strings.Replace(path, "-", "\\-", -1)
				path = strings.Replace(path, ".", "\\.", -1)
				path = strings.Replace(path, "+", "\\+", -1)

				new_source := regexp.MustCompile(`(?i)content=("|')` + path + `("|')`).ReplaceAllString(string(s.Body), "content=${1}" + string(body) + "${2}")
				s.Body = []byte(new_source)
			}
		}
	}

	matches_href := regexp.MustCompile(`(?i)href=["'][^"']+["']?`).FindAllString(string(s.Body), -1)

	for i := 0; i < len(matches_href); i++ {
		path := regexp.MustCompile(`(?i)href=["']([^"']+)["']?`).ReplaceAllString(string(matches_href[i]), "$1")

		if regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(path) == "" {
			address := pathToURL(path, s.Origin)

			var body []byte

			for a := 0; a < len(s.Resources); a++ {
				if address == s.Resources[a].Address {
					body = s.Resources[a].Body
				}
			}

			if len(body) > 0 {
				path = strings.Replace(path, "?", "\\?", -1)
				path = strings.Replace(path, "-", "\\-", -1)
				path = strings.Replace(path, ".", "\\.", -1)
				path = strings.Replace(path, "+", "\\+", -1)

				new_source := regexp.MustCompile(`(?i)href=("|')` + path + `("|')`).ReplaceAllString(string(s.Body), "href=${1}" + string(body) + "${2}")
				s.Body = []byte(new_source)			
			}
		}
	}

	matches_url := regexp.MustCompile(`(?i)url[(]["']?[^"')]+["']?[)]`).FindAllString(string(s.Body), -1)

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

			if len(body) > 0 {
				path = strings.Replace(path, "?", "\\?", -1)
				path = strings.Replace(path, "-", "\\-", -1)
				path = strings.Replace(path, ".", "\\.", -1)
				path = strings.Replace(path, "+", "\\+", -1)

				new_source := regexp.MustCompile(`(?i)url[(]("|'|)` + path + `("|'|)[)]`).ReplaceAllString(string(s.Body), "url(${1}" + string(body) + "${2})")
				s.Body = []byte(new_source)
			}
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

	if s.Recurse != 0 {
		resources := findResources(s)

		session.Depth++

		answer := ""

		if session.Depth == 1 {
			answer = log.Prompt("fetch at least " + strconv.Itoa(len(resources)) + " resource(s)? Y/n")
		}

		if session.Depth != 1 || answer != "n" && answer != "N" {
			for i := 0; i < len(resources); i++ {
				if resources[i] != "" && regexp.MustCompile(`(?i)^(?:data:|javascript:|#)`).FindString(resources[i]) == "" {
					var resource session.Resource

					address := pathToURL(resources[i], s.Origin)
					resource.Address = address

					if session.Depth <= s.Recurse {
						s.RequestQueue.Add(1)

						///// ASYNC /////
						go func(path string) {
							stripped_address := regexp.MustCompile(`(?:\?|#).*$`).ReplaceAllString(address, "")

							extension := regexp.MustCompile(`\.([a-zA-Z0-9)]+)$`).FindString(stripped_address)

							extension_mimetype := strings.Replace(mime.TypeByExtension(extension) , " ", "", -1)
							if extension_mimetype == "" { extension_mimetype = "unknown" }

							if containsString(&s.Accept, extension_mimetype) {
								body, content_type := net.SendRequest(address, s)

								content_type = strings.Replace(content_type, " ", "", -1)

								parsed_mimetype, err := filetype.Match(body)
								if err != nil { log.Info("could not determine filetype, using Content-Type header value: " + content_type + "(" + err.Error() + ")") }

								parsed_mimetype.MIME.Value = strings.Replace(parsed_mimetype.MIME.Value, " ", "", -1)

								if parsed_mimetype.MIME.Value != "" {
									content_type = parsed_mimetype.MIME.Value
								}

								content_type = regexp.MustCompile(`;.*`).ReplaceAllString(content_type, "")

								resource.Type = content_type

								if containsString(&s.Accept, content_type) {
									log.Success(strconv.Itoa(len(body)) + " B " + log.BOLD + "[" + content_type + "]" + log.RESET + " " + address)

									if regexp.MustCompile(`(?:text/(?:css|html)|image/svg\+xml)`).FindString(content_type) != "" {
										_s := session.SessionConfig { 
											resource.Address,             // Source string
											resource.Address,             // Origin string
											body,                         // Body []byte
											s.Accept,                     // Accept []string
											(s.Recurse - session.Depth),  // Recurse int
											[]session.Resource{},         // Resources []Resource
											sync.WaitGroup{},             // RequestQueue sync.WaitGroup
										}

										_s = Parse(&_s)

										resource.Body = _s.Body
									} else {
										resource.Body = body
									}

									s.Resources = append(s.Resources, resource)
								} else {
									log.Info("skipping response: " + strconv.Itoa(len(body)) + " B " + log.BOLD + "[" + content_type + "]" + log.RESET + " " + address)
								}
							} else {
								log.Info("skipping request: " + log.BOLD + "[" + extension_mimetype + "]" + log.RESET + " " + address)
							}

							s.RequestQueue.Done()
						}(resources[i])
						///// SYNC /////
					}
				}
			}

			s.RequestQueue.Wait()

			if len(s.Resources) > 0 {
				log.Info("generating base64 encoded data URLs ...")

				for i := 0; i < len(s.Resources); i++ {
					data_url := createDataURL(s.Resources[i].Type, &s.Resources[i].Body)

					s.Resources[i].Body = data_url
				}

				log.Info("embedding resources in " + log.BOLD + s.Source + log.RESET + " ...")

				*s = embedResources(s)
			}

			return *s
		} else {
			log.Raw("")
		}
	} else { // if s.Recurse is 0
		content_type := ""

		extension_mimetype := ""
		if s.Origin != "" {
			stripped_address := regexp.MustCompile(`(?:\?|#).*$`).ReplaceAllString(s.Origin, "")

			extension := regexp.MustCompile(`\.([a-zA-Z0-9)]+)$`).FindString(stripped_address)
			extension_mimetype = strings.Replace(mime.TypeByExtension(extension) , " ", "", -1)
		} else {
			extension := regexp.MustCompile(`\.([a-zA-Z0-9)]+)$`).FindString(s.Source)
			extension_mimetype = strings.Replace(mime.TypeByExtension(extension) , " ", "", -1)
		}

		parsed_mimetype, err := filetype.Match(s.Body)
		if err != nil { log.Info("could not determine filetype, using Content-Type header value: " + content_type + "(" + err.Error() + ")") }
		parsed_mimetype.MIME.Value = strings.Replace(parsed_mimetype.MIME.Value, " ", "", -1)

		if parsed_mimetype.MIME.Value != "" {
			content_type = parsed_mimetype.MIME.Value
		} else if extension_mimetype != "" {
			content_type = extension_mimetype
		}

		content_type = regexp.MustCompile(`;.*`).ReplaceAllString(content_type, "")

		if content_type != "" {
			data_url := "data:" + content_type + ";base64,"

			encoded_body := base64.StdEncoding.EncodeToString(s.Body)

			s.Body = []byte(data_url + encoded_body)

			return *s
		} else {
			log.Fatal("could not determine filetype, please specify it manually using -mimetype STRING")
		}
	}

	return *s
}
