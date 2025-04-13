package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
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
)

func apply_template(tmpl_name string, record map[string]any) (string, error) {
	if tmpl_name == "" {
		tmpl_name = "default"
	}
	tmplfile := "templates/" + tmpl_name + ".md"
	if _, err := os.Stat(tmplfile); errors.Is(err, os.ErrNotExist) {
		tmplfile = "templates/default.md"
	}
	tmpl, err := template.ParseFiles(tmplfile)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, record)
	if err != nil {
		return "", err
	}

	out, err := glamour.Render(buffer.String(), "dark")
	if err != nil {
		return "", err
	}
	return out, nil
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
			//fmt.Println(string(b))
			//fmt.Printf("map: %v\n", v.Interface())
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

func toString(m map[string]any) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%v\n", key, value)
	}
	return b.String()
}

func main() {
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

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

		var record map[string]any
		json.Unmarshal(buf, &record)
		if type_name, ok := record["type"].(string); ok {
			md, err := apply_template(type_name, record)
			if err != nil {
				log.Fatal(err)
			}
			md_strip := strings.TrimSpace(md)
			if len(md_strip) > 0 {
				fmt.Println(md_strip)
			} else {
				//fmt.Println("!!! " + toString(record))
			}
		} else {
			if len(record) > 0 {
				fmt.Println(record)
			}
		}
	}
}
