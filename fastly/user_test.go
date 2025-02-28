package fastly

import (
	"testing"
)

func TestClient_UsersCurrent(t *testing.T) {
	t.Parallel()

	var err error
	var u *User
	record(t, "users/get_current_user", func(c *Client) {
		u, err = c.GetCurrentUser()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", u)
}

func TestClient_Users(t *testing.T) {
	t.Parallel()

	fixtureBase := "users/"
	login := "go-fastly-test+user+20221104@example.com"

	// Create
	//
	// NOTE: When recreating the fixtures, update the login.
	var err error
	var u *User
	record(t, fixtureBase+"create", func(c *Client) {
		u, err = c.CreateUser(&CreateUserInput{
			Login: ToPointer(login),
			Name:  ToPointer("test user"),
			Role:  ToPointer("engineer"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			_ = c.DeleteUser(&DeleteUserInput{
				ID: *u.ID,
			})
		})
	}()

	if *u.Login != login {
		t.Errorf("bad login: %v", *u.Login)
	}

	if *u.Name != "test user" {
		t.Errorf("bad name: %v", *u.Name)
	}

	if *u.Role != "engineer" {
		t.Errorf("bad role: %v", *u.Role)
	}

	// List
	var us []*User
	record(t, fixtureBase+"list", func(c *Client) {
		us, err = c.ListCustomerUsers(&ListCustomerUsersInput{
			CustomerID: *u.CustomerID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) < 1 {
		t.Errorf("bad users: %v", us)
	}

	// Get
	var nu *User
	record(t, fixtureBase+"get", func(c *Client) {
		nu, err = c.GetUser(&GetUserInput{
			ID: *u.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *u.Name != *nu.Name {
		t.Errorf("bad name: %q (%q)", *u.Name, *nu.Name)
	}

	// Update
	var uu *User
	record(t, fixtureBase+"update", func(c *Client) {
		uu, err = c.UpdateUser(&UpdateUserInput{
			ID:   *u.ID,
			Name: ToPointer("updated user"),
			Role: ToPointer("superuser"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uu.Name != "updated user" {
		t.Errorf("bad name: %q", *uu.Name)
	}
	if *uu.Role != "superuser" {
		t.Errorf("bad role: %q", *uu.Role)
	}

	// Reset Password
	//
	// NOTE: This integration test can fail due to reCAPTCHA.
	// Which means you might have to manually correct the fixtures 😬
	record(t, fixtureBase+"reset_password", func(c *Client) {
		err = c.ResetUserPassword(&ResetUserPasswordInput{
			Login: *uu.Login,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeleteUser(&DeleteUserInput{
			ID: *u.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListCustomerUsers_validation(t *testing.T) {
	var err error
	_, err = testClient.ListCustomerUsers(&ListCustomerUsersInput{
		CustomerID: "",
	})
	if err != ErrMissingCustomerID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetUser_validation(t *testing.T) {
	var err error
	_, err = testClient.GetUser(&GetUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateUser_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateUser(&UpdateUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteUser_validation(t *testing.T) {
	err := testClient.DeleteUser(&DeleteUserInput{
		ID: "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ResetUser_validation(t *testing.T) {
	err := testClient.ResetUserPassword(&ResetUserPasswordInput{
		Login: "",
	})
	if err != ErrMissingLogin {
		t.Errorf("bad error: %s", err)
	}
}
