package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

//go:embed templates
var templates embed.FS

var type_uses = map[string]int{}

func inc_type_use(type_name string) {
	type_uses[type_name] += 1
}

func load_template(type_name string) (*template.Template, error) {
	if type_name == "" {
		type_name = "default"
	}
	template_file := "templates/" + type_name + ".md"
	data, err := templates.ReadFile(template_file)
	if os.IsNotExist(err) {
		slog.Debug("template not found", "name", type_name)
		data, err = templates.ReadFile("templates/default.md")
	}
	if err != nil {
		return nil, err
	}
	return template.New(type_name).Parse(string(data))
}

func apply_template(tmpl_name string, record map[string]any) (string, error) {
	if len(record) == 0 {
		return "", nil
	}
	tmpl, err := load_template(tmpl_name)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, record)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func parse(buf []byte) map[string]any {
	if len(buf) == 0 {
		return nil
	}
	var record map[string]any
	json.Unmarshal(buf, &record)
	return record
}

func resolve_type(record map[string]any) string {
	// check type and record.$type
	var type_name string = ""
	if r, ok := record["record"]; ok {
		rec := r.(map[string]any)
		type_name = rec["$type"].(string)
	} else if t, ok := record["type"]; ok {
		type_name = t.(string)
	} else if a, ok := record["action"]; ok {
		type_name = a.(string)
	}
	if type_name == "" {
		// type not found
		slog.Info("type not found", "record", record)
		inc_type_use("[empty]")
		return "default"
	} else {
		inc_type_use(type_name)
		return type_name
	}
}

func build_renderer(stream *os.File) (*glamour.TermRenderer, error) {
	// get terminal witdh for stdout
	width, _, err := term.GetSize(int(stream.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	// Create a new renderer.
	return glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
}

func render(type_name string, record map[string]any, renderer *glamour.TermRenderer) {
	md, err := apply_template(type_name, record)
	if err != nil {
		log.Fatal(err)
	}
	rendered, err := renderer.Render(md)
	if err != nil {
		log.Fatal(err)
	}
	trimmed := strings.TrimSpace(rendered)
	if len(trimmed) > 0 {
		fmt.Println(trimmed)
	}
}

func matches(type_name string, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if strings.Contains(type_name, t) {
			return true
		}
	}
	return false
}

func render_buffer(buf []byte, types []string, renderer *glamour.TermRenderer) {
	record := parse(buf)
	if record == nil {
		slog.Debug("skipping empty record")
		return
	}
	type_name := resolve_type(record)
	if matches(type_name, types) {
		render(type_name, record, renderer)
	} else {
		slog.Debug("skipping type", "type_name", type_name, "types", types)
	}
}

func type_report() {
	types_map := make(map[string]any)
	// TODO sort
	for k, v := range type_uses {
		types_map[k] = strconv.Itoa(v)
	}
	fmt.Println()
	for type_name, count := range type_uses {
		fmt.Println(type_name, count)
	}
}
