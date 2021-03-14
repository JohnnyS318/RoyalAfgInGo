package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/JohnnyS318/RoyalAfgInGo/pkg/dtos"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/errors"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/mw"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/poker/pokerConfig"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/poker/ticketToken"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/responses"
	"github.com/JohnnyS318/RoyalAfgInGo/services/poker-matchmaker/serviceconfig"
)

//TicketResponse is the successful response of a ticket request
type TicketResponse struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

//GetTicketWithParams requests a ticket with lobby params
func (h *Ticket) GetTicketWithParams(rw http.ResponseWriter, r *http.Request) {
	if err := mw.ValidateCSRF(r); err != nil {
		h.logger.Errorw("could not validate csrf token", "error", err)
		responses.JSONError(rw, &responses.ErrorResponse{Error: "wrong format decoding failed"}, http.StatusForbidden)
		return
	}

	vals := r.URL.Query()
	claims := mw.FromUserTokenContext(r.Context().Value("user"))

	class, err := strconv.Atoi(vals.Get("class"))
	if err != nil {
		h.logger.Errorw("Invalid Class", "error", err)
		http.Error(rw, "Either a valid class or a lobby class has to be given", http.StatusBadRequest)
		return
	}

	buyIn, err := strconv.Atoi(vals.Get("buyIn"))
	if err != nil {
		h.logger.Errorw("Invalid BuyIn", "error", err)
		http.Error(rw, "the buyIn has to be valid", http.StatusBadRequest)
		return
	}

	if viper.GetBool("include_bank_service_validation") {
		if err = VerifyBuyIn(claims.ID, buyIn); err != nil {
			h.logger.Errorw("Error during bank service validation", "error", err)
			http.Error(rw, "the buyIn has to be lower that the users wallet", http.StatusUnprocessableEntity)
			return
		}
	}

	res, err := h.manager.RequestTicket(class)

	if err != nil {
		h.logger.Errorw("Error during ticket request", "error", err)
		http.Error(rw, "something went wrong during a lobby search", http.StatusInternalServerError)
		return
	}

	h.logger.Infow("Generate Ticket", "username", claims.Username, "id", claims.ID, "lobbyId", res.LobbyId, "buyIn", buyIn)
	token, err := ticketToken.GenerateTicketToken(claims.Username, claims.ID, res.LobbyId, buyIn, viper.GetString(pokerConfig.MatchMakerJWTKey))

	json.NewEncoder(rw).Encode(&TicketResponse{Address: res.Address, Token: token})
}

//GetTicketWithID requests a ticket with lobby id
func (h *Ticket) GetTicketWithID(rw http.ResponseWriter, r *http.Request) {

	if err := mw.ValidateCSRF(r); err != nil {
		h.logger.Errorw("could not validate csrf token", "error", err)
		responses.JSONError(rw, &responses.ErrorResponse{Error: "wrong format decoding failed"}, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	claims := mw.FromUserTokenContext(r.Context().Value("user"))

	buyIn, err := strconv.Atoi(r.URL.Query().Get("buyIn"))
	if err != nil {
		h.logger.Errorw("Invalid BuyIn", "error", err)
		http.Error(rw, "the buyIn has to be valid", http.StatusBadRequest)
		return
	}

	if viper.GetBool("include_bank_service_validation") {
		if err = VerifyBuyIn(claims.ID, buyIn); err != nil {
			h.logger.Errorw("Error during bank service validation", "error", err)
			http.Error(rw, "the buyIn has to be lower that the users wallet", http.StatusUnprocessableEntity)
			return
		}
	}

	id, ok := vars["id"]
	if !ok || id == "" {
		http.Error(rw, "Either a valid class or a lobby class has to be given", http.StatusBadRequest)
		return
	}

	res, err := h.manager.Connect(id)

	if err != nil {
		h.logger.Errorw("error during connection", "error", err)
		http.Error(rw, "a lobby with the given id is not found", http.StatusNotFound)
		return
	}

	h.logger.Infof("Creating token for [%v;%v] to join lobby[%v] with %v", claims.Username, claims.ID, res.LobbyId, buyIn)
	token, err := ticketToken.GenerateTicketToken(claims.Username, claims.ID, res.LobbyId, buyIn, viper.GetString(pokerConfig.MatchMakerJWTKey))

	json.NewEncoder(rw).Encode(&TicketResponse{Address: res.Address, Token: token})
}

//VerifyBuyIn verifies the buy in amount against the user wallet using the bank service
func VerifyBuyIn(userId string, buyIn int) error {
	client := &http.Client{
		Timeout: 25 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/bank/verifyAmount", viper.GetString(serviceconfig.BankServiceUrl)), nil)

	q := req.URL.Query()
	q.Add("userId", userId)
	q.Add("amount", strconv.Itoa(buyIn))
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		text, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Bank service responded with a non 200 response. %v", text)
	}

	var result dtos.VerifyAmount
	err = json.NewDecoder(res.Body).Decode(&result)
	res.Body.Close()
	if err != nil {
		return err
	}

	if !result.VerificationResult {
		return &errors.InvalidBuyIn{}
	}

	return nil
}
