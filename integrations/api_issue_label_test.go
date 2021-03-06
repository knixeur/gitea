// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"fmt"
	"net/http"
	"testing"

	"code.gitea.io/gitea/models"
	api "code.gitea.io/sdk/gitea"

	"github.com/stretchr/testify/assert"
)

func TestAPIAddIssueLabels(t *testing.T) {
	assert.NoError(t, models.LoadFixtures())

	repo := models.AssertExistsAndLoadBean(t, &models.Repository{ID: 1}).(*models.Repository)
	issue := models.AssertExistsAndLoadBean(t, &models.Issue{RepoID: repo.ID}).(*models.Issue)
	label := models.AssertExistsAndLoadBean(t, &models.Label{RepoID: repo.ID}).(*models.Label)
	owner := models.AssertExistsAndLoadBean(t, &models.User{ID: repo.OwnerID}).(*models.User)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/labels",
		owner.Name, repo.Name, issue.Index)
	req := NewRequestWithJSON(t, "POST", urlStr, &api.IssueLabelsOption{
		Labels: []int64{label.ID},
	})
	session := loginUser(t, owner.Name)
	resp := session.MakeRequest(t, req)
	assert.EqualValues(t, http.StatusOK, resp.HeaderCode)
	var apiLabels []*api.Label
	DecodeJSON(t, resp, &apiLabels)
	assert.Len(t, apiLabels, models.GetCount(t, &models.IssueLabel{IssueID: issue.ID}))

	models.AssertExistsAndLoadBean(t, &models.IssueLabel{IssueID: issue.ID, LabelID: label.ID})
}

func TestAPIReplaceIssueLabels(t *testing.T) {
	assert.NoError(t, models.LoadFixtures())

	repo := models.AssertExistsAndLoadBean(t, &models.Repository{ID: 1}).(*models.Repository)
	issue := models.AssertExistsAndLoadBean(t, &models.Issue{RepoID: repo.ID}).(*models.Issue)
	label := models.AssertExistsAndLoadBean(t, &models.Label{RepoID: repo.ID}).(*models.Label)
	owner := models.AssertExistsAndLoadBean(t, &models.User{ID: repo.OwnerID}).(*models.User)

	urlStr := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/labels",
		owner.Name, repo.Name, issue.Index)
	req := NewRequestWithJSON(t, "PUT", urlStr, &api.IssueLabelsOption{
		Labels: []int64{label.ID},
	})
	session := loginUser(t, owner.Name)
	resp := session.MakeRequest(t, req)
	assert.EqualValues(t, http.StatusOK, resp.HeaderCode)
	var apiLabels []*api.Label
	DecodeJSON(t, resp, &apiLabels)
	assert.Len(t, apiLabels, 1)
	assert.EqualValues(t, label.ID, apiLabels[0].ID)

	models.AssertCount(t, &models.IssueLabel{IssueID: issue.ID}, 1)
	models.AssertExistsAndLoadBean(t, &models.IssueLabel{IssueID: issue.ID, LabelID: label.ID})
}
