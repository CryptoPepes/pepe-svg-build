package builder

import (
	"bytes"
	"html/template"
	"io"
	"path"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/svg"
	"cryptopepe.io/cryptopepe-svg/builder/look"
	"log"
	"os"
	"runtime"
)

type SVGBuilder struct {
	svgBuildTemplate *template.Template
	mini *minify.M
}

// Loads the SVG builder.
func (builder *SVGBuilder) Load() {
	t := template.New("builder")

	// Hack to make templates dynamically accessible by name.
	t.Funcs(map[string]interface{}{
		"CallTemplate": func(name string, data interface{}, placement string) (ret template.HTML, err error) {
			buf := bytes.NewBuffer([]byte{})
			err = t.ExecuteTemplate(buf, name + ">" + placement, data)
			if err != nil {
				buf.Reset()
				log.Println("Warning! Using fallback, failed to find pepe component: ", name, placement, err)
				err = t.ExecuteTemplate(buf, "error>fallback", data)
			}
			ret = template.HTML(buf.String())
			return
		},
	})

	basePath, envSetBasePath := os.LookupEnv("APP_PATH")
	if !envSetBasePath {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			panic("No caller info! Cannot get file path for SVG template loading.")
		}
		basePath = path.Dir(filename)
	}

	//add main builder template
	t = template.Must(t.ParseFiles(path.Join(basePath, "builder.tmpl")))
	//add all svg templates
	builder.svgBuildTemplate = template.Must(t.ParseGlob(basePath+"/tmpl/**/*.tmpl"))


	builder.mini = minify.New()
	builder.mini.Add("image/svg+xml", svg.DefaultMinifier)
}

// Converts to SVG, and writes minified result to `w`
func (builder *SVGBuilder) ConvertToSVG(w io.Writer, look *look.PepeLook) error {
	buf := new(bytes.Buffer)
	if err := builder.svgBuildTemplate.Execute(buf, look); err != nil {
		return err
	}
	minifyingW := builder.mini.Writer("image/svg+xml", w)
	if _, err := buf.WriteTo(minifyingW); err != nil {
		return err
	}
	if err := minifyingW.Close(); err != nil {
		return err
	}
	return nil
}
