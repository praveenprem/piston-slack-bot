package controllers

import (
	"fmt"
	"github.com/praveenprem/testbed-slack-bot/slack"
	"log"
	"net/http"
)

var badPage = `
<html>
	<body>
		<div>
			<h1 align="center">
    			<img src="https://raw.githubusercontent.com/praveenprem/testbed-slack-bot/master/images/code.png" width="25" height="25" alt="Bot icon">Test Bed
			</h1>
		</div>
		<div align="center">
			<p>
				<b>Something went wrong on the OAuth process!</b>
			</p>
			<p>Please try again...</p>
		</div>
	</body>
</html>
`

func Authorise(w http.ResponseWriter, r *http.Request) {
	var auth slack.Auth

	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("missing auth code")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html")
		log.Printf("%#v", w.Header())
		_, _ = fmt.Fprint(w, badPage)
		return
	}

	if err := auth.Authorise(code); err != nil {
		log.Printf("%#v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html")
		log.Printf("%#v", w.Header())
		_, _ = fmt.Fprint(w, badPage)
		return
	}

	log.Printf("authorisation completed")
	log.Printf("redirecting to %s", auth.Team.Name)
	http.Redirect(w, r, fmt.Sprintf("https://app.slack.com/client/%s", auth.Team.Id), http.StatusSeeOther)
}
