package withdrawal



// @summary List Banks
// @description Fetch a list of banks
// @tags withdrawal
// @accept json
// @produce json
// @success 200 {object} map[string]interface{} "List of banks"
// @failure 500 {object} map[string]interface{} "Internal server error"
// @Router /withdrawal/list-banks [get]
// @security BearerAuth
func _() {}


// @summary Retrieve Bank Details
// @description Retrieve details of a specific bank by its code
// @tags withdrawal
// @accept json
// @produce json
// @Param request body RetrieveAccountDetailsRequest true "Account Details Request"
// @success 200 {object} map[string]interface{} "Bank details"	
// @failure 400 {object} map[string]interface{} "Bad request"
// @failure 500 {object} map[string]interface{} "Internal server error"
// @Router /withdrawal/account-details [get]
// @security BearerAuth
func _() {}


// @summary Initiate Withdrawal
// @description Initiate a withdrawal request
// @tags withdrawal
// @accept json
// @produce json
// @Param request body InitiateWithdrawalRequest true "Withdrawal Request"
// @success 200 {object} map[string]interface{} "Withdrawal initiated successfully"
// @failure 400 {object} map[string]interface{} "Bad request"
// @failure 500 {object} map[string]interface{} "Internal server error"
// @Router /withdrawal/initiate-withdrawal [post]
// @security BearerAuth
func _() {}