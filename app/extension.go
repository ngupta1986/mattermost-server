// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/mattermost-server/model"
	"html/template"
	"net/http"
)

func (a *App) ExtensionSupportEnabled() bool {
	return *a.Config().ExtensionSettings.EnableExperimentalExtensions
}

func (a *App) ValidateExtension(extensionID string) bool {
	extensionIsValid := false
	extensionIDs := a.Config().ExtensionSettings.AllowedExtensionsIDs

	for _, id := range extensionIDs {
		if extensionID == id {
			extensionIsValid = true
		}
	}

	if !extensionIsValid {
		return false
	}

	return true
}

func (a *App) SendMessageToExtension(w http.ResponseWriter, extensionId string, token string) *model.AppError {
	var t *template.Template
	var err error
	if len(extensionId) == 0 {
		return model.NewAppError("completeSaml", "api.user.saml.extension_id.app_error", nil, "", http.StatusInternalServerError)
	}

	t = template.New("complete_saml_extension_body")
	t, err = t.ParseFiles("templates/complete_saml_extension_body.html")

	if err != nil {
		return model.NewAppError("completeSaml", "api.user.saml.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	var errMessage string
	if len(token) == 0 {
		loginError := model.NewAppError("completeSaml", "api.user.saml.app_error", nil, "", http.StatusInternalServerError)
		errMessage = loginError.Message
	}

	data := struct {
		ExtensionId string
		Token       string
		Error       string
	}{
		extensionId,
		token,
		errMessage,
	}

	if err := t.Execute(w, data); err != nil {
		return model.NewAppError("completeSaml", "api.user.saml.app_error", nil, "err="+err.Error(), http.StatusInternalServerError)
	}

	return nil
}
