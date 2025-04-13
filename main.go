package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

//go:embed templates
var templates embed.FS

func load_template(type_name string) ([]byte, error) {
	if type_name == "" {
		type_name = "default"
	}
	template_file := "templates/" + type_name + ".md"
	return templates.ReadFile(template_file)
}

func apply_template(tmpl_name string, record map[string]any) (string, error) {
	if len(record) == 0 {
		return "", nil
	}
	if tmpl_name == "" {
		tmpl_name = "default"
	}
	tmplfile := "templates/" + tmpl_name + ".md"
	tmpl, err := template.ParseFiles(tmplfile)
	if err != nil {
		// assume missing file, no way to stat the embed
		tmpl, err = template.ParseFiles("templates/default.md")
		if err != nil {
			return "", err
		}
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, record)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
	// out, err := glamour.Render(buffer.String(), "dark")
	// if err != nil {
	// 	return "", err
	// }
	// return strings.TrimSpace(out), nil
}

func render(buf string) {
	// Create a new colorprofile writer. We'll use it to detect the color
	// profile and downsample colors when necessary.
	w := colorprofile.NewWriter(os.Stdout, os.Environ())

	// While we're at it, let's jot down the detected color profile in the
	// markdown output while we're at it.
	//fmt.Fprintf(&buf, "\n\nBy the way, this was rendererd as _%s._\n", w.Profile)

	// Okay, now let's render some markdown.
	r, err := glamour.NewTermRenderer(glamour.WithEnvironmentConfig())
	if err != nil {
		log.Fatal(err)
	}
	md, err := r.RenderBytes([]byte(buf))
	if err != nil {
		log.Fatal(err)
	}

	// And finally, write it to stdout using the colorprofile writer. This will
	// ensure colors are downsampled if necessary.
	fmt.Fprintf(w, "%s\n", md)
}

func render_table(record map[string]any) {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	t := table.New()
	for key, value := range record {
		//fmt.Println(key, value.(string))
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Map:
			for subkey, subvalue := range value.(map[string]any) {
				b, err := json.Marshal(subvalue)
				if err != nil {
					log.Fatal(err)
				}
				t.Row(s(key+"."+subkey), string(b))
			}
		case reflect.Slice:
			b, err := json.Marshal(value)
			if err != nil {
				log.Fatal(err)
			}
			t.Row(key, string(b))

		default:
			b, err := json.Marshal(value)
			if err != nil {
				log.Fatal(err)
			}
			t.Row(s(key), s(string(b)))
		}
	}
	fmt.Println(t.Render())
}

func toString(record map[string]any) string {
	b := new(bytes.Buffer)
	for key, value := range record {
		fmt.Fprintf(b, "%s=%v\n", key, value)
	}
	return b.String()
}

func parse(buf []byte) map[string]any {
	var record map[string]any
	json.Unmarshal(buf, &record)
	return record
}

func resolve_type(record map[string]any) string {
	// check type and record.$type
	if t, ok := record["type"]; ok {
		return t.(string)
	} else if a, ok := record["action"]; ok {
		return a.(string)
	} else if r, ok := record["record"]; ok {
		rec := r.(map[string]any)
		return rec["$type"].(string)
	}
	return "default"
}

func main() {
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	//println("width:", width, "height:", height)

	// Create a new renderer.
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating renderer: %s\n", err)
	}

	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err) // FIXME
		}

		record := parse(buf)
		type_name := resolve_type(record)
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
}
