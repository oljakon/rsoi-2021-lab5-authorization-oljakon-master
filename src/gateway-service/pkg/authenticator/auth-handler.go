package authenticator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
)

func LoginHandler(store *sessions.FilesystemStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state := base64.StdEncoding.EncodeToString(b)

		session, err := store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["state"] = state
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authenticator, err := New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusFound)
	})
}

func CallbackHandler(store *sessions.FilesystemStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Get("state") != session.Values["state"] {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		authenticator, err := New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
		if err != nil {
			log.Printf("no token found: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 toke.", http.StatusUnauthorized)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: os.Getenv("AUTH0_CLIENT_ID"),
		}

		idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		session.Values["id_token"] = rawIDToken
		session.Values["profile"] = profile
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}
