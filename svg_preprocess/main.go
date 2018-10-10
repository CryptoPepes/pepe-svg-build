package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"path/filepath"
	"os"
	"encoding/xml"
	"path"
	"gopkg.in/yaml.v2"
	"strings"
	"bytes"
	"golang.org/x/net/html/charset"
	"encoding/json"
)

type Config struct {
	Categories map[string][]Part `yaml:"parts"`
}

type Part struct {
	Front string `yaml:"front"`
	Back string `yaml:"back"`
	Id string `yaml:"id"`
	Name string `yaml:"name"`
	Colors map[string]string `yaml:"colors"`
}

type JsonMappings struct {
	Id2Name map[string]string `json:"id2name"`
	Name2Id map[string]string `json:"name2id"`
}

func main() {
	config := new(Config)

	svgDir := path.Dir("pepe-svg-art/")
	outDir := path.Dir("builder/tmpl/")

	if err := os.RemoveAll(outDir); err != nil {
		log.Fatalf("Failed to clean output dir! %v", err)
	}

	if err := os.Mkdir(outDir, os.ModeDir | os.ModePerm); err != nil {
		log.Fatalf("Failed to create output dir! %v", err)
	}

	configFile, err := ioutil.ReadFile(path.Join(svgDir, "parts.yml"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal([]byte(configFile), &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	files, err := filepath.Glob(path.Join(svgDir, "**/*.svg"))
	if err != nil {
		log.Fatal(err)
	}

	parsedSvgs := map[string]*SvgFile{}
	for _, f := range files {
		log.Println("Parsing SVG: "+f)
		parsed := parseSVG(f)
		if parsed != nil {
			parsedSvgs[f] = parsed
		}
	}

	// Create a display-name mapping for the client to use
	jsonOut := JsonMappings{Id2Name: make(map[string]string), Name2Id: make(map[string]string)}

	for categoryName, category := range config.Categories {

		catOutDir := path.Join(outDir, categoryName)
		if err := os.Mkdir(catOutDir, os.ModeDir | os.ModePerm); err != nil {
			log.Fatalf("Failed to create output dir! %v", err)
		}

		for _, part := range category {
			outFilePath := path.Join(catOutDir, part.Id + ".tmpl")

			fullId := categoryName + ">" + part.Id
			jsonOut.Id2Name[fullId] = part.Name
			// Save name->id with lowercase keys, to enable case-insensitive lookups.
			lowercaseName := strings.ToLower(part.Name)
			jsonOut.Name2Id[lowercaseName] = fullId

			log.Println("Creating template: "+outFilePath)
			var svgFront, svgBack *SvgGroup
			if part.Front != "" {
				inputFrontFilePath := path.Join(svgDir, categoryName, part.Front)
				svgFront = &parsedSvgs[inputFrontFilePath].Groups[0]
			}
			if part.Back != "" {
				inputBackFilePath := path.Join(svgDir, categoryName, part.Back)
				svgBack = &parsedSvgs[inputBackFilePath].Groups[0]
			}
			tmplStr := createTemplate(fullId, part, svgFront, svgBack)
			err := ioutil.WriteFile(outFilePath, []byte(tmplStr), os.ModePerm)
			if err != nil {
				log.Println("Error! Failed to write! "+outFilePath)
				log.Println(err)
				return
			}
		}
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(jsonOut); err != nil {
		log.Println("Error! Failed to convert mappings to JSON!")
		log.Println(err)
		return
	}
	if err := ioutil.WriteFile(path.Join(outDir, "mappings.json"), buf.Bytes(), os.ModePerm);
		err != nil {
		log.Println("Error! Failed to write mappings file!")
		log.Println(err)
		return
	}

	log.Println("Done!")
}

type SvgFile struct {
	XMLName xml.Name `xml:"svg"`

	Groups []SvgGroup `xml:"g"`
}

type SvgGroup struct {
	XMLName xml.Name `xml:"g"`

	Id string `xml:"id,attr"`
	XML string `xml:",innerxml"`
}

func parseSVG(path string) *SvgFile {
	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer xmlFile.Close()

	fileContents, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	var parsedContent SvgFile

	decoder := xml.NewDecoder(bytes.NewReader(fileContents))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&parsedContent); err != nil {
		fmt.Println("Error parsing XML:", err)
		return nil
	}

	if len(parsedContent.Groups) < 1 {
		fmt.Println("Unexpected SVG problem: no svg groups! file: "+path)
		return nil
	}
	return &parsedContent
}

// Creates a golang template, to easily replace colors etc. with in production.
func createTemplate(fullId string, part Part, frontSvg *SvgGroup, backSvg *SvgGroup) (out string) {
	out = `
{{ define "` + fullId + `>back" }}
`
	if backSvg != nil {
		out += backSvg.XML
	}
	out += `
{{ end }}
`
	out += `
{{ define "` + fullId + `>front" }}
`
	if frontSvg != nil {
		out += frontSvg.XML
	}
	out += `
{{ end }}
`
	if part.Colors != nil {
		for k, v := range part.Colors {
			out = strings.Replace(out, "fill:" + v, "fill:{{- ."+k+" -}}", -1)
			out = strings.Replace(out, "fill: " + v, "fill:{{- ."+k+" -}}", -1)
		}
	}
	return out
}
