package helix

type GetDropEntitlementsParams struct {
	ID     string `query:"id"`
	UserID string `query:"user_id"`
	GameID string `query:"game_id"`
	After  string `query:"after"`
	First  int    `query:"first,20"` // Limit 1000
}

type Entitlement struct {
	ID        string `json:"id"`
	BenefitID string `json:"benefit_id"`
	Timestamp Time   `json:"timestamp"`
	UserID    string `json:"user_id"`
	GameID    string `json:"game_id"`
}

type ManyEntitlements struct {
	Entitlements []Entitlement `json:"data"`
}

type ManyEntitlementsWithPagination struct {
	ManyEntitlements
	Pagination `json:"pagination"`
}

type GetDropsEntitlementsResponse struct {
	ResponseCommon
	Data ManyEntitlementsWithPagination
}

// GetDropsEntitlements returns a list of entitlements, which have been awarded to users by your organization.
// Filtering by UserID returns all of the entitlements related to that specific user.
// Filtering by GameID returns all of the entitlements related to that game.
// Filtering by GameID and UserID returns all of the entitlements related to that game and that user.
// Entitlements are digital items that users are entitled to use. Twitch entitlements are granted based on viewership
// engagement with a content creator, based on the game developers' campaign.
func (c *Client) GetDropsEntitlements(params *GetDropEntitlementsParams) (*GetDropsEntitlementsResponse, error) {
	resp, err := c.get("/entitlements/drops", &ManyEntitlementsWithPagination{}, params)
	if err != nil {
		return nil, err
	}

	entitlements := &GetDropsEntitlementsResponse{}
	resp.HydrateResponseCommon(&entitlements.ResponseCommon)
	entitlements.Data.Entitlements = resp.Data.(*ManyEntitlementsWithPagination).Entitlements
	entitlements.Data.Pagination = resp.Data.(*ManyEntitlementsWithPagination).Pagination

	return entitlements, nil
}

type UpdateDropEntitlementsParams struct {
	EntitlementIDs     []string              `query:"entitlement_ids"`    // Limit 100
	FullfillmentStatus EntitlementCodeStatus `query:"fulfillment_status"` // CLAIMED or FULLFILLED
}

type EntitlementStatus struct {
	Status EntitlementCodeStatus `json:"status"` // SUCCESS, INVALID_ID, NOT_FOUND, UNAUTHORIZED, UPDATE_FAILED
	IDs    []string              `json:"ids"`
}

type ManyEntitlementStatuses struct {
	EntitlementStatuses []EntitlementStatus `json:"data"`
}

type UpdateDropsEntitlementsResponse struct {
	ResponseCommon
	Data ManyEntitlementStatuses
}

// UpdateDropsEntitlements updates the fulfillment status on a set of Drops entitlements, specified by their entitlement IDs.
//
// Requires User OAuth Token or App Access Token
//
// The client ID associated with the access token must have ownership of the game: Client ID > Organization ID > Game ID
//
// App Access OAuth Token returns all entitlements with benefits owned by your organization.
// User OAuth Token returns all entitlements owned by that user with benefits owned by your organization.
func (c *Client) UpdateDropsEntitlements(params *UpdateDropEntitlementsParams) (*UpdateDropsEntitlementsResponse, error) {
	resp, err := c.patchAsJSON("/entitlements/drops", &ManyEntitlementStatuses{}, params)
	if err != nil {
		return nil, err
	}

	entitlements := &UpdateDropsEntitlementsResponse{}
	resp.HydrateResponseCommon(&entitlements.ResponseCommon)
	entitlements.Data.EntitlementStatuses = resp.Data.(*ManyEntitlementStatuses).EntitlementStatuses

	return entitlements, nil
}
