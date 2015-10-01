package shadow_test
import (
	"testing"
	"github.com/akaspin/bar/shadow"
	"github.com/stretchr/testify/assert"
	"github.com/akaspin/bar/fixtures"
	"golang.org/x/crypto/sha3"
)

func Test_SHA3(t *testing.T) {

	hasher1 := sha3.New256()
	hasher1.Write([]byte("mama_myla_ramu"))

	hasher2 := sha3.New256()
	hasher2.Write([]byte("mama_"))
	hasher2.Write([]byte("myla_"))
	hasher2.Write([]byte("ramu"))

	assert.Equal(t, hasher1.Sum(nil), hasher2.Sum(nil))
}

func Test_Shadow_ToStringFull(t *testing.T) {
	id := "ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c"
	expect := fixtures.FixStream(`BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
		`)

	sh := &shadow.Shadow{
		false,
		"0.0.1",
		id,
		1234,
		[]shadow.Chunk{
			shadow.Chunk{id, 1234, 0},
		},
	}

	assert.Equal(t, expect, (*sh).String())
}

func Test_Shadow_New1(t *testing.T) {
	in := `BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
		`
	m, err := shadow.New(fixtures.CleanInput(in))
	assert.NoError(t, err)
	assert.Equal(t, fixtures.FixStream(in), (*m).String())
}

func Test_Shadow_NewFromBLOB_20M(t *testing.T)  {
	bn := fixtures.MakeBLOB(t, 1024 * 1024 * 20 + 5)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.FixStream(`BAR:SHADOW

	version 0.1.0
	id 9e39ad7cf632a038a5a2e0f9144f6ea4aff04ff11803c169cb24f60e56444f08
	size 20971525


	id 1d0debc4d598f7dc39cb53fb5e5ace3e310316231361bba6ef91e834cadd54f6
	size 2097152
	offset 0

	id a3ab965376a5e20a9b15a456ec9400cb16bd2c2570f2a0d53e249ea864307fbe
	size 2097152
	offset 2097152

	id 19493ebb0566afa93c72f0d4eaad59fa5f0326b57f4479e3a37a5d71af52e736
	size 2097152
	offset 4194304

	id c629e80e2205bfa680e0946a4fe4ac20ccb64c7b394d6486dcf1310f4b5dd9b6
	size 2097152
	offset 6291456

	id a6763718c0d3ef294d41c0c8ebe674ee2dac8432290aee1ff4d5e06bed33ad1c
	size 2097152
	offset 8388608

	id 48fb6b5457ce34099ce4f792519c3bc0d0cbdae006b8fc48e9d28850c07cd3fc
	size 2097152
	offset 10485760

	id e31ac3259e3ff446cc9a9e800104fa9d5f5cd74732bae954a6c52a718c4f5a58
	size 2097152
	offset 12582912

	id 4ff1c7f4d3b2c1ba53d2fa853b0b41af0faa9ce0e10e3c7bca9617eb48340af4
	size 2097152
	offset 14680064

	id beb0c46be81961250d2f84392bec70332d6cfa95011c29fa5d1cc176b5aa4feb
	size 2097152
	offset 16777216

	id cc50781f234ed5c6c7b337ec900491b5d8bb765c5769e1836c00e9ec0e43ce6b
	size 2097152
	offset 18874368

	id d9dc02c791744ee3db07776c1a149030c925ff2717ac07bde1fd10da475ac7e6
	size 5
	offset 20971520
	`), (*m).String())
}

func Test_Shadow_NewFromBLOB_2M(t *testing.T)  {
	bn := fixtures.MakeBLOB(t, 1024 * 1024 * 2 + 467)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.FixStream(`BAR:SHADOW

	version 0.1.0
	id fd76eb2f9866a12c6c3a2f884d5350b38319bc510106a7ba78789cc5507158b8
	size 2097619


	id 1d0debc4d598f7dc39cb53fb5e5ace3e310316231361bba6ef91e834cadd54f6
	size 2097152
	offset 0

	id 10c1a5bf6dee30935c0528049c4ccea0fa7f7d9d4d50fd361470b0affb0553f4
	size 467
	offset 2097152
	`), (*m).String())
}

func Test_Shadow_NewFromBLOB_2K(t *testing.T)  {
	bn := fixtures.MakeBLOB(t, 1024 * 2)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.FixStream(`BAR:SHADOW

		version 0.1.0
		id be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2
		size 2048


		id be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2
		size 2048
		offset 0
		`), (*m).String())
}

func Test_Shadow_NewFromBLOB_3b(t *testing.T)  {
	bn := fixtures.MakeBLOB(t, 3)
	defer fixtures.KillBLOB(bn)

	m, err := fixtures.NewShadowFromFile(bn)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.FixStream(`BAR:SHADOW

		version 0.1.0
		id 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253
		size 3


		id 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253
		size 3
		offset 0
		`), (*m).String())
}

func Benchmark_Shadow_NewFromBLOB_500MB(b *testing.B)  {
	bn, err := fixtures.MakeBLOBPure(1024 * 1024 * 500)
	if err != nil {
		b.Fail()
	}
	defer fixtures.KillBLOB(bn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = fixtures.NewShadowFromFile(bn)
	}

}
