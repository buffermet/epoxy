package session

/*
*	
*	Parser configuration
*	
*/

import(
	"os"
	"sync"
	"io/ioutil"
	"regexp"

	"github.com/yungtravla/epoxy/log"
)

type Resource struct {
	Type string
	Address string
	Body []byte
}

type SessionConfig struct {
	Source string
	Origin string
	Body []byte
	Accept []string
	Resources []Resource

	RequestQueue sync.WaitGroup
	SessionQueue sync.WaitGroup

	Recursive bool
}

func showOptions() {
	str := "usage: epoxy <options> -source <path> -origin <url>\n" + 
	       "\n" + 
	       "Options:\n" + 
	       "\n" + 
	       "  -source PATH    path to source file.\n" + 
	       "  -origin URL     full URL to source file.\n" + 
	       "  -recurse INT    recursion depth of resource fetching.\n" + 
	       "  -cores INT      limit of cores to use for async parsing.\n" + 
	       "\n" + 
	       "  -no-unknown     don't embed unknown filetypes.\n" + 
	       "  -no-svg         don't embed svg files.\n" + 
	       "  -no-jpg         don't embed jpg files.\n" + 
	       "  -no-png         don't embed png files.\n" + 
	       "  -no-gif         don't embed gif files.\n" + 
	       "  -no-webp        don't embed webp files.\n" + 
	       "  -no-cr2         don't embed cr2 files.\n" + 
	       "  -no-tif         don't embed tif files.\n" + 
	       "  -no-bmp         don't embed bmp files.\n" + 
	       "  -no-jxr         don't embed jxr files.\n" + 
	       "  -no-psd         don't embed psd files.\n" + 
	       "  -no-ico         don't embed ico files.\n" + 
	       "  -no-mp4         don't embed mp4 files.\n" + 
	       "  -no-m4v         don't embed m4v files.\n" + 
	       "  -no-mkv         don't embed mkv files.\n" + 
	       "  -no-webm        don't embed webm files.\n" + 
	       "  -no-mov         don't embed mov files.\n" + 
	       "  -no-avi         don't embed avi files.\n" + 
	       "  -no-wmv         don't embed wmv files.\n" + 
	       "  -no-mpg         don't embed mpg files.\n" + 
	       "  -no-flv         don't embed flv files.\n" + 
	       "  -no-mid         don't embed mid files.\n" + 
	       "  -no-mp3         don't embed mp3 files.\n" + 
	       "  -no-m4a         don't embed m4a files.\n" + 
	       "  -no-ogg         don't embed ogg files.\n" + 
	       "  -no-flac        don't embed flac files.\n" + 
	       "  -no-wav         don't embed wav files.\n" + 
	       "  -no-amr         don't embed amr files.\n" + 
	       "  -no-epub        don't embed epub files.\n" + 
	       "  -no-zip         don't embed zip files.\n" + 
	       "  -no-tar         don't embed tar files.\n" + 
	       "  -no-rar         don't embed rar files.\n" + 
	       "  -no-gz          don't embed gz files.\n" + 
	       "  -no-bz2         don't embed bz2 files.\n" + 
	       "  -no-7z          don't embed 7z files.\n" + 
	       "  -no-xz          don't embed xz files.\n" + 
	       "  -no-pdf         don't embed pdf files.\n" + 
	       "  -no-exe         don't embed exe files.\n" + 
	       "  -no-swf         don't embed swf files.\n" + 
	       "  -no-rtf         don't embed rtf files.\n" + 
	       "  -no-eot         don't embed eot files.\n" + 
	       "  -no-ps          don't embed ps files.\n" + 
	       "  -no-sqlite      don't embed sqlite files.\n" + 
	       "  -no-nes         don't embed nes files.\n" + 
	       "  -no-crx         don't embed crx files.\n" + 
	       "  -no-cab         don't embed cab files.\n" + 
	       "  -no-deb         don't embed deb files.\n" + 
	       "  -no-ar          don't embed ar files.\n" + 
	       "  -no-z           don't embed z files.\n" + 
	       "  -no-lz          don't embed lz files.\n" + 
	       "  -no-rpm         don't embed rpm files.\n" + 
	       "  -no-elf         don't embed elf files.\n" + 
	       "  -no-doc         don't embed doc files.\n" + 
	       "  -no-docx        don't embed docx files.\n" + 
	       "  -no-xls         don't embed xls files.\n" + 
	       "  -no-xlsx        don't embed xlsx files.\n" + 
	       "  -no-ppt         don't embed ppt files.\n" + 
	       "  -no-pptx        don't embed pptx files.\n" + 
	       "  -no-woff        don't embed woff files.\n" + 
	       "  -no-woff2       don't embed woff2 files.\n" + 
	       "  -no-ttf         don't embed ttf files.\n" + 
	       "  -no-otf         don't embed otf files.\n" + 
	       "  -no-css         don't embed css files.\n" + 
	       "  -no-html        don't embed html files.\n" + 
	       "  -no-js          don't embed js files.\n" + 
	       "  -no-json        don't embed json files.\n"

	log.Raw(str)

	os.Exit(0)
}

