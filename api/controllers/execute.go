package controllers

import (
	"fmt"
	"github.com/praveenprem/testbed-slack-bot/helper"
	"github.com/praveenprem/testbed-slack-bot/piston"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"log"
	"net/http"
)

func Execute(w http.ResponseWriter, r *http.Request) {
	var slk slack.Slack

	if err := r.ParseForm(); err != nil {
		resp, err := helper.ToJson(slk.Help())
		if err != nil {
			if resp, err := helper.ToJson(slk.Error()); err == nil {
				_, _ = fmt.Fprint(w, resp)
			}
			return
		}
		_, _ = fmt.Fprint(w, resp)
		return
	}

	log.Printf("%#v", r.Form)

	if err := slk.RawPayload.Parser(r); err != nil {
		log.Printf("%#v", err)
		slk.Stderr = err.Error()
		resp, err := helper.ToJson(slk.Output())
		if err != nil {
			if resp, err := helper.ToJson(slk.Error()); err == nil {
				log.Printf("%#v", err)
				_, _ = fmt.Fprint(w, resp)
			} else {
				log.Printf("%#v", err)
			}
			return
		}
		_, _ = fmt.Fprint(w, *resp)
		return
	}

	if err := slk.RawPayload.Validate(&slk); err != nil {
		log.Printf("%#v", err)
		slk.Stderr = err.Error()
		resp, err := helper.ToJson(slk.Output())
		if err != nil {
			if resp, err := helper.ToJson(slk.Error()); err == nil {
				log.Printf("%#v", err)
				_, _ = fmt.Fprint(w, resp)
			} else {
				log.Printf("%#v", err)
			}
			return
		}
		_, _ = fmt.Fprint(w, *resp)
		return
	}

	pstn := piston.Piston{
		Execute: &piston.Execute{},
	}
	pstn.Language = slk.Executor
	pstn.Args = slk.Args
	pstn.Source = slk.Code
	log.Printf("%#v", pstn.Execute)
	if err := pstn.Exec(); err != nil {
		log.Printf("%#v", err)
		slk.Stderr = err.Error()
		slk.Stdout = ""
		slk.Ran = pstn.Response.Ran
		if resp, err := helper.ToJson(slk.Output()); err != nil {
			log.Printf("%#v", err)
		} else {
			_, _ = fmt.Fprint(w, resp)
		}
		return
	}

	slk.Ran = pstn.Response.Ran
	slk.Stdout = pstn.Response.Stdout
	slk.Stderr = pstn.Response.Stderr
	resp, msgBuildErr := helper.ToJson(slk.Output())
	if msgBuildErr != nil {
		log.Printf("%#v", msgBuildErr)
		slk.Stderr = msgBuildErr.Error()
		resp, err := helper.ToJson(slk.Output())
		if err != nil {
			if resp, err := helper.ToJson(slk.Error()); err == nil {
				log.Printf("%#v", err)
				_, _ = fmt.Fprint(w, resp)
			} else {
				log.Printf("%#v", err)
			}
			return
		}
		_, _ = fmt.Fprint(w, *resp)
		return
	}

	_, _ = fmt.Fprint(w, *resp)

}
