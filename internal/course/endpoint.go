package course

import (
	"context"
	"errors"
	"github.com/og11423074s/go_lib_response/response"
	"github.com/og11423074s/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	UpdateReq struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
		ID        string  `json:"id"`
	}

	GetAllReq struct {
		Name  string `json:"name"`
		Limit int    `json:"limit"`
		Page  int    `json:"page"`
	}

	GetReq struct {
		ID string `json:"id"`
	}

	DeleteReq struct {
		ID string `json:"id"`
	}

	Config struct {
		LimPageDef string
	}
)

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		// validations

		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		course, err := s.Create(ctx, req.Name, req.StartDate, req.EndDate)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", course, nil), nil
	}

}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			name: req.Name,
		}

		// select count(*) from courses where name = ?
		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		metaResult, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)

		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		courses, err := s.GetAll(ctx, filters, metaResult.Offset(), metaResult.Limit())
		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		return response.OK("success", courses, metaResult), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		course, err := s.Get(ctx, req.ID)

		if err != nil {

			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", course, nil), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)

		if req.Name != nil && *req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		if req.StartDate != nil && *req.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateRequired.Error())
		}

		if req.EndDate != nil && *req.EndDate == "" {
			return nil, response.BadRequest(ErrEndDateRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.Name, req.StartDate, req.EndDate); err != nil {

			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		if err := s.Delete(ctx, req.ID); err != nil {

			if errors.As(err, &ErrorNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}