func skipMimetype(mimetype string, session *SessionConfig) {
	for i := 0; i < len(session.Accept); i++ {
		if mimetype == session.Accept[i] {
			session.Accept = append(session.Accept[:i], session.Accept[i+1:]...)
			log.Info("skipping file type: " + mimetype)
		}
	}
}

func NewSession() SessionConfig {
	accept := []string{"unknown", "application/octet-stream", "image/svg", "image/svg+xml", "image/jpeg", "image/png", "image/gif", "image/webp", "image/x-canon-cr2", "image/tiff", "image/bmp", "image/vnd.ms-photo", "image/vnd.adobe.photoshop", "image/vnd.microsoft.icon", "image/x-icon", "video/mp4", "video/x-m4v", "video/x-matroska", "video/webm", "video/quicktime", "video/x-msvideo", "video/x-ms-wmv", "video/mpeg", "video/x-flv", "audio/midi", "audio/mpeg", "audio/m4a", "audio/ogg", "audio/x-flac", "audio/x-wav", "audio/amr", "application/epub+zip", "application/zip", "application/x-tar", "application/x-rar-compressed", "application/gzip", "application/x-bzip2", "application/x-7z-compressed", "application/x-xz", "application/pdf", "application/x-msdownload", "application/x-shockwave-flash", "application/rtf", "application/vnd.ms-fontobject", "font/eot", "application/postscript", "application/x-sqlite3", "application/x-nintendo-nes-rom", "application/x-google-chrome-extension", "application/vnd.ms-cab-compressed", "application/x-deb", "application/x-unix-archive", "application/x-compress", "application/x-lzip", "application/x-rpm", "application/x-executable", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation", "application/font-woff", "font/woff", "application/font-woff", "font/woff2", "application/font-sfnt", "font/ttf", "application/font-sfnt", "font/otf", "text/css", "text/css;charset=UTF-8", "text/html", "text/html;charset=UTF-8", "text/javascript", "application/javascript", "application/x-javascript", "text/json", "application/json"}

	s := SessionConfig{ "", "", []byte(""), accept, []Resource{}, sync.WaitGroup{}, sync.WaitGroup{}, false }

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		if args[i] == "--help" || args[i] == "-help" {
			showOptions()
		} else if args[i] == "--no-unknown" || args[i] == "-no-unknown" {
			skipMimetype("unknown", &s)
			skipMimetype("application/octet-stream", &s)
		} else if args[i] == "--no-svg" || args[i] == "-no-svg" {
			skipMimetype("image/svg+xml", &s)
			skipMimetype("image/svg", &s)
		} else if args[i] == "--no-jpg" || args[i] == "-no-jpg" {
			skipMimetype("image/jpeg", &s)
		} else if args[i] == "--no-png" || args[i] == "-no-png" {
			skipMimetype("image/png", &s)
		} else if args[i] == "--no-gif" || args[i] == "-no-gif" {
			skipMimetype("image/gif", &s)
		} else if args[i] == "--no-webp" || args[i] == "-no-webp" {
			skipMimetype("image/webp", &s)
		} else if args[i] == "--no-cr2" || args[i] == "-no-cr2" {
			skipMimetype("image/x-canon-cr2", &s)
		} else if args[i] == "--no-tif" || args[i] == "-no-tif" {
			skipMimetype("image/tiff", &s)
		} else if args[i] == "--no-bmp" || args[i] == "-no-bmp" {
			skipMimetype("image/bmp", &s)
		} else if args[i] == "--no-jxr" || args[i] == "-no-jxr" {
			skipMimetype("image/vnd.ms-photo", &s)
		} else if args[i] == "--no-psd" || args[i] == "-no-psd" {
			skipMimetype("image/vnd.adobe.photoshop", &s)
		} else if args[i] == "--no-ico" || args[i] == "-no-ico" {
			skipMimetype("image/vnd.microsoft.icon", &s)
			skipMimetype("image/x-icon", &s)
		} else if args[i] == "--no-mp4" || args[i] == "-no-mp4" {
			skipMimetype("video/mp4", &s)
		} else if args[i] == "--no-m4v" || args[i] == "-no-m4v" {
			skipMimetype("video/x-m4v", &s)
		} else if args[i] == "--no-mkv" || args[i] == "-no-mkv" {
			skipMimetype("video/x-matroska", &s)
		} else if args[i] == "--no-webm" || args[i] == "-no-webm" {
			skipMimetype("video/webm", &s)
		} else if args[i] == "--no-mov" || args[i] == "-no-mov" {
			skipMimetype("video/quicktime", &s)
		} else if args[i] == "--no-avi" || args[i] == "-no-avi" {
			skipMimetype("video/x-msvideo", &s)
		} else if args[i] == "--no-wmv" || args[i] == "-no-wmv" {
			skipMimetype("video/x-ms-wmv", &s)
		} else if args[i] == "--no-mpg" || args[i] == "-no-mpg" {
			skipMimetype("video/mpeg", &s)
		} else if args[i] == "--no-flv" || args[i] == "-no-flv" {
			skipMimetype("video/x-flv", &s)
		} else if args[i] == "--no-mid" || args[i] == "-no-mid" {
			skipMimetype("audio/midi", &s)
		} else if args[i] == "--no-mp3" || args[i] == "-no-mp3" {
			skipMimetype("audio/mpeg", &s)
		} else if args[i] == "--no-m4a" || args[i] == "-no-m4a" {
			skipMimetype("audio/m4a", &s)
		} else if args[i] == "--no-ogg" || args[i] == "-no-ogg" {
			skipMimetype("audio/ogg", &s)
		} else if args[i] == "--no-flac" || args[i] == "-no-flac" {
			skipMimetype("audio/x-flac", &s)
		} else if args[i] == "--no-wav" || args[i] == "-no-wav" {
			skipMimetype("audio/x-wav", &s)
		} else if args[i] == "--no-amr" || args[i] == "-no-amr" {
			skipMimetype("audio/amr", &s)
		} else if args[i] == "--no-epub" || args[i] == "-no-epub" {
			skipMimetype("application/epub+zip", &s)
		} else if args[i] == "--no-zip" || args[i] == "-no-zip" {
			skipMimetype("application/zip", &s)
		} else if args[i] == "--no-tar" || args[i] == "-no-tar" {
			skipMimetype("application/x-tar", &s)
		} else if args[i] == "--no-rar" || args[i] == "-no-rar" {
			skipMimetype("application/x-rar-compressed", &s)
		} else if args[i] == "--no-gz" || args[i] == "-no-gz" {
			skipMimetype("application/gzip", &s)
		} else if args[i] == "--no-bz2" || args[i] == "-no-bz2" {
			skipMimetype("application/x-bzip2", &s)
		} else if args[i] == "--no-7z" || args[i] == "-no-7z" {
			skipMimetype("application/x-7z-compressed", &s)
		} else if args[i] == "--no-xz" || args[i] == "-no-xz" {
			skipMimetype("application/x-xz", &s)
		} else if args[i] == "--no-pdf" || args[i] == "-no-pdf" {
			skipMimetype("application/pdf", &s)
		} else if args[i] == "--no-exe" || args[i] == "-no-exe" {
			skipMimetype("application/x-msdownload", &s)
		} else if args[i] == "--no-swf" || args[i] == "-no-swf" {
			skipMimetype("application/x-shockwave-flash", &s)
		} else if args[i] == "--no-rtf" || args[i] == "-no-rtf" {
			skipMimetype("application/rtf", &s)
		} else if args[i] == "--no-eot" || args[i] == "-no-eot" {
			skipMimetype("application/vnd.ms-fontobject", &s)
			skipMimetype("font/eot", &s)
		} else if args[i] == "--no-ps" || args[i] == "-no-ps" {
			skipMimetype("application/postscript", &s)
		} else if args[i] == "--no-sqlite" || args[i] == "-no-sqlite" {
			skipMimetype("application/x-sqlite3", &s)
		} else if args[i] == "--no-nes" || args[i] == "-no-nes" {
			skipMimetype("application/x-nintendo-nes-rom", &s)
		} else if args[i] == "--no-crx" || args[i] == "-no-crx" {
			skipMimetype("application/x-google-chrome-extension", &s)
		} else if args[i] == "--no-cab" || args[i] == "-no-cab" {
			skipMimetype("application/vnd.ms-cab-compressed", &s)
		} else if args[i] == "--no-deb" || args[i] == "-no-deb" {
			skipMimetype("application/x-deb", &s)
		} else if args[i] == "--no-ar" || args[i] == "-no-ar" {
			skipMimetype("application/x-unix-archive", &s)
		} else if args[i] == "--no-z" || args[i] == "-no-z" {
			skipMimetype("application/x-compress", &s)
		} else if args[i] == "--no-lz" || args[i] == "-no-lz" {
			skipMimetype("application/x-lzip", &s)
		} else if args[i] == "--no-rpm" || args[i] == "-no-rpm" {
			skipMimetype("application/x-rpm", &s)
		} else if args[i] == "--no-elf" || args[i] == "-no-elf" {
			skipMimetype("application/x-executable", &s)
		} else if args[i] == "--no-doc" || args[i] == "-no-doc" {
			skipMimetype("application/msword", &s)
		} else if args[i] == "--no-docx" || args[i] == "-no-docx" {
			skipMimetype("application/vnd.openxmlformats-officedocument.wordprocessingml.document", &s)
		} else if args[i] == "--no-xls" || args[i] == "-no-xls" {
			skipMimetype("application/vnd.ms-excel", &s)
		} else if args[i] == "--no-xlsx" || args[i] == "-no-xlsx" {
			skipMimetype("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", &s)
		} else if args[i] == "--no-ppt" || args[i] == "-no-ppt" {
			skipMimetype("application/vnd.ms-powerpoint", &s)
		} else if args[i] == "--no-pptx" || args[i] == "-no-pptx" {
			skipMimetype("application/vnd.openxmlformats-officedocument.presentationml.presentation", &s)
		} else if args[i] == "--no-woff" || args[i] == "-no-woff" {
			skipMimetype("application/font-woff", &s)
			skipMimetype("font/woff", &s)
		} else if args[i] == "--no-woff2" || args[i] == "-no-woff2" {
			skipMimetype("application/font-woff", &s)
			skipMimetype("font/woff2", &s)
		} else if args[i] == "--no-ttf" || args[i] == "-no-ttf" {
			skipMimetype("application/font-sfnt", &s)
			skipMimetype("font/ttf", &s)
		} else if args[i] == "--no-otf" || args[i] == "-no-otf" {
			skipMimetype("application/font-sfnt", &s)
			skipMimetype("font/otf", &s)
		} else if args[i] == "--no-css" || args[i] == "-no-css" {
			skipMimetype("text/css", &s)
			skipMimetype("text/css;charset=UTF-8", &s)
		} else if args[i] == "--no-html" || args[i] == "-no-html" {
			skipMimetype("text/html", &s)
			skipMimetype("text/html;charset=UTF-8", &s)
		} else if args[i] == "--no-js" || args[i] == "-no-js" {
			skipMimetype("text/javascript", &s)
			skipMimetype("application/javascript", &s)
			skipMimetype("application/x-javascript", &s)
		} else if args[i] == "--no-json" || args[i] == "-no-json" {
			skipMimetype("application/json", &s)
		} else if args[i] == "--source" || args[i] == "-source" {
			if i < ( len(args) - 1 ) {
				s.Source = args[i+1]
			} else {
				log.Error("missing value for: " + args[i] + "\n")
				showOptions()
			}
		} else if args[i] == "--origin" || args[i] == "-origin" {
			if i < ( len(args) - 1 ) {
				s.Origin = args[i+1]
			} else {
				log.Error("missing value for: " + args[i] + "\n")
				showOptions()
			}
		} else if args[i-1] != "--origin" && args[i-1] != "-origin" && args[i-1] != "--source" && args[i-1] != "-source" {
			log.Error("invalid parameter: " + args[i] + "\n")
			showOptions()
		}
	}

	if s.Origin == "" {
		log.Error("missing parameter: -origin\n")
		showOptions()
	} else {
		if regexp.MustCompile(`(?i)^http[s]?://[a-z0-9]`).FindString(s.Origin) == "" {
			log.Error("invalid origin url: " + s.Origin + "\n")
			showOptions()
		}

		r := regexp.MustCompile(`(?i)^http[s]?://[^/]+`)
		host := r.FindString(s.Origin)
		path := r.ReplaceAllString(s.Origin, "")

		if path == "" {
			s.Origin = s.Origin + "/"
		} else {
			stripped_path := regexp.MustCompile(`[/]?[^/]*$`).ReplaceAllString(path, "/")

			s.Origin = host + stripped_path
		}
	}

	if s.Source == "" {
		log.Error("missing parameter: -source\n")
		showOptions()
	} else {
		source, err := ioutil.ReadFile(s.Source)
		if err != nil {
			log.Error( "invalid source file: " + s.Source + " (" + err.Error() + ")\n" )
			showOptions()
		}

		s.Body = source
	}

	return s
}
