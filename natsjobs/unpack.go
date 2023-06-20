package natsjobs

import (
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// consume all
	auto string = "deduced_by_rr"
)

func (c *Driver) unpack(data []byte, item *Item) error {
	err := json.Unmarshal(data, item)
	if err != nil {
		if c.consumeAll {
			c.log.Debug("unmarshal error", zap.Error(err))

			uid := uuid.NewString()
			c.log.Debug("get raw payload", zap.String("assigned ID", uid))

			if isJSONEncoded(data) != nil {
				data, err = json.Marshal(data)
				if err != nil {
					return err
				}
			}

			*item = Item{
				Job:     auto,
				Ident:   uid,
				Payload: bytesToStr(data),
				headers: make(map[string][]string, 2),
				Options: &Options{
					Priority: 10,
					Pipeline: (*c.pipeline.Load()).Name(),
					Queue:    c.stream,
				},
			}

			return nil
		}

		return err
	}

	return nil
}

func isJSONEncoded(data []byte) error {
	var a any
	return json.Unmarshal(data, &a)
}
