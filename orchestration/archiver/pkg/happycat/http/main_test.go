package http

import (
	"archiver/pkg/internal/db"
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

type testQuerier struct {
	list     []*db.HappycatFact
	listByID map[uuid.UUID]*db.HappycatFact
}

func newTestQuerier(list []*db.HappycatFact) *testQuerier {
	listByID := make(map[uuid.UUID]*db.HappycatFact)

	for _, hcf := range list {
		listByID[hcf.ID] = hcf
	}

	return &testQuerier{
		list,
		listByID,
	}
}

func (tq *testQuerier) SaveHappycatFact(ctx context.Context, arg db.SaveHappycatFactParams) error {
	cf := &db.HappycatFact{ID: arg.ID, Fact: arg.Fact, CreatedAt: time.Now()}

	tq.list = append(tq.list, cf)
	tq.listByID[arg.ID] = cf

	return nil
}

func (tq *testQuerier) GetHappycatFact(ctx context.Context, id uuid.UUID) (*db.HappycatFact, error) {
	hcf, ok := tq.listByID[id]
	if !ok {
		return nil, NoHappyCatFactFoundError{}
	}

	return hcf, nil
}

func (tq *testQuerier) ListHappycatFacts(ctx context.Context) ([]*db.HappycatFact, error) {
	if len(tq.list) == 0 {
		return nil, NoHappyCatFactFoundError{}
	}

	return tq.list, nil
}
