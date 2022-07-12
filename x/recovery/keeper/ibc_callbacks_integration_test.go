package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rhizomplatform/plateaus/app"
	"github.com/rhizomplatform/plateaus/testutil"
	claimtypes "github.com/rhizomplatform/plateaus/x/claims/types"
	"github.com/rhizomplatform/plateaus/x/recovery/types"
)

var _ = Describe("Recovery: Performing an IBC Transfer", Ordered, func() {
	coinXrz := sdk.NewCoin("axrz", sdk.NewInt(10000))
	coinOsmo := sdk.NewCoin("uosmo", sdk.NewInt(10))
	coinAtom := sdk.NewCoin("uatom", sdk.NewInt(10))

	var (
		sender, receiver       string
		senderAcc, receiverAcc sdk.AccAddress
		timeout                uint64
		claim                  claimtypes.ClaimsRecord
	)

	BeforeEach(func() {
		s.SetupTest()
	})

	Describe("from a non-authorized chain", func() {
		BeforeEach(func() {
			params := claimtypes.DefaultParams()
			params.AuthorizedChannels = []string{}
			s.PlateausChain.App.(*app.Plateaus).ClaimsKeeper.SetParams(s.PlateausChain.GetContext(), params)

			sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
			receiver = s.PlateausChain.SenderAccount.GetAddress().String()
			senderAcc = sdk.MustAccAddressFromBech32(sender)
			receiverAcc = sdk.MustAccAddressFromBech32(receiver)
		})
		It("should transfer and not recover tokens", func() {
			s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

			nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
			Expect(nativeXrz).To(Equal(coinXrz))
			ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
			Expect(ibcOsmo).To(Equal(sdk.NewCoin(uosmoIbcdenom, coinOsmo.Amount)))
		})
	})

	Describe("from an authorized, non-EVM chain (e.g. Osmosis)", func() {
		Describe("to a different account on Plateaus (sender != recipient)", func() {
			BeforeEach(func() {
				sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				receiver = s.PlateausChain.SenderAccount.GetAddress().String()
				senderAcc = sdk.MustAccAddressFromBech32(sender)
				receiverAcc = sdk.MustAccAddressFromBech32(receiver)
			})

			It("should transfer and not recover tokens", func() {
				s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

				nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
				Expect(nativeXrz).To(Equal(coinXrz))
				ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
				Expect(ibcOsmo).To(Equal(sdk.NewCoin(uosmoIbcdenom, coinOsmo.Amount)))
			})
		})

		Describe("to the sender's own eth_secp256k1 account on Plateaus (sender == recipient)", func() {
			BeforeEach(func() {
				sender = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				receiver = s.IBCOsmosisChain.SenderAccount.GetAddress().String()
				senderAcc = sdk.MustAccAddressFromBech32(sender)
				receiverAcc = sdk.MustAccAddressFromBech32(receiver)
			})

			Context("with disabled recovery parameter", func() {
				BeforeEach(func() {
					params := types.DefaultParams()
					params.EnableRecovery = false
					s.PlateausChain.App.(*app.Plateaus).RecoveryKeeper.SetParams(s.PlateausChain.GetContext(), params)
				})

				It("should not transfer or recover tokens", func() {
					s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)

					nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
					Expect(nativeXrz).To(Equal(coinXrz))
					ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
					Expect(ibcOsmo).To(Equal(sdk.NewCoin(uosmoIbcdenom, coinOsmo.Amount)))
				})
			})

			Context("with a sender's claims record", func() {
				Context("without completed actions", func() {
					BeforeEach(func() {
						amt := sdk.NewInt(int64(100))
						claim = claimtypes.NewClaimsRecord(amt)
						s.PlateausChain.App.(*app.Plateaus).ClaimsKeeper.SetClaimsRecord(s.PlateausChain.GetContext(), senderAcc, claim)
					})

					It("should not transfer or recover tokens", func() {
						// Prevent further funds from getting stuck
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)

						nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz).To(Equal(coinXrz))
						ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
					})
				})

				Context("with completed actions", func() {
					// Already has stuck funds
					BeforeEach(func() {
						amt := sdk.NewInt(int64(100))
						coins := sdk.NewCoins(sdk.NewCoin("axrz", sdk.NewInt(int64(75))))
						claim = claimtypes.NewClaimsRecord(amt)
						claim.MarkClaimed(claimtypes.ActionIBCTransfer)
						s.PlateausChain.App.(*app.Plateaus).ClaimsKeeper.SetClaimsRecord(s.PlateausChain.GetContext(), senderAcc, claim)

						// update the escrowed account balance to maintain the invariant
						err := testutil.FundModuleAccount(s.PlateausChain.App.(*app.Plateaus).BankKeeper, s.PlateausChain.GetContext(), claimtypes.ModuleName, coins)
						s.Require().NoError(err)

						// axrz & ibc tokens that originated from the sender's chain
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
						timeout = uint64(s.PlateausChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
					})

					It("should transfer tokens to the recipient and perform recovery", func() {
						// Escrow before relaying packets
						balanceEscrow := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), transfertypes.GetEscrowAddress("transfer", "channel-0"), "axrz")
						Expect(balanceEscrow).To(Equal(coinXrz))
						ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())

						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10000", "axrz", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Check that the axrz were recovered
						nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz.IsZero()).To(BeTrue())
						ibcPlateaus := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, axrzIbcdenom)
						Expect(ibcPlateaus).To(Equal(sdk.NewCoin(axrzIbcdenom, coinXrz.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo = s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))
					})

					It("should not claim/migrate/merge claims records", func() {
						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10000", "axrz", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						claimAfter, _ := s.PlateausChain.App.(*app.Plateaus).ClaimsKeeper.GetClaimsRecord(s.PlateausChain.GetContext(), senderAcc)
						Expect(claim).To(Equal(claimAfter))
					})
				})
			})

			Context("without a sender's claims record", func() {
				When("recipient has no ibc vouchers that originated from other chains", func() {
					It("should transfer and recover tokens", func() {
						// axrz & ibc tokens that originated from the sender's chain
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, coinOsmo.Denom, coinOsmo.Amount.Int64(), sender, receiver, 1)
						timeout = uint64(s.PlateausChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())

						// Escrow before relaying packets
						balanceEscrow := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), transfertypes.GetEscrowAddress("transfer", "channel-0"), "axrz")
						Expect(balanceEscrow).To(Equal(coinXrz))
						ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())

						// Relay both packets that were sent in the ibc_callback
						err := s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10000", "axrz", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Check that the axrz were recovered
						nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz.IsZero()).To(BeTrue())
						ibcPlateaus := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, axrzIbcdenom)
						Expect(ibcPlateaus).To(Equal(sdk.NewCoin(axrzIbcdenom, coinXrz.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo = s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))
					})
				})

				// Do not recover uatom sent from Cosmos when performing recovery through IBC transfer from Osmosis
				When("recipient has additional ibc vouchers that originated from other chains", func() {
					BeforeEach(func() {
						params := types.DefaultParams()
						params.EnableRecovery = false
						s.PlateausChain.App.(*app.Plateaus).RecoveryKeeper.SetParams(s.PlateausChain.GetContext(), params)

						// Send uatom from Cosmos to Plateaus
						s.SendAndReceiveMessage(s.pathCosmosPlateaus, s.IBCCosmosChain, coinAtom.Denom, coinAtom.Amount.Int64(), s.IBCCosmosChain.SenderAccount.GetAddress().String(), receiver, 1)

						params.EnableRecovery = true
						s.PlateausChain.App.(*app.Plateaus).RecoveryKeeper.SetParams(s.PlateausChain.GetContext(), params)
					})
					It("should not recover tokens that originated from other chains", func() {
						// Send uosmo from Osmosis to Plateaus
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 1)

						// Relay both packets that were sent in the ibc_callback
						timeout := uint64(s.PlateausChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err := s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10000", "axrz", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)

						// Axrz was recovered from user address
						nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz.IsZero()).To(BeTrue())
						ibcPlateaus := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, axrzIbcdenom)
						Expect(ibcPlateaus).To(Equal(sdk.NewCoin(axrzIbcdenom, coinXrz.Amount)))

						// Check that the uosmo were retrieved
						ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						// Check that the atoms were not retrieved
						ibcAtom := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, uatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(uatomIbcdenom, coinAtom.Amount)))

						// Repeat transaction from Osmosis to Plateaus
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 2)

						timeout = uint64(s.PlateausChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 3, timeout))
						s.Require().NoError(err)

						// No further tokens recovered
						nativeXrz = s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz.IsZero()).To(BeTrue())
						ibcPlateaus = s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, axrzIbcdenom)
						Expect(ibcPlateaus).To(Equal(sdk.NewCoin(axrzIbcdenom, coinXrz.Amount)))

						ibcOsmo = s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo = s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						ibcAtom = s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, uatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(uatomIbcdenom, coinAtom.Amount)))
					})
				})

				// Recover ibc/uatom that was sent from Osmosis back to Osmosis
				When("recipient has additional non-native ibc vouchers that originated from senders chains", func() {
					BeforeEach(func() {
						params := types.DefaultParams()
						params.EnableRecovery = false
						s.PlateausChain.App.(*app.Plateaus).RecoveryKeeper.SetParams(s.PlateausChain.GetContext(), params)

						s.SendAndReceiveMessage(s.pathOsmosisCosmos, s.IBCCosmosChain, coinAtom.Denom, coinAtom.Amount.Int64(), s.IBCCosmosChain.SenderAccount.GetAddress().String(), receiver, 1)

						// Send IBC transaction of 10 ibc/uatom
						transferMsg := transfertypes.NewMsgTransfer(s.pathOsmosisPlateaus.EndpointA.ChannelConfig.PortID, s.pathOsmosisPlateaus.EndpointA.ChannelID, sdk.NewCoin(uatomIbcdenom, sdk.NewInt(10)), sender, receiver, timeoutHeight, 0)
						_, err := s.IBCOsmosisChain.SendMsgs(transferMsg)
						s.Require().NoError(err) // message committed
						transfer := transfertypes.NewFungibleTokenPacketData("transfer/channel-1/uatom", "10", sender, receiver)
						packet := channeltypes.NewPacket(transfer.GetBytes(), 1, s.pathOsmosisPlateaus.EndpointA.ChannelConfig.PortID, s.pathOsmosisPlateaus.EndpointA.ChannelID, s.pathOsmosisPlateaus.EndpointB.ChannelConfig.PortID, s.pathOsmosisPlateaus.EndpointB.ChannelID, timeoutHeight, 0)
						// Receive message on the plateaus side, and send ack
						err = s.pathOsmosisPlateaus.RelayPacket(packet)
						s.Require().NoError(err)

						// Check that the ibc/uatom are available
						osmoIBCAtom := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uatomOsmoIbcdenom)
						s.Require().Equal(osmoIBCAtom.Amount, coinAtom.Amount)

						params.EnableRecovery = true
						s.PlateausChain.App.(*app.Plateaus).RecoveryKeeper.SetParams(s.PlateausChain.GetContext(), params)
					})
					It("should not recover tokens that originated from other chains", func() {
						s.SendAndReceiveMessage(s.pathOsmosisPlateaus, s.IBCOsmosisChain, "uosmo", 10, sender, receiver, 2)

						// Relay packets that were sent in the ibc_callback
						timeout := uint64(s.PlateausChain.GetContext().BlockTime().Add(time.Hour * 4).Add(time.Second * -20).UnixNano())
						err := s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10000", "axrz", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 1, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/transfer/channel-1/uatom", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 2, timeout))
						s.Require().NoError(err)
						err = s.pathOsmosisPlateaus.RelayPacket(CreatePacket("10", "transfer/channel-0/uosmo", sender, receiver, "transfer", "channel-0", "transfer", "channel-0", 3, timeout))
						s.Require().NoError(err)

						// Axrz was recovered from user address
						nativeXrz := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), senderAcc, "axrz")
						Expect(nativeXrz.IsZero()).To(BeTrue())
						ibcPlateaus := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, axrzIbcdenom)
						Expect(ibcPlateaus).To(Equal(sdk.NewCoin(axrzIbcdenom, coinXrz.Amount)))

						// Check that the uosmo were recovered
						ibcOsmo := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uosmoIbcdenom)
						Expect(ibcOsmo.IsZero()).To(BeTrue())
						nativeOsmo := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), receiverAcc, "uosmo")
						Expect(nativeOsmo).To(Equal(coinOsmo))

						// Check that the ibc/uatom were retrieved
						osmoIBCAtom := s.PlateausChain.App.(*app.Plateaus).BankKeeper.GetBalance(s.PlateausChain.GetContext(), receiverAcc, uatomOsmoIbcdenom)
						Expect(osmoIBCAtom.IsZero()).To(BeTrue())
						ibcAtom := s.IBCOsmosisChain.GetSimApp().BankKeeper.GetBalance(s.IBCOsmosisChain.GetContext(), senderAcc, uatomIbcdenom)
						Expect(ibcAtom).To(Equal(sdk.NewCoin(uatomIbcdenom, sdk.NewInt(10))))
					})
				})
			})
		})
	})
})
