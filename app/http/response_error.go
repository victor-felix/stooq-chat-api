package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-felix/chat-service/app/models"
)

type Title string

const (
	TitleInvalidArgument  Title = "invalid argument"
	TitleNotFound         Title = "resource not found"
	TitleRuleNotSatisfied Title = "internal rule was not satisfied"
	TitleDuplicated       Title = "duplicated record"
	TypeBlank             string       = "about:blank"
	TitleBadRequest       Title = "bad request"
	TitleForbidden        Title = "forbidden"
	TitleUnauthorized     Title = "unauthorized"
)

var (
	ErrInternalServer = NewFailure(TypeBlank, "the server encountered an unexpected condition that prevented it from fulfilling the request", "", "", http.StatusInternalServerError, "", nil)
)

// NewFailure creates a new instance of Failure
func NewFailure(tp, tittle, detail, instance string, status int, code string, args map[string]interface{}) *Failure {
	err := &Failure{}
	err.Type = tp
	err.Title = tittle
	err.Status = status
	err.Instance = instance
	err.Detail = detail
	err.Args = args
	err.Code = code
	return err
}

type Failure struct {
	Type string `json:"type"`
	Status int `json:"status,omitempty"`
	Code string `json:"code,omitempty"`
	Title string `json:"title"`
	Detail string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Args map[string]interface{} `json:"arguments,omitempty"`
}

// Error is the type of error
func (f *Failure) Error() string {
	return f.Detail
}

func (h *handler) responseProblemBuilder(c *gin.Context, err error, status int, title Title) {
	describer, ok := models.ErrorDescriberCast(err)
	var failure *Failure
	if ok {
		failure = &Failure{
			Type:   TypeBlank,
			Status: status,
			Code:   string(describer.GetCode()),
			Title:  string(title),
			Detail: describer.GetMessage(),
			Args:   describer.GetArgs(),
		}
	} else {
		failure = ErrInternalServer
	}

	_ = c.Error(err)
	h.responseProblemWriter(c, failure)
}

func (h *handler) responseProblemWriter(c *gin.Context, rfc *Failure) {
	if problemJSON, err := json.Marshal(rfc); err == nil {
		c.Data(rfc.Status, "application/problem+json; charset=utf-8", problemJSON)
		return
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}

func (h *handler) responseProblem(c *gin.Context, err error) {
	if nf, ok := models.ErrNotFoundCast(err); ok {
		h.responseProblemBuilder(c, nf, http.StatusNotFound, TitleNotFound)
		return
	}

	if ia, ok := models.ErrInvalidArgumentCast(err); ok {
		h.responseProblemBuilder(c, ia, http.StatusBadRequest, TitleInvalidArgument)
		return
	}

	if rns, ok := models.ErrRulesNotSatisfiedCast(err); ok {
		h.responseProblemBuilder(c, rns, http.StatusBadRequest, TitleRuleNotSatisfied)
		return
	}

	if br, ok := models.ErrBadRequestCast(err); ok {
		h.responseProblemBuilder(c, br, http.StatusBadRequest, TitleBadRequest)
		return
	}

	if fc, ok := models.ErrForbiddenCast(err); ok {
		h.responseProblemBuilder(c, fc, http.StatusForbidden, TitleForbidden)
		return
	}

	if un, ok := models.ErrUnauthorizedCast(err); ok {
		h.responseProblemBuilder(c, un, http.StatusUnauthorized, TitleUnauthorized)
		return
	}

	h.responseProblemBuilder(c, err, http.StatusInternalServerError, "")
}
