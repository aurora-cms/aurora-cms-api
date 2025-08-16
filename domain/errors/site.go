package errors

import "errors"

var ErrSiteNameEmpty = errors.New("site name cannot be empty")
var ErrSiteDomainEmpty = errors.New("site domain cannot be empty")
var ErrSitePageEmpty = errors.New("site page cannot be empty")
var ErrSitePageWithSlugAlreadyExists = errors.New("site page with slug already exists")
var ErrSitePageNotFound = errors.New("site page not found")
var ErrSiteEmpty = errors.New("site cannot be empty")
