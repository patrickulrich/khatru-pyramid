package main

import (
	"fmt"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/theplant/htmlgo"
)

func inviteTreeHandler(w http.ResponseWriter, r *http.Request) {
	content := inviteTreePageHTML(r.Context(), InviteTreePageParams{
		loggedUser: getLoggedUser(r),
	})
	htmlgo.Fprint(w, baseHTML(content), r.Context())
}

func addToWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)

	pubkey := r.PostFormValue("pubkey")
	if pfx, value, err := nip19.Decode(pubkey); err == nil && pfx == "npub" {
		pubkey = value.(string)
	}

	if loggedUser != s.RelayPubkey && hasInvitedAtLeast(loggedUser, s.MaxInvitesPerPerson) {
		http.Error(w, fmt.Sprintf("cannot invite more than %d", s.MaxInvitesPerPerson), 403)
		return
	}

	if err := addToWhitelist(pubkey, loggedUser); err != nil {
		http.Error(w, "failed to add to whitelist: "+err.Error(), 500)
		return
	}

	content := inviteTreeComponent(r.Context(), "", loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

func removeFromWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)
	pubkey := r.PostFormValue("pubkey")
	if err := removeFromWhitelist(pubkey, loggedUser); err != nil {
		http.Error(w, "failed to remove from whitelist: "+err.Error(), 500)
		return
	}
	content := inviteTreeComponent(r.Context(), "", loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

func reportsViewerHandler(w http.ResponseWriter, r *http.Request) {
	events, err := db.QueryEvents(r.Context(), nostr.Filter{
		Kinds: []int{1984},
		Limit: 52,
	})
	if err != nil {
		http.Error(w, "failed to query reports: "+err.Error(), 500)
		return
	}

	content := reportsPageHTML(r.Context(), ReportsPageParams{
		reports:    events,
		loggedUser: getLoggedUser(r),
	})
	htmlgo.Fprint(w, content, r.Context())
}
