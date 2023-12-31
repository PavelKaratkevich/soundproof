package eth

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/storyicon/sigverify"
)

func ParseMetamaskSignedString(signedMessage, signature string) (string, error) {
	// convert message to bytes
	msg := common.FromHex(signedMessage)

	// decode signature
	sig, err := sigverify.HexDecode(signature)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	// recover address
	addr, err := sigverify.EcRecoverEx(msg, sig)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}

	return addr.String(), nil
}

/////////////////////////////////////////////////////////////////////////
// ifValid helps us validate if the message was ideed signed by the address //
// this function is unused in the present code, but can be reused as a pkg in other projects
// where one would need to ensure propoer validation
/////////////////////////////////////////////////////////////////////////

func ifValid(addr, msg, signature string) (bool, error) {
	return sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(addr),
		common.FromHex(msg),
		signature,
	)
}
