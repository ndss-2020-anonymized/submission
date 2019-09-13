package dac

import (
	"encoding/asn1"
	"sort"
	"strconv"
	"sync"

	"github.com/ndss-2020-anonymized/fabric-amcl/amcl/FP256BN"
)

var _ParallelOptimization = true
var _OptimizeTate = true

type proofMarshal struct {
	C      []byte
	RPrime [][]byte
	ResS   [][]byte
	ResT   [][][]byte
	ResA   [][][]byte
	ResCpk [][]byte
	ResCsk []byte
	ResNym []byte
}

// ProofFromBytes un-marshals the proof
func ProofFromBytes(input []byte) (proof *Proof) {
	var marshal proofMarshal
	asn1.Unmarshal(input, &marshal)

	proof = &Proof{}

	proof.c = FP256BN.FromBytes(marshal.C)
	proof.resCsk = FP256BN.FromBytes(marshal.ResCsk)
	proof.resNym = FP256BN.FromBytes(marshal.ResNym)

	proof.rPrime = make([]interface{}, len(marshal.RPrime))
	for i := 0; i < len(marshal.RPrime); i++ {
		proof.rPrime[i], _ = pointFromBytes(marshal.RPrime[i])
	}

	proof.resS = make([]interface{}, len(marshal.ResS))
	for i := 0; i < len(marshal.ResS); i++ {
		proof.resS[i], _ = pointFromBytes(marshal.ResS[i])
	}

	proof.resCpk = make([]interface{}, len(marshal.ResCpk))
	for i := 0; i < len(marshal.ResCpk); i++ {
		proof.resCpk[i], _ = pointFromBytes(marshal.ResCpk[i])
	}

	proof.resT = make([][]interface{}, len(marshal.ResT))
	for i := 0; i < len(marshal.ResT); i++ {
		proof.resT[i] = make([]interface{}, len(marshal.ResT[i]))
		for j := 0; j < len(marshal.ResT[i]); j++ {
			proof.resT[i][j], _ = pointFromBytes(marshal.ResT[i][j])
		}
	}

	proof.resA = make([][]interface{}, len(marshal.ResA))
	for i := 0; i < len(marshal.ResA); i++ {
		proof.resA[i] = make([]interface{}, len(marshal.ResA[i]))
		for j := 0; j < len(marshal.ResA[i]); j++ {
			proof.resA[i][j], _ = pointFromBytes(marshal.ResA[i][j])
		}
	}

	return
}

// ToBytes marshlas the proof
func (proof *Proof) ToBytes() (result []byte) {
	var marshal proofMarshal

	marshal.C = bigToBytes(proof.c)
	marshal.ResCsk = bigToBytes(proof.resCsk)
	marshal.ResNym = bigToBytes(proof.resNym)

	marshal.RPrime = make([][]byte, len(proof.rPrime))
	for i := 0; i < len(proof.rPrime); i++ {
		marshal.RPrime[i] = pointToBytes(proof.rPrime[i])
	}

	marshal.ResS = make([][]byte, len(proof.resS))
	for i := 0; i < len(proof.resS); i++ {
		marshal.ResS[i] = pointToBytes(proof.resS[i])
	}

	marshal.ResCpk = make([][]byte, len(proof.resCpk))
	for i := 0; i < len(proof.resT); i++ {
		marshal.ResCpk[i] = pointToBytes(proof.resCpk[i])
	}

	marshal.ResT = make([][][]byte, len(proof.resT))
	for i := 0; i < len(proof.resT); i++ {
		marshal.ResT[i] = make([][]byte, len(proof.resT[i]))
		for j := 0; j < len(proof.resT[i]); j++ {
			marshal.ResT[i][j] = pointToBytes(proof.resT[i][j])
		}
	}

	marshal.ResA = make([][][]byte, len(proof.resA))
	for i := 0; i < len(proof.resA); i++ {
		marshal.ResA[i] = make([][]byte, len(proof.resA[i]))
		for j := 0; j < len(proof.resA[i]); j++ {
			marshal.ResA[i][j] = pointToBytes(proof.resA[i][j])
		}
	}

	result, _ = asn1.Marshal(marshal)

	return
}

// Equals checks the equality of two proofs
func (proof *Proof) Equals(other Proof) (result bool) {

	if !bigEqual(proof.c, other.c) {
		return
	}

	if !bigEqual(proof.resCsk, other.resCsk) {
		return
	}

	if !bigEqual(proof.resNym, other.resNym) {
		return
	}

	if !pointListEquals(proof.rPrime, other.rPrime) {
		return
	}

	if !pointListEquals(proof.resS, other.resS) {
		return
	}

	if !pointListEquals(proof.resCpk, other.resCpk) {
		return
	}

	if !pointListOfListEquals(proof.resT, other.resT) {
		return
	}

	if !pointListOfListEquals(proof.resA, other.resA) {
		return
	}

	return true
}

type grothSignatureMarshal struct {
	R  []byte
	S  []byte
	Ts [][]byte
}

type credentialsMarshal struct {
	Signatures []grothSignatureMarshal
	Attributes [][][]byte
	PublicKeys [][]byte
}

