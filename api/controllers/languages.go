package controllers

import (
	"fmt"
	"github.com/praveenprem/testbed-slack-bot/helper"
	"github.com/praveenprem/testbed-slack-bot/piston"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"log"
	"net/http"
)

func Languages(w http.ResponseWriter, r *http.Request) {
	var (
		slk  slack.Slack
		pstn piston.Piston
	)

	if err := pstn.Lang(); err != nil {
		log.Printf("%#v", err)
		w.WriteHeader(http.StatusBadRequest)
	}


	for _, l := range *pstn.Versions {
		slk.Langs = append(slk.Langs, l.Name)
	}

	resp, err := helper.ToJson(slk.Languages())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, *resp)
}
