package wallet

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"mini-wallet/domain/constant"
	historymodel "mini-wallet/domain/model/history"
	model "mini-wallet/domain/model/wallet"
	"mini-wallet/domain/util"
	"net/http"
	"strconv"
)

type handler struct {
	service Service
}

func NewHandler(s Service) HandlerWallet {
	return &handler{service: s}
}

type HandlerWallet interface {
	InitAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	EnableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DisableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	DepositWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	WithdrawWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	GetTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

func (h *handler) InitAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	customerXId := r.FormValue("customer_xid")

	if customerXId == "" {
		logrus.Error(constant.ErrorXid)
		util.FailedResponseWriter(w, constant.ErrorXid, http.StatusBadRequest)
		return
	}

	generatedToken, _ := util.GenerateToken(customerXId)

	token := map[string]string{
		"token": generatedToken,
	}

	util.SuccessResponseWriter(w, token, http.StatusOK)
}

func (h *handler) EnableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ownerId := ctx.Value("ownerId").(string)

	args := model.OwnerWalletParam{
		Owner: ownerId,
	}

	result, err := h.service.EnableWallet(ctx, args)
	if err != nil {
		logrus.Errorf("Handler error : %v", err.Error())
		util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SuccessResponseWriter(w, result, http.StatusOK)

}

func (h *handler) DisableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ownerId := ctx.Value("ownerId").(string)
	isDisabled := r.FormValue("is_disabled")

	if isDisabled == "true" {
		args := model.OwnerWalletParam{
			Owner: ownerId,
		}

		result, err := h.service.DisableWallet(ctx, args)
		if err != nil {
			logrus.Errorf("Handler error : %v", err.Error())
			util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
			return
		}

		util.SuccessResponseWriter(w, result, http.StatusOK)
		return
	}
	util.FailedResponseWriter(w, constant.FailDisable, http.StatusBadRequest)

}

func (h *handler) GetWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ownerId := ctx.Value("ownerId").(string)

	args := model.OwnerWalletParam{
		Owner: ownerId,
	}

	result, err := h.service.GetWallet(ctx, args)
	if err != nil {
		logrus.Errorf("Handler error : %v", err.Error())
		util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SuccessResponseWriter(w, result, http.StatusOK)
}

func (h *handler) DepositWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	amount := r.FormValue("amount")
	amountNumber, _ := strconv.Atoi(amount)

	referenceID := r.FormValue("reference_id")

	args := historymodel.TransactionParams{
		Amount:      int64(amountNumber),
		ReferenceID: referenceID,
	}

	result, err := h.service.DepositWallet(ctx, args)
	if err != nil {
		logrus.Errorf("Handler error : %v", err.Error())
		util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := util.TransformTransactionResult(constant.Deposit, result)
	util.SuccessResponseWriter(w, data, http.StatusOK)

}

func (h *handler) WithdrawWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	amount := r.FormValue("amount")
	amountNumber, _ := strconv.Atoi(amount)

	referenceID := r.FormValue("reference_id")

	args := historymodel.TransactionParams{
		Amount:      int64(amountNumber),
		ReferenceID: referenceID,
	}

	result, err := h.service.WithdrawalWallet(ctx, args)
	if err != nil {
		logrus.Errorf("Handler error : %v", err.Error())
		util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := util.TransformTransactionResult(constant.Withdrawal, result)
	util.SuccessResponseWriter(w, data, http.StatusOK)

}

func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	ownerId := ctx.Value("ownerId").(string)

	args := historymodel.TransactionParams{
		TransactionBy: ownerId,
	}

	result, err := h.service.GetTransaction(ctx, args)
	if err != nil {
		logrus.Errorf("Handler error : %v", err.Error())
		util.FailedResponseWriter(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := util.TransformTransactionResult(constant.Transaction, result)
	util.SuccessResponseWriter(w, data, http.StatusOK)
}
