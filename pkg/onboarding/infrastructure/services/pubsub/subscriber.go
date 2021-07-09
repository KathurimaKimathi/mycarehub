package pubsubmessaging

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/common"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

// ReceivePubSubPushMessages receives and processes a Pub/Sub push message.
func (ps ServicePubSubMessaging) ReceivePubSubPushMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	message, err := ps.baseExt.VerifyPubSubJWTAndDecodePayload(w, r)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	topicID, err := ps.baseExt.GetPubSubTopic(message)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusBadRequest,
		)
		return
	}

	ctx := r.Context()
	switch topicID {
	case ps.AddPubSubNamespace(common.CreateCustomerTopic):
		var data dto.CustomerPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err := ps.erp.CreateERPCustomer(
			ctx,
			data.CustomerPayload,
			data.UID,
		); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	case ps.AddPubSubNamespace(common.CreateSupplierTopic):
		var data dto.SupplierPubSubMessage
		err := json.Unmarshal(message.Message.Data, &data)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err := ps.erp.CreateERPSupplier(
			ctx,
			data.SupplierPayload,
			data.UID,
		); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	case ps.AddPubSubNamespace(common.CreateCRMContact):
		var CRMContact domain.CRMContact
		err := json.Unmarshal(message.Message.Data, &CRMContact)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err = ps.crm.CreateContact(CRMContact); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	case ps.AddPubSubNamespace(common.UpdateCRMContact):
		var contactProperties dto.UpdateContactPSMessage
		err := json.Unmarshal(message.Message.Data, &contactProperties)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}
		if _, err = ps.crm.UpdateContact(
			contactProperties.Phone,
			contactProperties.Properties,
		); err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	case ps.AddPubSubNamespace(common.LinkCoverTopic):
		var userDetails dto.LinkCoverPubSubMessage
		err := json.Unmarshal(message.Message.Data, &userDetails)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

		_, err = ps.edi.LinkCover(ctx, userDetails.PhoneNumber, userDetails.UID)
		if err != nil {
			ps.baseExt.WriteJSONResponse(
				w,
				ps.baseExt.ErrorMap(err),
				http.StatusBadRequest,
			)
			return
		}

	default:
		errMsg := fmt.Sprintf(
			"pub sub handler error: unknown topic `%s`",
			topicID,
		)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	resp := map[string]string{"status": "success"}
	marshalledSuccessMsg, err := json.Marshal(resp)
	if err != nil {
		ps.baseExt.WriteJSONResponse(
			w,
			ps.baseExt.ErrorMap(err),
			http.StatusInternalServerError,
		)
		return
	}
	_, _ = w.Write(marshalledSuccessMsg)
}