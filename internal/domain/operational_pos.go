package domain

// cspell: ignore fdapi fdpkg eles

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"

	fdapi "github.com/pinguo-icc/field-definitions/api"
	fdpkg "github.com/pinguo-icc/field-definitions/pkg"
	"github.com/pinguo-icc/kratos-library/v2/trace"
	oppapi "github.com/pinguo-icc/operational-positions-svc/api"
	"golang.org/x/text/language"
)

func NewParserFactory(c fdapi.FieldDefinitionsClient) *fdpkg.ParserFactory {
	return fdpkg.NewParserFactory(c, fdpkg.WithTTL(60))
}

var _ json.Marshaler = (*Activity)(nil)

type ActivityPlan struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Period     period      `json:"period"`
	Activities []*Activity `json:"activities"`
}

// Activity an activity with contents
type Activity struct {
	ID        string
	PID       string
	RootID    string
	TrackID   uint64
	FieldCode string
	Name      string
	Period    period

	b, l []fdpkg.E
}

type period struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
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
		a.writeBaseKV(buf, "trackId", strconv.FormatUint(a.TrackID, 10))
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
	pFac      *fdpkg.ParserFactory
	trFactory *trace.Factory
}

func NewActivitiesParser(p *fdpkg.ParserFactory, trFactory *trace.Factory) *ActivitiesParser {
	return &ActivitiesParser{pFac: p, trFactory: trFactory}
}

func (ap *ActivitiesParser) Parse(ctx context.Context, lm language.Matcher, data map[string]*oppapi.PlacingResponse_Plans) (map[string][]*ActivityPlan, error) {
	fps, err := ap.getFieldParser(ctx, data)
	if err != nil {
		return nil, err
	}

	formatActivity := func(ac *oppapi.PlacingResponse_Activity) (*Activity, error) {
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
			return nil, fmt.Errorf("parse activity failed, id=%s %w", ac.Id, err)
		}
		return tmp, nil
	}

	res := make(map[string][]*ActivityPlan, len(data))
	for posCode, pPlan := range data {
		outPlans := make([]*ActivityPlan, len(pPlan.Plans))
		for i, plan := range pPlan.Plans {
			outPlans[i] = &ActivityPlan{
				ID:   plan.Id,
				Name: plan.Name,
				Period: period{
					Begin: plan.Period.GetBegin(),
					End:   plan.Period.GetEnd(),
				},
				Activities: make([]*Activity, len(plan.Activities)),
			}
			for j, ac := range plan.Activities {
				tmp, err := formatActivity(ac)
				if err != nil {
					return nil, err
				} else {
					tmp.TrackID = ap.generateTrackID(plan.Id, plan.ContentId, tmp.ID)
					outPlans[i].Activities[j] = tmp
				}
			}
		}

		res[posCode] = outPlans
	}

	return res, nil
}

func (ap *ActivitiesParser) getFieldParser(ctx context.Context, data map[string]*oppapi.PlacingResponse_Plans) (map[string]*fdpkg.Parser, error) {
	dataset := make(map[string]struct{}, 8)
	fDefIDs := make([]string, 0, 8)
	for _, pPlans := range data {
		for _, p := range pPlans.Plans {
			for _, ac := range p.Activities {
				if _, ok := dataset[ac.FieldDefCode]; ok {
					continue
				}
				dataset[ac.FieldDefCode] = struct{}{}
				fDefIDs = append(fDefIDs, ac.FieldDefCode)
			}

		}
	}

	if len(fDefIDs) == 0 {
		return nil, nil
	}

	return ap.pFac.MGet(ctx, fDefIDs)
}

func (ap *ActivitiesParser) generateTrackID(planID, contentID, activityID string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(planID + contentID + activityID))
	return h.Sum64()
}
