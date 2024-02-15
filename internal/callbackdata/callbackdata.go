package callbackdata

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type CallBackData struct {
	prefix  string
	builder *strings.Builder
}

const Separator = ":"

func NewCallBackDataFromData(dataStr string) CallBackData {
	data := strings.Split(dataStr, Separator)
	prefix := withSeparator(data[0])
	c := CallBackData{prefix: prefix, builder: &strings.Builder{}}

	if len(data) > 1 {
		for i, d := range data[1:] {
			if i != len(data[1:])-1 {
				c.builder.WriteString(withSeparator(d))
				continue
			}
			c.builder.WriteString(d)
		}
	}

	return c
}

func NewCallBackData(prefix string) CallBackData {
	c := CallBackData{prefix: withSeparator(prefix), builder: &strings.Builder{}}
	c.builder.WriteString(c.prefix)
	return c
}

func (c *CallBackData) Add(data string) *CallBackData {
	c.builder.WriteString(withSeparator(data))
	return c
}

func (c *CallBackData) String() string {
	return c.builder.String()
}

func (c *CallBackData) IntData() []int {
	data := c.Data()
	var result []int
	for _, d := range data {
		i, err := strconv.Atoi(d)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}
		result = append(result, i)
	}
	return result
}

func (c *CallBackData) Data() []string {
	data := strings.Split(c.String(), Separator)

	for i, d := range data {
		if d == "" {
			data = append(data[:i], data[i+1:]...)
		}
	}

	return data
}

func withSeparator(data string) string {
	return data + Separator
}
