package settings

import "github.com/number571/hidden-lake/internal/service/pkg/request"

type SPubKey struct {
	FKEMPKey string `json:"kem_pkey"`
	FDSAPKey string `json:"dsa_pkey"`
}

type SFriend struct {
	FAliasName string `json:"alias_name"`
	FPublicKey string `json:"public_key"`
}

type SRequest struct {
	FReceiver string            `json:"receiver"` // alias_name
	FReqData  *request.SRequest `json:"req_data"`
}
