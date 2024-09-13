package traefik_plugin_sec_hasura

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServeBatchQuery(t *testing.T) {
	tests := []struct {
		Name       string
		Query      string
		StatusCode int
	}{
		{
			Name: "SimpleQuery",
			Query: `{
  "query": "query {  blocks(limit: 1) { height } }",
  "variables": null
}`,
			StatusCode: 200,
		},

		{
			Name: "BatchQuery",
			Query: `[
 {"operationName":"ProposalsActive","variables":{"where":{"status":{"_neq":"PROPOSAL_STATUS_INVALID"}},"limit":16,"offset":0},"query":"query ProposalsActive($limit: Int!, $offset: Int!, $where: proposal_bool_exp = {}) {\n  all_proposals: proposal(\n    order_by: {active_first_order: asc}\n    limit: $limit\n    offset: $offset\n    where: $where\n  ) {\n    content\n    deposit_end_time\n    description\n    id\n    proposal_type\n    proposal_deposits {\n      amount\n      depositor_address\n      __typename\n    }\n    proposal_votes_aggregate(where: {is_valid: {_eq: true}}) {\n      aggregate {\n        count\n        __typename\n      }\n      __typename\n    }\n    proposer_address\n    status\n    submit_time\n    title\n    voting_end_time\n    voting_start_time\n    __typename\n  }\n  proposal_aggregate(where: $where) {\n    aggregate {\n      count\n      __typename\n    }\n    __typename\n  }\n}"},{"operationName":"ProposalsActive","variables":{"where":{"status":{"_neq":"PROPOSAL_STATUS_INVALID"}},"limit":16,"offset":0},"query":"query ProposalsActive($limit: Int!, $offset: Int!, $where: proposal_bool_exp = {}) {\n  all_proposals: proposal(\n    order_by: {active_first_order: asc}\n    limit: $limit\n    offset: $offset\n    where: $where\n  ) {\n    content\n    deposit_end_time\n    description\n    id\n    proposal_type\n    proposal_deposits {\n      amount\n      depositor_address\n      __typename\n    }\n    proposal_votes_aggregate(where: {is_valid: {_eq: true}}) {\n      aggregate {\n        count\n        __typename\n      }\n      __typename\n    }\n    proposer_address\n    status\n    submit_time\n    title\n    voting_end_time\n    voting_start_time\n    __typename\n  }\n  proposal_aggregate(where: $where) {\n    aggregate {\n      count\n      __typename\n    }\n    __typename\n  }\n}"}
]`,
			StatusCode: 403,
		},
	}

	cfg := CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			handler, err := New(ctx, next, cfg, "traefik-plugin-sec-hasura")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			data := bytes.NewBuffer([]byte(test.Query))

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost/v1/graphql", data)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			require.Equal(t, recorder.Code, test.StatusCode)

		})
	}
}
