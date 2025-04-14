package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

// func render(buf string) {
// 	// Create a new colorprofile writer. We'll use it to detect the color
// 	// profile and downsample colors when necessary.
// 	w := colorprofile.NewWriter(os.Stdout, os.Environ())

// 	// While we're at it, let's jot down the detected color profile in the
// 	// markdown output while we're at it.
// 	//fmt.Fprintf(&buf, "\n\nBy the way, this was rendererd as _%s._\n", w.Profile)

// 	// Okay, now let's render some markdown.
// 	r, err := glamour.NewTermRenderer(glamour.WithEnvironmentConfig())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	md, err := r.RenderBytes([]byte(buf))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// And finally, write it to stdout using the colorprofile writer. This will
// 	// ensure colors are downsampled if necessary.
// 	fmt.Fprintf(w, "%s\n", md)
// }

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

func parse_args() []string {
	// wordPtr := flag.String("word", "foo", "a string")
	// numbPtr := flag.Int("numb", 42, "an int")
	// forkPtr := flag.Bool("fork", false, "a bool")
	// var svar string
	// flag.StringVar(&svar, "svar", "bar", "a string var")

	flag.Parse()
	return flag.Args()
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

func render_buffer(buf []byte, types []string, renderer *glamour.TermRenderer) {
	record := parse(buf)
	if record == nil {
		slog.Debug("skipping empty record")
		return
	}
	type_name := resolve_type(record)

	if types == nil || len(types) == 0 || slices.Contains(types, type_name) {
		render(type_name, record, renderer)
	} else {
		slog.Debug("skipping type", "type_name", type_name, "types", types)
	}
}

func cleanup() {
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

func main() {
	// handle the sigterm
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	// parse cli args
	names := parse_args()

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	// create a new renderer
	renderer, err := build_renderer(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// fill the buffer
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

		// output the buffer
		render_buffer(buf, names, renderer)
		time.Sleep(1 * time.Millisecond) // self throttle
	}
}
