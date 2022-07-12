package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/rhizomplatform/plateaus/x/claims/types"
)

func (suite *KeeperTestSuite) TestsClaimsRecords() {
	addr1, err := sdk.AccAddressFromBech32("xrz1777uf283z705lmeztqxpme7w7e6g4qjlwdy2ch")
	suite.Require().NoError(err)
	addr2, err := sdk.AccAddressFromBech32("xrz1hmvjlsellna3khhfn2n5d2mmxucr9jmyqul9s0")
	suite.Require().NoError(err)

	cr1 := types.NewClaimsRecord(sdk.NewInt(1000))
	cr2 := types.NewClaimsRecord(sdk.NewInt(200))
	cr2.MarkClaimed(types.ActionDelegate)

	expRecords := []types.ClaimsRecordAddress{
		{
			Address:                addr2.String(),
			InitialClaimableAmount: cr2.InitialClaimableAmount,
			ActionsCompleted:       cr2.ActionsCompleted,
		},
		{
			Address:                addr1.String(),
			InitialClaimableAmount: cr1.InitialClaimableAmount,
			ActionsCompleted:       cr1.ActionsCompleted,
		},
	}

	suite.app.ClaimsKeeper.SetClaimsRecord(suite.ctx, addr1, cr1)
	suite.app.ClaimsKeeper.SetClaimsRecord(suite.ctx, addr2, cr2)

	records := suite.app.ClaimsKeeper.GetClaimsRecords(suite.ctx)
	suite.Require().Equal(expRecords, records)
}
