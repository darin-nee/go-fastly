package fastly

import (
	"testing"
)

func TestClient_ListTokens(t *testing.T) {
	t.Parallel()

	var tokens []*Token
	var err error
	record(t, "tokens/list", func(c *Client) {
		tokens, err = c.ListTokens(&ListTokensInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 1 {
		t.Errorf("bad tokens: %v", tokens)
	}
}

func TestClient_ListCustomerTokens(t *testing.T) {
	t.Parallel()

	var tokens []*Token
	var err error
	record(t, "tokens/list_customer", func(c *Client) {
		tokens, err = c.ListCustomerTokens(&ListCustomerTokensInput{
			CustomerID: "XXXXXXXXXXXXXXXXXXXXXX",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 1 {
		t.Errorf("bad tokens: %v", tokens)
	}
}

func TestClient_GetTokenSelf(t *testing.T) {
	t.Parallel()

	var token *Token
	var err error
	record(t, "tokens/get_self", func(c *Client) {
		token, err = c.GetTokenSelf()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", token)
}

func TestClient_CreateToken(t *testing.T) {
	t.Parallel()

	input := &CreateTokenInput{
		Name:     ToPointer("my-test-token"),
		Scope:    ToPointer(GlobalScope),
		Username: ToPointer("XXXXXXXXXXXXXXXXXXXXXX"),
		Password: ToPointer("XXXXXXXXXXXXXXXXXXXXXX"),
	}

	var token *Token
	var err error
	record(t, "tokens/create", func(c *Client) {
		token, err = c.CreateToken(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if *token.Name != *input.Name {
		t.Errorf("returned invalid name, got %s, want %s", *token.Name, *input.Name)
	}
	if *token.Scope != *input.Scope {
		t.Errorf("returned invalid scope, got %s, want %s", *token.Scope, *input.Scope)
	}
}

func TestClient_DeleteToken(t *testing.T) {
	t.Parallel()

	input := &DeleteTokenInput{
		TokenID: "XXXXXXXXXXXXXXXXXXXXXX",
	}

	var err error
	record(t, "tokens/delete", func(c *Client) {
		err = c.DeleteToken(input)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteTokenSelf(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tokens/delete_self", func(c *Client) {
		err = c.DeleteTokenSelf()
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateAndBulkDeleteTokens(t *testing.T) {
	t.Parallel()

	var deleteErr error

	record(t, "tokens/create_and_bulk_delete", func(c *Client) {
		token1, err := c.CreateToken(&CreateTokenInput{
			Name:     ToPointer("my-test-token-1"),
			Scope:    ToPointer(GlobalScope),
			Username: ToPointer("testing@fastly.com"),
			Password: ToPointer("foobar"),
			Services: []string{"0Us63sb8R1BpWQBIhluncu", "7frORaFZvHgC6eRAJdA7kf"},
		})
		if err != nil {
			t.Fatal(err)
		}

		token2, err := c.CreateToken(&CreateTokenInput{
			Name:     ToPointer("my-test-token-2"),
			Scope:    ToPointer(GlobalScope),
			Username: ToPointer("testing@fastly.com"),
			Password: ToPointer("foobar"),
			Services: []string{"0Us63sb8R1BpWQBIhluncu", "7frORaFZvHgC6eRAJdA7kf"},
		})
		if err != nil {
			t.Fatal(err)
		}

		deleteErr = c.BatchDeleteTokens(&BatchDeleteTokensInput{
			Tokens: []*BatchToken{
				{ID: *token1.ID},
				{ID: *token2.ID},
			},
		})
	})

	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
}
