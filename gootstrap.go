package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type templateDescriptor struct {
	dirs  []string
	files []string
}

func getTemplatesToOutput(
	templatesdir string,
	outputdir string,
) map[string]string {
	templateDesc := newTemplateDescriptor(templatesdir)
	fmt.Println(templateDesc)
	prepareDirectories(outputdir, templateDesc.dirs)
	return nil
}

func relativePath(path string) (string, bool) {
	parsedPath := strings.Split(path, "/templates/")
	if len(parsedPath) > 1 {
		return parsedPath[len(parsedPath)-1], true
	}

	return "", false
}

func newTemplateDescriptor(dir string) templateDescriptor {
	desc := templateDescriptor{}
	if d, ok := relativePath(dir); ok {
		desc.dirs = append(desc.dirs, d)
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if path == dir {
				return err
			}
			recursiveDesc := newTemplateDescriptor(path)
			desc.dirs = append(desc.dirs, recursiveDesc.dirs...)
			desc.files = append(desc.files, recursiveDesc.files...)
		} else {
			desc.files = append(desc.files, path)
		}
		return err
	})

	return desc
}

func prepareDirectories(output string, dirs []string) {
	// TODO: check for empty dirs, create if do not exists
}

func applyTemplates(
	templateToOut map[string]string,
	project string,
	docker_registry string,
) {
}

func fatalerr(err error, context string) {
	if err != nil {
		fmt.Printf("error[%s] while[%s]\n", err, context)
		os.Exit(1)
	}
}

func getGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return os.Getenv("HOME") + "/go"
}

func getTemplatesDir() string {
	templatesdir := getGoPath() + "/src/github.com/NeowayLabs/gootstrap/templates"
	if _, err := os.Stat(templatesdir); os.IsNotExist(err) {
		errmsg := fmt.Sprintf(
			"templates directory[%s] do not exist, probably "+
				"a bad/corrupted gootstrap installation\n"+
				"are you sure you installed with 'go get' ?",
			templatesdir,
		)
		panic(fmt.Sprintf("%s\nerror[%s]", errmsg, err))
	}
	return templatesdir
}

func requiredOption(v string, msg string) {
	if v == "" {
		fmt.Println(msg + ":\n")
		flag.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}
}

func guaranteeEmptyDir(outputdir string) {

	info, err := os.Stat(outputdir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(outputdir, 0775)
		fatalerr(err, fmt.Sprintf("creating output dir[%s]", outputdir))
		return
	}

	if !info.IsDir() {
		fmt.Printf("outputdir[%s] already exists and it is not a directory\n", outputdir)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(outputdir)
	fatalerr(err, fmt.Sprintf("listing outputdir[%s]", outputdir))

	if len(files) > 0 {
		fmt.Printf("outputdir[%s] already exists and it is not empty\n", outputdir)
		os.Exit(1)
	}
}

func main() {
	outputdir := ""
	project := ""
	docker_registry := ""

	cwd, err := os.Getwd()
	fatalerr(err, "getting default cwd")
	templatesdir := getTemplatesDir()

	flag.StringVar(&project, "project", "", "project name (required)")
	flag.StringVar(&docker_registry, "docker-registry", "", "docker registry base name (required)")
	flag.StringVar(&outputdir, "output", cwd, "output directory where project will be created (optional)")
	flag.Parse()

	requiredOption(project, "'project' is a required option")
	requiredOption(docker_registry, "'docker-registry' is a required option")

	fmt.Printf("\ncreating: project[%s] docker-registry[%s] output[%s]\n\n", project, docker_registry, outputdir)

	outputdir += "/" + project

	guaranteeEmptyDir(outputdir)
	templateToOut := getTemplatesToOutput(templatesdir, outputdir)
	applyTemplates(templateToOut, project, docker_registry)
}
