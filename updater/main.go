package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/mcuadros/go-version"
	"golang.org/x/net/html"
)

var findEnvParameter = regexp.MustCompile(`ENV(\s+([^=\s]+)=([^\s=]+)(\s\\)?)+`)
var findImage = regexp.MustCompile(`(?P<pre>[ \t]*image:[ \t]*)(?P<image>\S+):(?P<tag>\S+)`)
var envParameter = regexp.MustCompile(`(?P<key>[^=\s]+)=(?P<value>[^\s=]+)`)

const url = "https://www.factorio.com/download-headless/experimental"

func main() {
	composefilePath := flag.String("composefile", "docker-compose.yml", "Path to compose file if not same directory")
	dockerfilePath := flag.String("dockerfile", "Dockerfile", "Path to dockerfile if not same directory")
	forceUpdate := flag.Bool("force", false, "Force the update of the files")
	flag.Parse()

	fmt.Println("- Checking for update...")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	versions := getElementsByTag(doc, "h3")
	latestVersion := ""

	re := regexp.MustCompile(`(\d+\.\d+\.\d+)`)

	for _, node := range versions {
		v := re.FindString(node.Data)

		if latestVersion == "" || version.CompareSimple(v, latestVersion) > 0 {
			latestVersion = v
		}
	}

	envParamter := getEnvParameterFromDockerfile(*dockerfilePath)
	fmt.Println("  Current version: " + envParamter["VERSION"])
	fmt.Println("  Latest version:  " + latestVersion)

	if version.CompareSimple(latestVersion, envParamter["VERSION"]) <= 0 && !*forceUpdate {
		fmt.Println()
		fmt.Println("- No update available")
		os.Exit(1)
	}

	if *forceUpdate {
		fmt.Println()
		fmt.Println("- Forced update")
	}

	fmt.Println()
	fmt.Println("- Retrieving hash...")
	respfile, err := http.Get(fmt.Sprintf("https://www.factorio.com/get-download/%s/headless/linux64", latestVersion))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer respfile.Body.Close()

	h := sha1.New()
	if _, err := io.Copy(h, respfile.Body); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sha1 := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("  SHA1: " + sha1)

	fmt.Println()
	fmt.Println("- Updating Dockerfile...")
	updateDockerfile(*dockerfilePath, envParamter, latestVersion, sha1)

	fmt.Println("- Updating docker-compose file...")
	updateComposefile(*composefilePath, latestVersion)

	fmt.Println("- Finished")
	os.Exit(0)
}

func getElementsByTag(doc *html.Node, tag string) (nodes []*html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == tag {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nodes
}

func getEnvParameterFromDockerfile(filepath string) map[string]string {
	dockerfile, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	dockerfileStr := string(dockerfile)

	parameter := make(map[string]string)
	env := findEnvParameter.FindString(dockerfileStr)
	for _, match := range envParameter.FindAllStringSubmatch(env, -1) {
		parameter[match[1]] = match[2]
	}

	return parameter
}

func updateComposefile(filepath string, newVersion string) {
	composefile, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	composefileStr := string(composefile)

	match := findImage.FindStringSubmatch(composefileStr)
	if len(match) == 0 {
		return
	}

	composefileStr = strings.Replace(composefileStr, match[0], match[1]+match[2]+":"+newVersion, 1)
	ioutil.WriteFile(filepath, []byte(composefileStr), 0666)
}

func updateDockerfile(filepath string, parameter map[string]string, newVersion string, hash string) {
	dockerfile, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	dockerfileStr := string(dockerfile)

	env := findEnvParameter.FindString(dockerfileStr)
	envNew := strings.Replace(env, parameter["VERSION"], newVersion, 1)
	envNew = strings.Replace(envNew, parameter["SHA1"], hash, 1)
	dockerfileStr = strings.Replace(dockerfileStr, env, envNew, 1)

	ioutil.WriteFile(filepath, []byte(dockerfileStr), 0666)
}
