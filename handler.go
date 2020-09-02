package mailsender

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Name    string
	Email   string
	Message string
}

type Response struct {
	Code int
	Msg  string
}

func MailHandler(sender *Smtp, fromMail string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			resp := &Response{
				Code: http.StatusMethodNotAllowed,
				Msg:  wrongMethod,
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(resp)
			return
		}

		req := &Request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			resp := &Response{
				Code: http.StatusBadRequest,
				Msg:  badRequest,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		fmt.Println(req.Name)
		fmt.Println(req.Email)
		fmt.Println(req.Message)

		subj := fmt.Sprintf("Landing mail from %s", req.Name)

		if err := sender.Send(fromMail, req.Email, subj, req.Message); err != nil {
			fmt.Println(err)

			resp := &Response{
				Code: http.StatusInternalServerError,
				Msg:  internalServerError,
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := &Response{
			Code: http.StatusOK,
			Msg:  "Sucess",
		}

		json.NewEncoder(w).Encode(resp)
	}
}
