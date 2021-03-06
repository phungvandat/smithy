package endpoints

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/dwarvesf/smithy/backend/service"
)

// Endpoints .
type Endpoints struct {
	AgentSync       endpoint.Endpoint
	AvailableModels endpoint.Endpoint
	AddHook         endpoint.Endpoint
	DBQuery         endpoint.Endpoint
	DBCreate        endpoint.Endpoint
	DBUpdate        endpoint.Endpoint
	DBDelete        endpoint.Endpoint
	ListVersion     endpoint.Endpoint
	RevertVersion   endpoint.Endpoint
	Login           endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		AgentSync:       makeAgentSyncEndpoint(s),
		DBQuery:         makeDBQueryEndpoint(s),
		DBCreate:        makeDBCreateEndpoint(s),
		DBUpdate:        makeDBUpdateEndpoint(s),
		DBDelete:        makeDBDeleteEndpoint(s),
		AvailableModels: makeAvailableModelsEndpoint(s),
		AddHook:         makeAddHookEndpoint(s),
		ListVersion:     makeListVersionEndpoint(s),
		RevertVersion:   makeRevertVersionEndpoint(s),
		Login:           makeLoginEndpoint(s),
	}
}
