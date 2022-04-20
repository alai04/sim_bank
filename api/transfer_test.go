package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/alai04/sim_bank/db/mock"
	db "github.com/alai04/sim_bank/db/sqlc"
	"github.com/alai04/sim_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	fromAccount, toAccount, transfer := randomTransfer()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          transfer.Amount,
				"currency":        fromAccount.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).
					Times(1).
					Return(fromAccount, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).
					Times(1).
					Return(toAccount, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{
						FromAccountID: fromAccount.ID,
						ToAccountID:   toAccount.ID,
						Amount:        transfer.Amount,
					})).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomTransfer() (fromAccount db.Account, toAccount db.Account, transfter db.Transfer) {
	fromAccount = randomAccount()
	toAccount = randomAccount()
	toAccount.Currency = fromAccount.Currency
	transfter = db.Transfer{
		ID:            util.RandomInt(1, 1000),
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomInt(1, 10),
	}
	return
}
