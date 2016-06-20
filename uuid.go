package sirpent

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

type UUID struct {
	satori_uuid uuid.UUID
}

func NewUUID() UUID {
	return UUID{satori_uuid: uuid.NewV4()}
}

func UUIDFromString(id_str string) (UUID, error) {
	satori_uuid, err := uuid.FromString(id_str)
	return UUID{satori_uuid: satori_uuid}, err
}

func (id UUID) String() string {
	return id.satori_uuid.String()
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *UUID) UnmarshalJSON(b []byte) error {
	var id_str string
	err := json.Unmarshal(b, &id_str)
	if err == nil {
		var id2 UUID
		id2, err = UUIDFromString(id_str)
		id.satori_uuid = id2.satori_uuid
	}
	return err
}
