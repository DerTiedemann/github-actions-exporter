package metrics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v38/github"

	"github-actions-exporter/pkg/config"
)

var (
	workflows map[string]map[int64]github.Workflow
)

// workflowCache - used for limit calls to github api
func workflowCache() {
	for {
		cache, err := buildWorkflowMap()
		if err != nil {
			log.Println(err)
		} else {
			workflows = cache
		}

		time.Sleep(time.Duration(60) * time.Second)
	}
}

func buildWorkflowMap() (map[string]map[int64]github.Workflow, error) {

	ww := make(map[string]map[int64]github.Workflow)

	for _, repo := range config.Github.Repositories.Value() {
		r := strings.Split(repo, "/")
		s := make(map[int64]github.Workflow)
		opt := &github.ListOptions{PerPage: int(config.Github.PageSize)}

		for {
			resp, rr, err := client.Actions.ListWorkflows(context.Background(), r[0], r[1], opt)
			if err != nil {
				return nil, fmt.Errorf("ListWorkflows error for %s: %s", repo, err.Error())
			}
			for _, w := range resp.Workflows {
				s[*w.ID] = *w
			}

			if rr.NextPage == 0 {
				break
			}
			opt.Page = rr.NextPage
		}

		ww[repo] = s

	}

	return ww, nil

}
