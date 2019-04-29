package types

//fzmsupply action
const (
	FzmsupplyActionType = iota
	FzmsupplyInitPlatformAction = 1; //平台初始化
	FzmsupplyCreateRoleAction = 2; //创建后台角色
	FzmsupplyUpdateRoleAction = 3; //修改角色权限
	FzmsupplyDeleteRoleAction = 4; //删除角色
	FzmsupplyCreateOfficialAction = 5; //创建后台人员
	FzmsupplyUpdateOfficialAction = 6; //修改后台人员权限
	FzmsupplyDeleteOfficialAction = 7; //删除后台人员
	FzmsupplyCreateCompanyAction = 8; //企业用户注册
	FzmsupplyCreateAssetAction = 9; //企业录入或修改资产
	FzmsupplyExamineAssetAction = 10; //资产审核录入
	FzmsupplyCancelAssetAction = 11; //企业撤销资产登记
	FzmsupplySetCreditAction = 12; //设置企业信用总额度
	FzmsupplyExamineDepositAction = 13; //审核企业入金
	FzmsupplyApplyWithdrawAction = 14; //企业申请出金
	FzmsupplyExamineWithdrawAction = 15; //企业出金初审
	FzmsupplyReviewWithdrawAction = 16; //企业出金复审
	FzmsupplyMortgageAction = 17; //申请质押资产
	FzmsupplyCancelMortgageAction = 18; //撤销质押
	FzmsupplyConfirmMortgageAction = 19; //摘取质押
	FzmsupplyRedeemAction = 20; //赎回资产
	FzmsupplyOverdueAction = 21; //设置逾期质押
	FzmsupplyConfirmWithdrawAction = 22; //确认出金结果
	FzmsupplyReceiveAction = 23; //企业签收资产
	FzmsupplySellAction = 24;//挂牌
	FzmsupplyCancelAction = 25; //撤销挂牌
	FzmsupplyInvestAssetAction = 26; //摘牌
	FzmsupplyPayAction = 27; //支付
	FzmsupplyRecoverAction = 28; //解除冻结状态
	FzmsupplyCancelPayAction = 29; //撤销票据支付
	FzmsupplyBatchCashAction = 30; //批量兑付
	FzmsupplyOfflineCashAction = 31;//线下兑付
)

const (
	//log for fzmsupply


)

