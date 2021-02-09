package controllers

import (
	"fmt"
	"github.com/praveenprem/testbed-slack-bot/helper"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"net/http"
)

func Help(w http.ResponseWriter, r *http.Request) {
	s := slack.Slack{}
	resp, err := helper.ToJson(s.Help())
	if err != nil {
		if resp, err := helper.ToJson(s.Error()); err == nil {
			_, _ = fmt.Fprint(w, resp)
		}
		return
	}
	_, _ = fmt.Fprint(w, *resp)
}
