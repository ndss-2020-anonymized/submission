package dac

import (
	"encoding/asn1"
	"fmt"

	"github.com/ndss-2020-anonymized/fabric-amcl/amcl"
	"github.com/ndss-2020-anonymized/fabric-amcl/amcl/FP256BN"
)

// NymSignature is signature / NIZK proof of knowledge of pseudonym's
// secret key sk and randomness skNym
type NymSignature struct {
	resSk      *FP256BN.BIG
	resSkNym   *FP256BN.BIG
	commitment interface{}
}

// GenerateNymKeys generates a fresh pair of pseudonym keys
// Nym object is needed to commit to a secret key without revealing it
func GenerateNymKeys(prg *amcl.RAND, sk SK, h *FP256BN.ECP) (skNym SK, pkNym PK) {
	q := FP256BN.NewBIGints(FP256BN.CURVE_Order)
	g := FP256BN.ECP_generator()

	skNym = FP256BN.Randomnum(q, prg)
	pkNym = productOfExponents(g, sk, h, skNym)

	return
}

// SignNym generates a proof of knowledge of pseudonym's secret key sk and randomness skNym
func SignNym(prg *amcl.RAND, pkNym PK, skNym SK, sk SK, h *FP256BN.ECP, m []byte) (signature NymSignature) {
	q := FP256BN.NewBIGints(FP256BN.CURVE_Order)
	g := FP256BN.ECP_generator()

	t1 := FP256BN.Randomnum(q, prg)
	t2 := FP256BN.Randomnum(q, prg)

	signature.commitment = productOfExponents(g, t1, h, t2)

	c := hashNym(q, signature.commitment, pkNym, m)

	signature.resSk = FP256BN.Modmul(sk, c, q).Plus(t1)
	signature.resSk.Mod(q)

	signature.resSkNym = FP256BN.Modmul(skNym, c, q).Plus(t2)
	signature.resSkNym.Mod(q)

	return
}

// VerifyNym verifies the proof of knowledge of pseudonym's secret key sk and randomness skNym
func (signature *NymSignature) VerifyNym(h *FP256BN.ECP, pkNym PK, m []byte) (e error) {
	q := FP256BN.NewBIGints(FP256BN.CURVE_Order)
	g := FP256BN.ECP_generator()

	c := hashNym(q, signature.commitment, pkNym, m)

	LHS := pointMultiply(pkNym, c)
	pointAdd(LHS, signature.commitment)

	RHS := productOfExponents(g, signature.resSk, h, signature.resSkNym)

	if !pointEqual(LHS, RHS) {
		return fmt.Errorf("VerifyNym: verification failed")
	}

	return
}

func hashNym(q *FP256BN.BIG, commitment interface{}, pkNym PK, m []byte) *FP256BN.BIG {
	var raw []byte
	raw = append(raw, pointToBytes(commitment)...)
	raw = append(raw, pointToBytes(pkNym)...)
	raw = append(raw, m...)

	return sha3(q, raw)
}

type nymSignatureMarshal struct {
	ResSk      []byte
	ResSkNym   []byte
	Commitment []byte
}

// ToBytes marshals the NIZK object using ASN1 encoding
func (signature *NymSignature) ToBytes() (result []byte) {
	var marshal nymSignatureMarshal

	marshal.ResSk = bigToBytes(signature.resSk)
	marshal.ResSkNym = bigToBytes(signature.resSkNym)
	marshal.Commitment = pointToBytes(signature.commitment)

	result, _ = asn1.Marshal(marshal)

	return
}

// NymSignatureFromBytes un-marshals the NIZK object using ASN1 encoding
func NymSignatureFromBytes(input []byte) (signature *NymSignature) {
	var marshal nymSignatureMarshal
	asn1.Unmarshal(input, &marshal)

	signature = &NymSignature{}

	signature.commitment, _ = pointFromBytes(marshal.Commitment)
	signature.resSk = FP256BN.FromBytes(marshal.ResSk)
	signature.resSkNym = FP256BN.FromBytes(marshal.ResSkNym)

	return
}

func (signature *NymSignature) equals(other *NymSignature) (result bool) {
	if !pointEqual(signature.commitment, other.commitment) {
		return
	}

	if !bigEqual(signature.resSk, other.resSk) {
		return
	}

	if !bigEqual(signature.resSkNym, other.resSkNym) {
		return
	}

	return true
}
