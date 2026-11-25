// Package githubbilling reads the billing usage rows of the GitHub user that
// gh is logged in as, for one month, through two read-only gh api calls. It
// needs the gh scope user and never changes the login or its scopes.
package githubbilling

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/ghapi"
)

// Item is one billing usage row: one day, product, SKU and repository.
type Item struct {
	Product, SKU, UnitType, Repository string
	Quantity                           float64
}

// Usage is the login and its billing rows of one month.
type Usage struct {
	Login string
	Items []Item
}

var login = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

// Read returns the billing usage of the logged-in gh user in year and month.
// Errors carry one of the sentinels above and never gh's output.
func Read(ctx context.Context, runner process.Runner, year int, month time.Month) (Usage, error) {
	var user struct {
		Login string `json:"login"`
	}
	if err := ghapi.Get(ctx, runner, "/user", nil, &user); err != nil {
		if errors.Is(err, ghapi.ErrNotPermitted) {
			return Usage{}, fmt.Errorf("%w: gh could not read the logged-in user", ghapi.ErrAuth)
		}
		return Usage{}, err
	}
	if !login.MatchString(user.Login) {
		return Usage{}, fmt.Errorf("%w: unreadable login", ghapi.ErrFailed)
	}
	var body struct {
		UsageItems []struct {
			Product        string  `json:"product"`
			SKU            string  `json:"sku"`
			UnitType       string  `json:"unitType"`
			RepositoryName string  `json:"repositoryName"`
			Quantity       float64 `json:"quantity"`
		} `json:"usageItems"`
	}
	path := fmt.Sprintf("/users/%s/settings/billing/usage?year=%d&month=%d", user.Login, year, int(month))
	if err := ghapi.Get(ctx, runner, path, nil, &body); err != nil {
		if errors.Is(err, ghapi.ErrNotPermitted) {
			return Usage{}, fmt.Errorf("%w: billing usage needs the gh scope user", ghapi.ErrScope)
		}
		return Usage{}, err
	}
	usage := Usage{Login: user.Login, Items: []Item{}}
	for _, i := range body.UsageItems {
		usage.Items = append(usage.Items, Item{Product: i.Product, SKU: i.SKU, UnitType: i.UnitType, Repository: i.RepositoryName, Quantity: i.Quantity})
	}
	return usage, nil
}