// CredentialsFromBytes un-marshals the credentials object using ASN1 encoding
func CredentialsFromBytes(input []byte) (creds *Credentials) {
	var marshal credentialsMarshal
	asn1.Unmarshal(input, &marshal)

	creds = &Credentials{}

	creds.signatures = make([]GrothSignature, len(marshal.Signatures))
	for i := 0; i < len(marshal.Signatures); i++ {
		creds.signatures[i].r, _ = pointFromBytes(marshal.Signatures[i].R)
		creds.signatures[i].s, _ = pointFromBytes(marshal.Signatures[i].S)
		creds.signatures[i].ts = make([]interface{}, len(marshal.Signatures[i].Ts))
		for j := 0; j < len(marshal.Signatures[i].Ts); j++ {
			creds.signatures[i].ts[j], _ = pointFromBytes(marshal.Signatures[i].Ts[j])
		}
	}

	creds.publicKeys = make([]interface{}, len(marshal.PublicKeys))
	for i := 0; i < len(marshal.PublicKeys); i++ {
		creds.publicKeys[i], _ = pointFromBytes(marshal.PublicKeys[i])
	}

	creds.attributes = make([][]interface{}, len(marshal.Attributes))
	for i := 0; i < len(marshal.Attributes); i++ {
		creds.attributes[i] = make([]interface{}, len(marshal.Attributes[i]))
		for j := 0; j < len(marshal.Attributes[i]); j++ {
			creds.attributes[i][j], _ = pointFromBytes(marshal.Attributes[i][j])
		}
	}

	return
}

// ToBytes marshals the credentials object using ASN1 encoding
func (creds *Credentials) ToBytes() (result []byte) {
	var marshal credentialsMarshal

	marshal.Signatures = make([]grothSignatureMarshal, len(creds.signatures))
	for i := 0; i < len(marshal.Signatures); i++ {
		marshal.Signatures[i].R = pointToBytes(creds.signatures[i].r)
		marshal.Signatures[i].S = pointToBytes(creds.signatures[i].s)
		marshal.Signatures[i].Ts = make([][]byte, len(creds.signatures[i].ts))
		for j := 0; j < len(creds.signatures[i].ts); j++ {
			marshal.Signatures[i].Ts[j] = pointToBytes(creds.signatures[i].ts[j])
		}
	}

	marshal.PublicKeys = make([][]byte, len(creds.publicKeys))
	for i := 0; i < len(creds.publicKeys); i++ {
		marshal.PublicKeys[i] = pointToBytes(creds.publicKeys[i])
	}

	marshal.Attributes = make([][][]byte, len(creds.attributes))
	for i := 0; i < len(creds.attributes); i++ {
		marshal.Attributes[i] = make([][]byte, len(creds.attributes[i]))
		for j := 0; j < len(creds.attributes[i]); j++ {
			marshal.Attributes[i][j] = pointToBytes(creds.attributes[i][j])
		}
	}

	result, _ = asn1.Marshal(marshal)

	return
}

// Equals checks the equality of two credentials objects
func (creds *Credentials) Equals(other *Credentials) (result bool) {

	if !pointListEquals(creds.publicKeys, other.publicKeys) {
		return
	}

	if !pointListOfListEquals(creds.attributes, other.attributes) {
		return
	}

	if len(creds.signatures) != len(other.signatures) {
		return
	}
	for index := 0; index < len(creds.signatures); index++ {
		if !creds.signatures[index].equals(other.signatures[index]) {
			return
		}
	}

	return true
}

// Index holds the attribute with its position in credentials
type Index struct {
	i, j      int
	attribute interface{}
}

// Indices is an abstraction over the set of Index objects
type Indices []Index

func (indices Indices) Len() int {
	return len(indices)
}
func (indices Indices) Swap(i, j int) {
	indices[i], indices[j] = indices[j], indices[i]
}
func (indices Indices) Less(i, j int) bool {
	return indices[i].i < indices[j].i || indices[i].j < indices[j].j
}

func (indices Indices) contains(i, j int) (attribute interface{}) {
	for _, ij := range indices {
		if ij.i == i && ij.j == j {
			return ij.attribute
		}
	}
	return
}

func (indices Indices) hash() (result []byte) {
	d := make(Indices, len(indices))
	copy(d, indices)
	sort.Sort(d)

	for i := 0; i < len(d); i++ {
		result = append(result, []byte(strconv.Itoa(d[i].i))...)
		result = append(result, []byte(strconv.Itoa(d[i].j))...)
		result = append(result, pointToBytes(d[i].attribute)...)
	}

	return result
}

type eResult struct {
	result *FP256BN.FP12
	i      int
	j      int
}

func eProductParallel(
	wg *sync.WaitGroup,
	i int,
	j int,
	communication chan eResult,
	arguments ...*eArg,
) {
	routine := func() {
		defer wg.Done()

		result := eProduct(arguments...)

		communication <- eResult{result, i, j}
	}

	if _ParallelOptimization {
		go routine()
	} else {
		routine()
	}
}
