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

func load_file(file_name string) (*template.Template, error) {
	data, err := os.ReadFile(file_name)
	if err != nil {
		return nil, err
	}
	return template.New(file_name).Parse(string(data))
}

func load_embed(type_name string) (*template.Template, error) {
	template_file := "templates/" + type_name + ".md"
	data, err := templates.ReadFile(template_file)
	if os.IsNotExist(err) {
		slog.Debug("template not found", "name", template_file)
		data, err = templates.ReadFile("templates/default.md")
	}
	if err != nil {
		return nil, err
	}
	return template.New(type_name).Parse(string(data))
}

func resolve_template(type_name string, type_params map[string]string) string {
	template_name := ""
	if len(type_params) == 0 {
		template_name = type_name
	} else {
		for key, value := range type_params {
			if strings.Contains(type_name, key) {
				if value == "" {
					template_name = type_name
				} else {
					template_name = value
				}
			}
		}
	}
	return template_name
}

func load_template(type_name string, type_params map[string]string) (*template.Template, error) {
	template_name := resolve_template(type_name, type_params)
	if strings.HasSuffix(template_name, ".md") {
		// assumes actual file
		return load_file(template_name)
	} else {
		return load_embed(template_name)
	}
}

func apply_template(tmpl_name string, type_params map[string]string, record map[string]any) (string, error) {
	if len(record) == 0 {
		return "", nil
	}
	tmpl, err := load_template(tmpl_name, type_params)
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
	} else {
		// extra-special cases
		// logs: has msg and level
		if _, ok := record["level"]; ok {
			if _, ok2 := record["msg"]; ok2 {
				type_name = "log"
			}
		}
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

func render(type_name string, record map[string]any, params map[string]string, renderer *glamour.TermRenderer) {
	md, err := apply_template(type_name, params, record)
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

func matches(type_name string, type_params map[string]string) bool {
	if len(type_params) == 0 {
		return true
	}
	for key := range type_params {
		// t contains name, name=other, name=path
		// process: contains =, split and test is path is a path
		if strings.Contains(type_name, key) {
			return true
		}
	}
	return false
}

func render_buffer(buf []byte, type_params map[string]string, renderer *glamour.TermRenderer) {
	record := parse(buf)
	if record == nil {
		slog.Debug("skipping empty record")
		return
	}
	type_name := resolve_type(record)
	if matches(type_name, type_params) {
		render(type_name, record, type_params, renderer)
	} else {
		slog.Debug("skipping type", "type_name", type_name, "types", type_params)
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
