package wallet

import (
	"context"
	"io"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
)

const server = "https://edge.qiwi.com"

// Wallet is a client wrapper for qiwi wallet
type Wallet struct {
	token  string
	client ClientWithResponsesInterface
	// personID int
}

// New creates a new qiwi wallet
func New(token string) (w *Wallet, err error) {

	var bearer *securityprovider.SecurityProviderBearerToken
	if bearer, err = securityprovider.NewSecurityProviderBearerToken(token); err == nil {

		var client ClientWithResponsesInterface
		if client, err = NewClientWithResponses(server,
			WithRequestEditorFn(bearer.Intercept)); err == nil {

			w = &Wallet{
				token:  token,
				client: client,
			}
		}
	}

	return
}

// /////////////////////////////////////////////////////////////////////////////

// GetPersonID returns phone number by current token
func (w *Wallet) GetPersonID(ctx context.Context) (personID int, err error) {
	var (
		f                 = false
		getPersonIDParams = &GetPersonProfileParams{
			ContractInfoEnabled: &f,
			UserInfoEnabled:     &f,
		}
	)

	var response *GetPersonProfileResponse
	if response, err = w.client.GetPersonProfileWithResponse(ctx, getPersonIDParams); err == nil {

		personID = *response.JSON200.AuthInfo.PersonId
	}

	return
}

// /////////////////////////////////////////////////////////////////////////////

// GetChequeBytes возвращает файл квитанции о транзакции
func (w *Wallet) GetChequeBytes(ctx context.Context,
	transactionID int, transactionType string, format string) (r io.Reader, err error) {

	var response *GetChequeBytesResponse
	if response, err = w.client.GetChequeBytesWithResponse(ctx, transactionID,
		&GetChequeBytesParams{
			Type:   &transactionType,
			Format: &format,
		}); err == nil {

		r = strings.NewReader((*response.JSON200)[0])
	}

	return
}

// SendCheque отправляет квитанцию о транзакции на email
func (w *Wallet) SendCheque(ctx context.Context,
	transactionID int, transactionType string, email string) (err error) {

	_, err = w.client.SendChequeWithResponse(ctx, transactionID,
		&SendChequeParams{
			Type: &transactionType,
		},
		SendChequeJSONRequestBody{
			Email: email,
		})

	return
}
