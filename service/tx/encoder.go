package tx

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	sdkTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/encoding/protowire"
)

// Decoder defines the interface for decoding transaction bytes into a DecodedTx.
//type Decoder interface {
//	Decode(txBytes []byte) (*txdecode.DecodedTx, error)
//}

type MsgEncoder interface {
	Decode(txBytes []byte) (sdk.Tx, error)
	Encode(tx sdk.Tx) ([]byte, error)
	EncodeMsgs(msgs ...sdk.Msg) ([]byte, error)
}

type DefaultCoder struct {
	cdc     codec.Codec
	decoder func(txBytes []byte) (sdk.Tx, error)
}

func NewDefaultDecoder(cdc codec.Codec) *DefaultCoder {
	return &DefaultCoder{
		cdc:     cdc,
		decoder: sdkTx.DefaultTxDecoder(cdc),
	}
}

func (d *DefaultCoder) Decode(txBytes []byte) (sdk.Tx, error) {
	// Make sure txBytes follow ADR-027.
	err := rejectNonADR027TxRaw(txBytes)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	var raw txtypes.TxRaw

	// reject all unknown proto fields in the root TxRaw
	err = unknownproto.RejectUnknownFieldsStrict(txBytes, &raw, d.cdc.InterfaceRegistry())
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	err = d.cdc.Unmarshal(txBytes, &raw)
	if err != nil {
		return nil, err
	}

	var body txtypes.TxBody

	// allow non-critical unknown fields in TxBody
	txBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true, d.cdc.InterfaceRegistry())
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	err = d.cdc.Unmarshal(raw.BodyBytes, &body)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	var authInfo txtypes.AuthInfo

	// reject all unknown proto fields in AuthInfo
	err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo, d.cdc.InterfaceRegistry())
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	err = d.cdc.Unmarshal(raw.AuthInfoBytes, &authInfo)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	theTx := &txtypes.Tx{
		Body:       &body,
		AuthInfo:   &authInfo,
		Signatures: raw.Signatures,
	}

	return &Wrapper{
		tx:                           theTx,
		bodyBz:                       raw.BodyBytes,
		authInfoBz:                   raw.AuthInfoBytes,
		txBodyHasUnknownNonCriticals: txBodyHasUnknownNonCriticals,
		cdc:                          d.cdc,
	}, nil
}

func (d *DefaultCoder) Encode(tx sdk.Tx) ([]byte, error) {
	txWrapper, ok := tx.(*Wrapper)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &Wrapper{}, tx)
	}

	raw := &txtypes.TxRaw{
		BodyBytes:     txWrapper.getBodyBytes(),
		AuthInfoBytes: txWrapper.getAuthInfoBytes(),
		Signatures:    txWrapper.tx.Signatures,
	}

	return proto.Marshal(raw)
}

func (d *DefaultCoder) EncodeMsgs(msgs ...sdk.Msg) ([]byte, error) {
	builder := NewBuilder(d.cdc)
	err := builder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}

	return d.Encode(builder.GetTx())
}

// rejectNonADR027TxRaw rejects txBytes that do not follow ADR-027. This is NOT
// a generic ADR-027 checker, it only applies decoding TxRaw. Specifically, it
// only checks that:
// - field numbers are in ascending order (1, 2, and potentially multiple 3s),
// - and varints are as short as possible.
// All other ADR-027 edge cases (e.g. default values) are not applicable with
// TxRaw.
func rejectNonADR027TxRaw(txBytes []byte) error {
	// Make sure all fields are ordered in ascending order with this variable.
	prevTagNum := protowire.Number(0)

	for len(txBytes) > 0 {
		tagNum, wireType, m := protowire.ConsumeTag(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// TxRaw only has bytes fields.
		if wireType != protowire.BytesType {
			return fmt.Errorf("expected %d wire type, got %d", protowire.BytesType, wireType)
		}
		// Make sure fields are ordered in ascending order.
		if tagNum < prevTagNum {
			return fmt.Errorf("txRaw must follow ADR-027, got tagNum %d after tagNum %d", tagNum, prevTagNum)
		}
		prevTagNum = tagNum

		// All 3 fields of TxRaw have wireType == 2, so their next component
		// is a varint, so we can safely call ConsumeVarint here.
		// Byte structure: <varint of bytes length><bytes sequence>
		// Inner  fields are verified in `DefaultTxDecoder`
		lengthPrefix, m := protowire.ConsumeVarint(txBytes[m:])
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		// We make sure that this varint is as short as possible.
		n := varintMinLength(lengthPrefix)
		if n != m {
			return fmt.Errorf("length prefix varint for tagNum %d is not as short as possible, read %d, only need %d", tagNum, m, n)
		}

		// Skip over the bytes that store fieldNumber and wireType bytes.
		_, _, m = protowire.ConsumeField(txBytes)
		if m < 0 {
			return fmt.Errorf("invalid length; %w", protowire.ParseError(m))
		}
		txBytes = txBytes[m:]
	}

	return nil
}

// varintMinLength returns the minimum number of bytes necessary to encode an
// uint using varint encoding.
func varintMinLength(n uint64) int {
	switch {
	// Note: 1<<N == 2**N.
	case n < 1<<(7):
		return 1
	case n < 1<<(7*2):
		return 2
	case n < 1<<(7*3):
		return 3
	case n < 1<<(7*4):
		return 4
	case n < 1<<(7*5):
		return 5
	case n < 1<<(7*6):
		return 6
	case n < 1<<(7*7):
		return 7
	case n < 1<<(7*8):
		return 8
	case n < 1<<(7*9):
		return 9
	default:
		return 10
	}
}
