package send

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/udomday/api-service-backend/iternal/lib/api/response"
	"github.com/udomday/api-service-backend/iternal/lib/logger/sl"
)

type Request struct {
	Url     string            `json:"url" validate:"required,url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

type Response struct {
	resp.Response
	Body string `json:"body"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.request.send.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode reuest body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		payload := []byte(req.Body)

		rq, err := http.NewRequest(req.Method, req.Url, bytes.NewBuffer(payload))
		if err != nil {
			log.Error("invalid create request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid create request"))
			return
		}

		for title, value := range req.Headers {
			rq.Header.Add(title, value)
		}

		client := &http.Client{}

		rs, err := client.Do(rq)
		if err != nil {
			log.Error("invalid create request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid create request"))
			return
		}
		defer rs.Body.Close()

		body, err := io.ReadAll(rs.Body)
		if err != nil {
			log.Error("Do not read answer", sl.Err(err))

			render.JSON(w, r, resp.Error("Do not read answer"))
			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Body:     string(body),
		})
	}
}
