package storage_test

import (
	"github.com/akaspin/bar/proto"
	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Storage_UUID(t *testing.T) {
	u, _ := uuid.NewV4()

	var id proto.ID
	err := (&id).UnmarshalBinary((*u)[:])
	assert.NoError(t, err)

}
