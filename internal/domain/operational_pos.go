package domain

// cspell: ignore fdapi fdpkg eles

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	fdapi "github.com/pinguo-icc/field-definitions/api"
	fdpkg "github.com/pinguo-icc/field-definitions/pkg"
	oppapi "github.com/pinguo-icc/operational-positions-svc/api"
	"golang.org/x/text/language"
)

func NewParserFactory(c fdapi.FieldDefinitionsClient) *fdpkg.ParserFactory {
	return fdpkg.NewParserFactory(c, fdpkg.WithTTL(60))
}

var _ json.Marshaler = (*Activity)(nil)

// Activity an activity with contents
type Activity struct {
	ID        string
	PID       string
	RootID    string
	FieldCode string
	Name      string
	Period    period

	b, l []fdpkg.E
}

type period struct {
	Begin, End int64
}

var parseOpts = []fdpkg.ParseOption{
	fdpkg.WithConcise(true),
}

func (a *Activity) ParseContents(parser *fdpkg.Parser, lm language.Matcher, contents *fdpkg.FieldsCollection) (err error) {
	a.b, a.l, err = parser.Parse(lm, contents, parseOpts...)
	return
}

func (a *Activity) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.Grow(512)

	buf.WriteByte('{')
	{
		a.writeBaseKV(buf, "id", a.ID)
		a.writeBaseKV(buf, "pid", a.PID)
		a.writeBaseKV(buf, "rootId", a.RootID)
		a.writeBaseKV(buf, "fieldCode", a.FieldCode)
		a.writeBaseKV(buf, "name", a.Name)

		a.writePeriod(buf)

		if err := a.writeEles(buf, a.b); err != nil {
			return nil, err
		}
		if err := a.writeEles(buf, a.l); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (a *Activity) writeBaseKV(buf *bytes.Buffer, k, v string) {
	buf.WriteByte('"')
	buf.WriteString(k)
	buf.WriteString(`":"`)
	buf.WriteString(v)
	buf.WriteByte('"')
	buf.WriteByte(',')
}

func (a *Activity) writePeriod(buf *bytes.Buffer) {
	buf.WriteString(`"period":{"begin":`)
	buf.WriteString(strconv.FormatInt(a.Period.Begin, 10))

	buf.WriteString(`,"end":`)
	buf.WriteString(strconv.FormatInt(a.Period.End, 10))

	buf.WriteByte('}')
}

func (a *Activity) writeEles(buf *bytes.Buffer, data []fdpkg.E) error {
	if len(data) == 0 {
		return nil
	}

	en := json.NewEncoder(buf)

	for _, v := range data {
		buf.WriteByte(',')
		buf.WriteByte('"')
		buf.WriteString(v.Key)
		buf.WriteByte('"')
		buf.WriteByte(':')

		if err := en.Encode(v.Value); err != nil {
			return err
		}
	}
	return nil
}

type ActivitiesParser struct {
	pFac *fdpkg.ParserFactory
}

func NewActivitiesParser(p *fdpkg.ParserFactory) *ActivitiesParser {
	return &ActivitiesParser{pFac: p}
}

func (ap *ActivitiesParser) Parse(ctx context.Context, lm language.Matcher, data map[string]*oppapi.PlacingResponse_Activities) (map[string][]*Activity, error) {
	fps, err := ap.getFieldParser(ctx, data)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]*Activity, len(data))
	for k, v := range data {
		vv := v.Data
		if len(vv) > 0 {
			activities := make([]*Activity, len(vv))
			for i, ac := range vv {
				if fps[ac.FieldDefCode] == nil {
					return nil, fmt.Errorf("domain: ActivitiesParser#Parse field parse not found, fieldDefCode=%s", ac.FieldDefCode)
				}

				contents := new(fdpkg.FieldsCollection)
				if err := contents.Unmarshal(ac.Contents); err != nil {
					return nil, err
				}

				tmp := &Activity{
					ID:        ac.Id,
					PID:       ac.Pid,
					RootID:    ac.RootId,
					FieldCode: ac.FieldDefCode,
					Name:      ac.Name,
					Period: period{
						Begin: ac.Period.GetBegin(),
						End:   ac.Period.GetEnd(),
					},
				}

				if err := tmp.ParseContents(fps[ac.FieldDefCode], lm, contents); err != nil {
					return nil, err
				}
				activities[i] = tmp
			}

			res[k] = activities
		}
	}

	return res, nil
}

func (ap *ActivitiesParser) getFieldParser(ctx context.Context, data map[string]*oppapi.PlacingResponse_Activities) (map[string]*fdpkg.Parser, error) {
	dataset := make(map[string]struct{}, 8)
	fDefIDs := make([]string, 0, 8)
	for _, v := range data {
		for _, vv := range v.Data {
			if _, ok := dataset[vv.FieldDefCode]; ok {
				continue
			}

			dataset[vv.FieldDefCode] = struct{}{}
			fDefIDs = append(fDefIDs, vv.FieldDefCode)
		}
	}

	if len(fDefIDs) == 0 {
		return nil, nil
	}

	return ap.pFac.MGet(ctx, fDefIDs)
}
