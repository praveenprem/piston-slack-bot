package slack

import (
	"errors"
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var (
	colourError   = "#f20c44"
	colourWarning = "#f2c744"
	colourSuccess = "#44f20c"
)

type (
	Slack struct {
		Executor   string
		Args       []string
		Code       string
		Ran        bool
		Stdout     string
		Stderr     string
		Langs      []string
		RawPayload Payload
	}

	Payload struct {
		ApiAppId            string `schema:"api_app_id"`
		ChannelId           string `schema:"channel_id"`
		ChannelName         string `schema:"channel_name"`
		Command             string `schema:"command"`
		IsEnterpriseInstall string `schema:"is_enterprise_install"`
		ResponseUrl         string `schema:"response_url"`
		TeamDomain          string `schema:"team_domain"`
		TeamId              string `schema:"team_id"`
		Text                string `schema:"text"`
		Token               string `schema:"token"`
		TriggerId           string `schema:"trigger_id"`
		UserId              string `schema:"user_id"`
		UserName            string `schema:"user_name"`
		EnterpriseId        string `schema:"enterprise_id"`
		EnterpriseName      string `schema:"enterprise_name"`
	}

	Slacker interface {
		Help() *Message
		Languages() *Message
		Output() *Message
		Error() *Message
		execSuccess() []attachment
		execFail() []attachment
		multiFieldSets() []block
	}

	Payloader interface {
		Parser(r *http.Request) error
		Validate() error
	}

	Config struct {
		AppId        string `json:"appId"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		SigningSig   string `json:"signingSig"`
	}

	Message struct {
		ResponseType string       `json:"response_type"`
		Attachments  []attachment `json:"attachments"`
	}

	blocks struct {
		Blocks []block `json:"blocks"`
	}

	block struct {
		Type string `json:"type,omitempty"`
		Text *text  `json:"text,omitempty"`
		*fields
		*elements
	}

	text struct {
		Type  string `json:"type"`
		Text  string `json:"text"`
		Emoji *bool  `json:"emoji,omitempty"`
	}

	fields struct {
		Fields []text `json:"fields"`
	}

	elements struct {
		Elements *[]text `json:"elements,omitempty"`
	}

	attachment struct {
		Color string `json:"color"`
		*blocks
	}
)

var (
	CFG Config
)

func (s *Slack) Help() *Message {
	log.Printf("responding to help instructions")
	var body = Message{
		ResponseType: "ephemeral",
		Attachments: []attachment{
			{
				Color: colourWarning,
				blocks: &blocks{
					[]block{
						{
							Type: "header",
							Text: &text{
								Type: "plain_text",
								Text: "Usage /tb-lang",
							},
						},
						{
							Type: "section",
							fields: &fields{
								[]text{
									{
										Type: "plain_text",
										Text: "python3",
									},
									{
										Type: "plain_text",
										Text: "python2",
									},
									{
										Type: "plain_text",
										Text: "rust",
									},
									{
										Type: "plain_text",
										Text: "typescript",
									},
									{
										Type: "plain_text",
										Text: "php",
									},
									{
										Type: "plain_text",
										Text: "paradoc",
									},
									{
										Type: "plain_text",
										Text: "go",
									},
									{
										Type: "plain_text",
										Text: "java",
									},
								},
							},
						},
						{
							Type: "divider",
						},
						{
							Type: "header",
							Text: &text{
								Type: "plain_text",
								Text: "Usage /tb",
							},
						},
						{
							Type: "section",
							Text: &text{
								Type: "mrkdwn",
								Text: "/tb [lang]\n```code```\n[arg1] [arg2] [arg3]...",
							},
						},
						{
							Type: "context",
							elements: &elements{
								Elements: &[]text{
									{
										Type: "mrkdwn",
										Text: "*Sample Python test*",
									},
								},
							},
						},
						{
							Type: "context",
							elements: &elements{
								Elements: &[]text{
									{
										Type: "mrkdwn",
										Text: "Execute",
									},
									{
										Type: "mrkdwn",
										Text: "/tb-exec python3\n```from sys import argv\nprint(f'Hello {argv[1]} {argv[2]}')```\nJohn Smith",
									},
								},
							},
						},
						{
							Type: "context",
							elements: &elements{
								Elements: &[]text{
									{
										Type: "mrkdwn",
										Text: "Output",
									},
									{
										Type: "mrkdwn",
										Text: ">Stdout\n```Hello John Smith```\n>Stderr \n\n",
									},
								},
							},
						},
						footer(),
					},
				},
			},
		},
	}
	return &body
}

func (s *Slack) Languages() *Message {
	log.Printf("responding to language options request")
	blks := []block{
		{
			Type: "header",
			Text: &text{
				Type: "plain_text",
				Text: "Supported languages",
			},
		},
	}

	for _, b := range s.multiFieldSets() {
		blks = append(blks, b)
	}

	blks = append(blks, footer())

	var body = Message{
		ResponseType: "ephemeral",
		Attachments: []attachment{
			{
				Color: colourSuccess,
				blocks: &blocks{
					blks,
				},
			},
		},
	}

	return &body
}

func (s *Slack) Output() *Message {
	log.Printf("generating output")
	body := Message{Attachments: []attachment{}}
	if s.Stderr != "" {
		body.ResponseType = "ephemeral"
		body.Attachments = s.execFail()
	} else {
		body.ResponseType = "in_channel"
		body.Attachments = s.execSuccess()
	}
	return &body
}

func (s *Slack) Error() *Message {
	log.Printf("sending generic error message")
	return &Message{
		Attachments: []attachment{
			{
				Color: colourError,
				blocks: &blocks{
					[]block{
						{
							Type: "section",
							Text: &text{
								Type: "mrkdwn",
								Text: "*`Something went wrong...`*\nTry again later...",
							},
						},
						footer(),
					},
				},
			},
		},
	}
}

func (s *Slack) execSuccess() []attachment {
	log.Printf("generating success message")
	body := []attachment{
		{
			Color: colourSuccess,
			blocks: &blocks{
				[]block{
					{
						Type: "section",
						Text: &text{
							Type: "plain_text",
							Text: "Test Bed execution results",
						},
					},
					{
						Type: "context",
						elements: &elements{
							Elements: &[]text{
								{
									Type: "mrkdwn",
									Text: "Execute",
								},
								{
									Type: "mrkdwn",
									Text: fmt.Sprintf("%s %s", s.RawPayload.Command, removeNewLineEscape(s.RawPayload.Text)),
								},
							},
						},
					},
					{
						Type: "context",
						elements: &elements{
							Elements: &[]text{
								{
									Type: "mrkdwn",
									Text: "Output",
								},
								{
									Type: "mrkdwn",
									Text: fmt.Sprintf(">*Stdout*\n```%s```\n>_Stderr_\n\n", s.Stdout),
								},
							},
						},
					},
					footer(),
				},
			},
		},
	}
	return body
}

func (s *Slack) execFail() []attachment {
	log.Printf("generating error message")
	body := []attachment{
		{
			Color: colourError,
			blocks: &blocks{
				[]block{
					{
						Type: "header",
						Text: &text{
							Type: "plain_text",
							Text: "Error occurred",
						},
					},
					{
						Type: "section",
						Text: &text{
							Type: "mrkdwn",
							Text: "`Something went wrong, please validate the provided code below`",
						},
					},
					{
						Type: "section",
						Text: &text{
							Type: "mrkdwn",
							Text: fmt.Sprintf("%s %s", s.RawPayload.Command, removeNewLineEscape(s.RawPayload.Text)),
						},
					},
					{
						Type: "context",
						elements: &elements{
							Elements: &[]text{
								{
									Type: "mrkdwn",
									Text: "Output",
								},
								{
									Type: "mrkdwn",
									Text: fmt.Sprintf(">*Stdout*\n```null```\n>_Stderr_\n```%s```", s.Stderr),
								},
							},
						},
					},
					footer(),
				},
			},
		},
	}
	return body
}

func (s *Slack) multiFieldSets() []block {

	var (
		blks []block
	)

	blk := block{
		Type:   "section",
		fields: &fields{Fields: []text{}},
	}

	sort.Strings(s.Langs)

	counter := 0
	for i, lang := range s.Langs {
		if counter == 10 {
			blks = append(blks, blk)
			blk = block{Type: "section", fields: &fields{Fields: []text{}}}
			counter = 0
		}

		txt := text{
			Type: "plain_text",
			Text: lang,
		}

		blk.fields.Fields = append(blk.fields.Fields, txt)
		counter++

		if i == len(s.Langs)-1 {
			blks = append(blks, blk)
		}

	}

	return blks
}

func (p *Payload) Parser(r *http.Request) error {
	log.Printf("parsing slack request body")
	if err := schema.NewDecoder().Decode(p, r.Form); err != nil {
		log.Printf("%#v", err)
		return err
	}

	p.Token = strings.TrimRight(p.Token, "\n")
	p.Command = strings.TrimRight(p.Command, "\n")
	p.Text = strings.TrimRight(strings.ReplaceAll(p.Text, "\n", `\n`), `\n`)
	p.ResponseUrl = strings.TrimRight(p.ResponseUrl, "\n")
	p.TriggerId = strings.TrimRight(p.TriggerId, "\n")
	p.UserId = strings.TrimRight(p.UserId, "\n")
	p.UserName = strings.TrimRight(p.UserName, "\n")
	p.TeamId = strings.TrimRight(p.TeamId, "\n")
	p.EnterpriseId = strings.TrimRight(p.EnterpriseId, "\n")
	p.ChannelId = strings.TrimRight(p.ChannelId, "\n")
	p.ChannelName = strings.TrimRight(p.ChannelName, "\n")
	p.ApiAppId = strings.TrimRight(p.ApiAppId, "\n")
	p.EnterpriseName = strings.TrimRight(p.EnterpriseName, "\n")
	p.TeamDomain = strings.TrimRight(p.TeamDomain, "\n")

	log.Printf("form data parse completed")

	return nil
}

func (p *Payload) Validate(slack *Slack) error {
	log.Printf("validating form fields")
	regex := regexp.MustCompile(`^(.*)\\n\x60\x60\x60(.*)\x60\x60\x60\\?n?(.*)$`)
	matches := regex.FindAllStringSubmatch(p.Text, -1)
	if len(matches) < 1 {
		log.Printf("validation failed")
		log.Println(p.Text)
		return errors.New("message did not meet the expectation")
	}

	slack.Executor = matches[0][1]
	slack.Code = strings.ReplaceAll(matches[0][2], "\\n", "\n") //matches[0][3]
	if len(matches[0]) == 4 {
		slack.Args = strings.Split(matches[0][3], " ")
	}

	if slack.Executor == "" {
		slack.Stderr = "Malformed command. Missing execution language"
		return errors.New(fmt.Sprintf("Malformed command. Missing execution language"))
	}
	if slack.Code == "" {
		slack.Stderr = "Malformed command. Missing execution code"
		return errors.New(fmt.Sprintf("Malformed command. Missing execution code"))
	}

	log.Printf("validation completed")

	return nil
}

func removeNewLineEscape(text string) string {
	return strings.ReplaceAll(text, "\\n", "\n")

}

func footer() block {
	return block{
		Type: "context",
		elements: &elements{
			Elements: &[]text{
				{
					Type: "mrkdwn",
					Text: "Powered by <https://github.com/engineer-man/piston-bot|Piston - Code Execution engine>",
				},
			},
		},
	}
}
