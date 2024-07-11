package main

import (
	"context"
	"fmt"

	"github.com/fiatjaf/khatru"
	"github.com/nbd-wtf/go-nostr/nip86"
)

func allowPubKeyHandler(ctx context.Context, pubkey, reason string) error {
	loggedUser := khatru.GetAuthed(ctx)

	if loggedUser != s.RelayPubkey && hasInvitedAtLeast(loggedUser, s.MaxInvitesPerPerson) {
		return fmt.Errorf("cannot invite more than %d", s.MaxInvitesPerPerson)
	}
	if err := addToWhitelist(pubkey, loggedUser); err != nil {
		return fmt.Errorf("failed to add to whitelist: %w", err)
	}

	return nil
}

func banPubKeyHandler(ctx context.Context, pubkey, reason string) error {
	loggedUser := khatru.GetAuthed(ctx)

	// check if this user is a descendant of the user who issued the delete command
	if !isAncestorOf(loggedUser, pubkey) {
		return fmt.Errorf("insufficient permissions to delete this")
	}

	// if we got here that means we have permission to delete the target
	delete(whitelist, pubkey)

	// delete all people who were invited by the target
	removeDescendantsFromWhitelist(pubkey)

	return saveWhitelist()
}

func listAllowedPubKeysHandler(ctx context.Context) ([]nip86.PubKeyReason, error) {
	list := make([]nip86.PubKeyReason, len(whitelist))
	i := 0
	for pubkey, inviter := range whitelist {
		list[i] = nip86.PubKeyReason{PubKey: pubkey, Reason: fmt.Sprintf("invited by %s", inviter)}
		i++
	}
	return list, nil
}
