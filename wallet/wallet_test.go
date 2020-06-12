package wallet

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type WalletTestSuite struct {
	suite.Suite
}

func TestWalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) TestGetPersonID() {
	const (
		token        = "test-token"
		personID int = 79101234567
	)
	defer gock.Off()

	gock.New("https://edge.qiwi.com").
		Get("/person-profile/v1/profile/current").
		MatchHeaders(map[string]string{
			"Authorization": fmt.Sprintf(`^Bearer %s$`, token),
			// TODO header not present
			// "Accept": `^application/json`,
		}).
		// MatchParam("authInfoEnabled", "true").
		MatchParam("contractInfoEnabled", "false").
		MatchParam("userInfoEnabled", "false").
		Reply(200).
		AddHeader("Content-type", "application/json").
		BodyString(fmt.Sprintf(`{"authInfo":{"personId":%d}}`, personID))

	w, err := New("test-token")
	s.Nil(err)

	p, err := w.GetPersonID(context.Background())

	s.Nil(err)
	s.Equal(personID, p)

	s.True(gock.IsDone())
}

// /////////////////////////////////////////////////////////////////////////////

func (s *WalletTestSuite) TestGetChequeBytes() {
	const (
		token           = "test-token"
		transactionID   = 666
		transactionType = "IN"
		chequeFormat    = "PDF"
		chequeFile      = "SOME BINARY STUFF"
	)
	defer gock.Off()

	gock.New("https://edge.qiwi.com").
		Get(fmt.Sprintf("/payment-history/v1/transactions/%d/cheque/file", transactionID)).
		MatchHeaders(map[string]string{
			"Authorization": fmt.Sprintf(`^Bearer %s$`, token),
			// "Accept":        `^application/json`,
		}).
		MatchParam("type", transactionType).
		Reply(200).
		AddHeader("Content-type", "application/json").
		BodyString(fmt.Sprintf(`["%s"]`, chequeFile))

	w, err := New(token)
	s.Nil(err)

	res, err := w.GetChequeBytes(context.Background(), transactionID, transactionType, chequeFormat)
	s.Nil(err)

	buf, err := ioutil.ReadAll(res)
	s.Nil(err)
	s.Equal(chequeFile, string(buf))

	s.True(gock.IsDone())
}

func (s *WalletTestSuite) TestSendCheque() {
	const (
		token           = "test-token"
		transactionID   = 666
		transactionType = "IN"
		notifyEmail     = "postmaster@example.com"
	)
	defer gock.Off()

	gock.New("https://edge.qiwi.com").
		Post(fmt.Sprintf("/payment-history/v1/transactions/%d/cheque/send", transactionID)).
		MatchHeaders(map[string]string{
			"Authorization": fmt.Sprintf(`^Bearer %s$`, token),
			// "Accept":        `^application/json`,
			"Content-type": `^application/json`,
		}).
		MatchParam("type", transactionType).
		BodyString(fmt.Sprintf(`{"email":"%s"}`, notifyEmail)).
		Reply(201)

	w, err := New(token)
	s.Nil(err)

	err = w.SendCheque(context.Background(), transactionID, transactionType, notifyEmail)
	s.Nil(err)

	s.True(gock.IsDone())
}
