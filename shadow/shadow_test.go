package shadow_test
import (
	"testing"
	"github.com/akaspin/bar/shadow"
	"github.com/stretchr/testify/assert"
	"bytes"
	"os"
	"strings"
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

	sh := &shadow.Shadow{
		false,
		"0.0.1",
		id,
		1234,
		[]shadow.Chunk{
			shadow.Chunk{id, 1234, 0},
		},
	}

	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
		`), (*sh).String())
}

func Test_Shadow_ToStringShort(t *testing.T) {
	id := "ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c"

	sh := &shadow.Shadow{
		false,
		"0.0.1",
		id,
		1234,
		nil,
	}

	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		`), (*sh).String())
}

func Test_Shadow_FromManifestShort(t *testing.T) {
	in := `BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		`
	m := &shadow.Shadow{}
	err := m.FromManifest(bytes.NewReader(
		[]byte(in),
	), false)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(in), (*m).String())
	assert.False(t, m.HasChunks())
}

func Test_Shadow_FromManifestFull(t *testing.T) {
	in := `BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234


		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		offset 0
		`
	m := &shadow.Shadow{}
	err := m.FromManifest(bytes.NewReader(
		[]byte(in),
	), true)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(in), (*m).String())
	assert.True(t, m.HasChunks())
}

func Test_Shadow_FromFullManifestShort(t *testing.T) {
	in := `BAR:SHADOW

		version 0.1.0
		id ac934d9a88b42aa3b40ef7c81a9dee1aad5a2cddccb00ae6abab9c38095fc15c
		size 1234
		`
	m := &shadow.Shadow{}
	err := m.FromManifest(bytes.NewReader(
		[]byte(in),
	), false)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(in), (*m).String())
	assert.False(t, m.HasChunks())
}



func Test_Shadow_FromBLOB_20M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 20 + 5)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id 9e39ad7cf632a038a5a2e0f9144f6ea4aff04ff11803c169cb24f60e56444f08
		size 20971525


		id 7dd05dc961a9bfd1d7d3310e890d8be20a7f0667b17b1380a520adc491f202ce
		size 1048576
		offset 0

		id ce3ab016cb42aa12dc1eab3976ad3a49e134966e5cdf24c2c48f93b56239460e
		size 1048576
		offset 1048576

		id 39342fd4b4fe70f6b6640f67f86a8fbc50dc6fdf70a414fa517bd2d57073bddf
		size 1048576
		offset 2097152

		id 824b76f05b554f8ff29f5a0cf67bc173a47d3919050de53c9515e961b25f0060
		size 1048576
		offset 3145728

		id 5f2da039a71d1693b7df7b3eaaad25981cd50e733ea474c0243e47add598d451
		size 1048576
		offset 4194304

		id ec03eb3211ed6e9f94257c2be969ecda232f8b8b4d5c5ea40bc359623db2e709
		size 1048576
		offset 5242880

		id 5f64cd7323313a50291930e1c75253e0f102bfc21ba8816498ef2b5aac232fa3
		size 1048576
		offset 6291456

		id cc31b92e7bf09e7f8e02c8b23665f691d2ee6ec4cb72062209daa3e59efb3ebb
		size 1048576
		offset 7340032

		id 4c9a19b8eaf8c2beb7c129cd26534f0d577453681b1ff20b68dabbbdb3c8dc5c
		size 1048576
		offset 8388608

		id e09e11abc2960645073d73fc8dca221975d658a16a48d08f563a42bf5f20257b
		size 1048576
		offset 9437184

		id d40571b37c189a809f9eee35a33ce4d16c733505004b4ffb2eef6f94e10140c7
		size 1048576
		offset 10485760

		id 1e9ff301260c450c3061dc61cb019653b648572508fd24dc802a75eb6b0ddb3a
		size 1048576
		offset 11534336

		id 7e16a4047ddd72526f5957d02912108f1cb8873e497a361f3e9b49bdd65c1eba
		size 1048576
		offset 12582912

		id ce33106c66b53300d24f4dee48efdb63eaaa42eb0d900509330c59bf64cfb099
		size 1048576
		offset 13631488

		id d322a2b51b11a1c01f8a558020df3eed52ea6013cb753b2b1f0db03cb7ab5e0b
		size 1048576
		offset 14680064

		id 4eeba906e1845e0edd99d60753290760ffb9a92dee3a3387ebab6847ade2e59c
		size 1048576
		offset 15728640

		id 2f1607b9305e7f07c9bf1544cd2e04113f618c3eb55db1023475e546d9c6c2a6
		size 1048576
		offset 16777216

		id 8dba1480af0468641177b9d42f3c5f87fa114e05a3f98f881853eb6c6a70c370
		size 1048576
		offset 17825792

		id 80c2308266f0e83c0ac360124c6cbef8f0e414852ff24ac93cb86b28eb9e4ab9
		size 1048576
		offset 18874368

		id 3daed2017dd78450046d83c11906b07942baf8c122ba53ea7fb72d07f7d5cfe3
		size 1048576
		offset 19922944

		id d9dc02c791744ee3db07776c1a149030c925ff2717ac07bde1fd10da475ac7e6
		size 5
		offset 20971520
		`), (*m).String())
}

func Test_Shadow_FromBLOB_2M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 2 + 467)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id fd76eb2f9866a12c6c3a2f884d5350b38319bc510106a7ba78789cc5507158b8
		size 2097619


		id 7dd05dc961a9bfd1d7d3310e890d8be20a7f0667b17b1380a520adc491f202ce
		size 1048576
		offset 0

		id ce3ab016cb42aa12dc1eab3976ad3a49e134966e5cdf24c2c48f93b56239460e
		size 1048576
		offset 1048576

		id 10c1a5bf6dee30935c0528049c4ccea0fa7f7d9d4d50fd361470b0affb0553f4
		size 467
		offset 2097152
		`), (*m).String())
}

