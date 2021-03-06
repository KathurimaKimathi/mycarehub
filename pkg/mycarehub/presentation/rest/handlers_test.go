package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

func createSendOTPPayload(phonenumber string, flavour feedlib.Flavour) []byte {
	payload := &dto.SendOTPInput{
		PhoneNumber: phonenumber,
		Flavour:     flavour,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return marshalled
}

func createLoginPayload(phonenumber *string, pin string, flavour feedlib.Flavour) []byte {
	payload := &dto.LoginInput{
		PhoneNumber: phonenumber,
		PIN:         &pin,
		Flavour:     flavour,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return marshalled
}

func TestMyCareHubHandlersInterfacesImpl_SendOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/send_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/send_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Close = true

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			defer resp.Body.Close()

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_SendRetryOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/send_retry_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/send_retry_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Close = true

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			defer resp.Body.Close()

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_RequestPINReset(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/request_pin_reset",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/request_pin_reset",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_LoginByPhone(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	phoneNumber := interserviceclient.TestUserPhoneNumber
	invalidPayload := createLoginPayload(&phoneNumber, "1234", feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createLoginPayload(nil, "1234", feedlib.FlavourConsumer)
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/login_by_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/login_by_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_VerifyPhone(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_phone",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_VerifyOTP(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := createSendOTPPayload(interserviceclient.TestUserPhoneNumber, feedlib.Flavour("invalid flavour"))
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid flavour defined",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/verify_otp",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_GetUserRespondedSecurityQuestions(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}
	invalidPayload1 := createSendOTPPayload("", feedlib.FlavourConsumer)

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Missing phonenumber",
			args: args{
				url: fmt.Sprintf(
					"%s/get_user_responded_security_questions",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidPayload1),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_RefreshToken(t *testing.T) {
	ctx := context.Background()
	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	invalidPayload := &dto.RefreshTokenPayload{
		UserID: nil,
	}
	marshalled, err := json.Marshal(invalidPayload)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Missing user ID",
			args: args{
				url: fmt.Sprintf(
					"%s/refresh_token",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(marshalled),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.DefaultClient
			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			dataResponse, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}
			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}