func Test_Shadow_FromBLOB_2K(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2
		size 2048


		id be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2
		size 2048
		offset 0
		`), (*m).String())
}

func Test_Shadow_FromBLOB_3b(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(3)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromBlob(r, true, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253
		size 3


		id 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253
		size 3
		offset 0
		`), (*m).String())
}

func Test_Shadow_FromAny_Manifest(t *testing.T) {
	in := `BAR:SHADOW

		version 0.1.0
		id 82783ef12d68fd4c57fd7a8d7e42e7b71fc0fd13e5e30d459f15bc64a298395c
		size 4
		`
	m := &shadow.Shadow{}
	err := m.FromAny(strings.NewReader(in), false, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(in), (*m).String())
}

func Test_Shadow_FromAny_20M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 20 + 5)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r, false, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id 9e39ad7cf632a038a5a2e0f9144f6ea4aff04ff11803c169cb24f60e56444f08
		size 20971525
		`), (*m).String())
}

func Test_Shadow_FromAny_2M(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 1024 * 2 + 467)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r, false, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id fd76eb2f9866a12c6c3a2f884d5350b38319bc510106a7ba78789cc5507158b8
		size 2097619
		`), (*m).String())
}

func Test_Shadow_FromAny_2K(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(1024 * 2)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r, false, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id be4215176932949d887fa82241bbe0b03a44dc16ee2d23eedbc973e511ae8bb2
		size 2048
		`), (*m).String())
}

func Test_Shadow_FromAny_3b(t *testing.T)  {
	bn, err := fixtures.MakeBLOB(3)
	assert.NoError(t, err)
	defer fixtures.KillBLOB(bn)

	m := &shadow.Shadow{}
	r, err := os.Open(bn)
	assert.NoError(t, err)
	defer r.Close()

	err = m.FromAny(r, false, shadow.CHUNK_SIZE)
	assert.NoError(t, err)
	assert.Equal(t, fixtures.CleanManifest(`BAR:SHADOW

		version 0.1.0
		id 1186d49a4ad620618f760f29da2c593b2ec2cc2ced69dc16817390d861e62253
		size 3
		`), (*m).String())
}
