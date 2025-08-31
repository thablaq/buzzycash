package transaction

import (
	"net/http"
	"log"
	"strconv"
	"strings"
	"fmt"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/gin-gonic/gin"
)



func GetTransactionHistoryHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	log.Printf("INFO: GetTransactionHistoryHandler called for user ID: %d", currentUser.ID)

	// Pagination: page number from query, default is 1
	pageStr := ctx.Query("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			log.Printf("WARN: GetTransactionHistoryHandler: User ID: %d, Invalid page parameter '%s', defaulting to 1", currentUser.ID, pageStr)
		}
	}
	limit := 20
	offset := (page - 1) * limit
	log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Page: %d, Limit: %d, Offset: %d", currentUser.ID, page, limit, offset)

	// Collect filters dynamically
	appliedFilters := make(map[string]string)
	if status := ctx.Query("payment_status"); status != "" {
		appliedFilters["payment_status"] = status
	}
	if pType := ctx.Query("payment_type"); pType != "" {
		appliedFilters["payment_type"] = pType
	}
	if category := ctx.Query("category"); category != "" {
		appliedFilters["category"] = category
	}
	if tType := ctx.Query("transaction_type"); tType != "" {
		appliedFilters["transaction_type"] = tType
	}
	if tmethod := ctx.Query("payment_method"); tmethod != "" {
		appliedFilters["payment_method"] = tmethod
	}
	if tcurrency := ctx.Query("currency"); tcurrency != "" {
		appliedFilters["currency"] = tcurrency
	}
	log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Applied filters: %+v", currentUser.ID, appliedFilters)

	// Helper function to run query with given filters
	runQuery := func(filters map[string]string) ([]models.TransactionHistory, int64, error) {
		var histories []models.TransactionHistory
		var totalCount int64

		q := config.DB.Where("user_id = ?", currentUser.ID)
	 // ✅ Case-insensitive for text-based filters, case-sensitive for others
        for k, v := range filters {
        switch k {
        case "payment_status", "category", "payment_type", "transaction_type", "payment_method", "currency":
            q = q.Where(fmt.Sprintf("LOWER(%s) = ?", k), strings.ToLower(v))
        default:
            q = q.Where(fmt.Sprintf("%s = ?", k), v)
        }
         }
		log.Printf("DEBUG: GetTransactionHistoryHandler: User ID: %d, Running query with filters: %+v", currentUser.ID, filters)

		if err := q.Model(&models.TransactionHistory{}).Count(&totalCount).Error; err != nil {
			log.Printf("ERROR: GetTransactionHistoryHandler: User ID: %d, Failed to count transactions with filters %+v, error: %v", currentUser.ID, filters, err)
			return nil, 0, err
		}
		log.Printf("DEBUG: GetTransactionHistoryHandler: User ID: %d, Total count for filters %+v: %d", currentUser.ID, filters, totalCount)


		if err := q.Order("paid_at desc, created_at desc").
			Limit(limit).Offset(offset).Find(&histories).Error; err != nil {
			log.Printf("ERROR: GetTransactionHistoryHandler: User ID: %d, Failed to fetch transactions with filters %+v, error: %v", currentUser.ID, filters, err)
			return nil, 0, err
		}
		log.Printf("DEBUG: GetTransactionHistoryHandler: User ID: %d, Found %d transactions for filters %+v", currentUser.ID, len(histories), filters)

		return histories, totalCount, nil
	}

	// 1. Try with all filters
	log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Attempting to fetch transactions with all applied filters.", currentUser.ID)
	histories, totalCount, err := runQuery(appliedFilters)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch transaction history")
		return
	}
	if len(histories) > 0 {
		log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Found %d transactions with all applied filters.", currentUser.ID, len(histories))
	} else {
		log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, No transactions found with all applied filters.", currentUser.ID)
	}


	// 2. Progressive fallback: remove filters one by one until something is found
	if len(histories) == 0 && len(appliedFilters) > 0 {
		log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, No results with all filters, initiating progressive fallback.", currentUser.ID)
		keys := make([]string, 0, len(appliedFilters))
		for k := range appliedFilters {
			keys = append(keys, k)
		}

		// Try removing one filter at a time
		for i := 0; i < len(keys); i++ {
			tmp := make(map[string]string)
			removedKey := keys[i]
			for j, k := range keys {
				if j != i {
					tmp[k] = appliedFilters[k]
				}
			}
			log.Printf("DEBUG: GetTransactionHistoryHandler: User ID: %d, Fallback attempt: removing filter '%s'. New filters: %+v", currentUser.ID, removedKey, tmp)

			histories, totalCount, err = runQuery(tmp)
			if err != nil {
				utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch transaction history")
				return
			}
			if len(histories) > 0 {
				log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Found %d transactions after removing filter '%s'.", currentUser.ID, len(histories), removedKey)
				break
			}
			log.Printf("DEBUG: GetTransactionHistoryHandler: User ID: %d, Still no transactions after removing filter '%s'.", currentUser.ID, removedKey)
		}
	}

	// 3. Absolute last fallback: no filters at all
	if len(histories) == 0 {
		log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, No transactions found after progressive fallback. Attempting to fetch without any filters.", currentUser.ID)
		histories, totalCount, err = runQuery(map[string]string{})
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch transaction history")
			return
		}
		if len(histories) > 0 {
			log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Found %d transactions with no filters.", currentUser.ID, len(histories))
		} else {
			log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Still no transactions found even without filters.", currentUser.ID)
		}
	}

	// Map to response
	var response []TransactionHistoryResponse
	for _, h := range histories {
		response = append(response, TransactionHistoryResponse{
			ID:                   h.ID,
			TicketPurchaseID:     h.TicketPurchaseID,
			Amount:               float64(h.Amount),
			TransactionReference: h.TransactionReference,
			Reference:            h.Reference,
			CustomerEmail:        h.CustomerEmail,
			PaymentStatus:        string(h.PaymentStatus),
			PaymentType:          string(h.PaymentType),
			Currency:             string(h.Currency),
			PaidAt:               h.PaidAt,
			TransactionType:      string(h.TransactionType),
			PaymentMethod:        string(h.PaymentMethod),
			Category:             string(h.Category),
		})
	}
	log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Mapped %d transaction history records to response.", currentUser.ID, len(response))

	hasMore := int64(offset+limit) < totalCount
	log.Printf("INFO: GetTransactionHistoryHandler: User ID: %d, Final response: transactions count %d, total_count %d, page %d, has_more %t", currentUser.ID, len(response), totalCount, page, hasMore)

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": response,
		"page":         page,
		"has_more":     hasMore,
		"total_count":  totalCount,
	})
}



func SearchTransactionHistoryHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	log.Printf("INFO: SearchTransactionHistoryHandler called for user ID: %d", currentUser.ID)

	// Pagination
	pageStr := ctx.Query("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	limit := 20
	offset := (page - 1) * limit
	log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, Page: %d, Limit: %d", currentUser.ID, page, limit)

	// Search term
	search := ctx.Query("search")
	if search == "" {
		log.Printf("WARN: SearchTransactionHistoryHandler: User ID: %d, search query is empty", currentUser.ID)
		utils.Error(ctx, http.StatusBadRequest, "Search query cannot be empty")
		return
	}
	log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, Search term: '%s'", currentUser.ID, search)

	// Progressive fallback query
	var histories []models.TransactionHistory
	var totalCount int64

	// Helper to run a search query
	runSearch := func(term string) ([]models.TransactionHistory, int64, error) {
		var hs []models.TransactionHistory
		var count int64
		like := "%" + term + "%"
		log.Printf("DEBUG: SearchTransactionHistoryHandler: User ID: %d, Running search for term: '%s'", currentUser.ID, term)

		q := config.DB.Where("user_id = ?", currentUser.ID).
			Where(`
				transaction_reference ILIKE ? OR 
				reference ILIKE ? OR 
				customer_email ILIKE ? OR 
				category ILIKE ? OR
				payment_method ILIKE ? OR
				payment_type ILIKE ? OR
				payment_status ILIKE ? OR
				transaction_type ILIKE ? OR
				currency ILIKE ?`,
				like, like, like, like, like, like, like, like,like)

		if err := q.Model(&models.TransactionHistory{}).Count(&count).Error; err != nil {
			log.Printf("ERROR: SearchTransactionHistoryHandler: User ID: %d, Failed to count transactions for term '%s', error: %v", currentUser.ID, term, err.Error())
			return nil, 0, err
		}
		if err := q.Order("paid_at desc, created_at desc").
			Limit(limit).Offset(offset).Find(&hs).Error; err != nil {
			log.Printf("ERROR: SearchTransactionHistoryHandler: User ID: %d, Failed to find transactions for term '%s', error: %v", currentUser.ID, term, err)
			return nil, 0, err
		}
		log.Printf("DEBUG: SearchTransactionHistoryHandler: User ID: %d, Found %d transactions for term '%s'", currentUser.ID, len(hs), term)

		return hs, count, nil
	}

	// Run initial search
	histories, totalCount, err := runSearch(search)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to search transaction history")
		return
	}
	log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, Initial search for '%s' returned %d results", currentUser.ID, search, len(histories))


	// Smart fallback: progressively strip characters from search term
	if len(histories) == 0 && len(search) > 3 {
		log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, No results for initial search, attempting smart fallback", currentUser.ID)
		for i := len(search) - 1; i >= 3; i-- { // don’t go below 3 chars
			fallbackSearchTerm := search[:i]
			log.Printf("DEBUG: SearchTransactionHistoryHandler: User ID: %d, Fallback search with term: '%s'", currentUser.ID, fallbackSearchTerm)
			histories, totalCount, err = runSearch(fallbackSearchTerm)
			if err != nil {
				utils.Error(ctx, http.StatusInternalServerError, "Failed to search transaction history")
				return
			}
			if len(histories) > 0 {
				log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, Found %d results with fallback term '%s'", currentUser.ID, len(histories), fallbackSearchTerm)
				break
			}
		}
	}

	// Still nothing? return empty but valid response
	if len(histories) == 0 {
		log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, No transactions found after all search attempts for '%s'", currentUser.ID, search)
		ctx.JSON(http.StatusOK, gin.H{
			"transactions": []TransactionHistoryResponse{}, // Returns an empty array 
			"page":         page,
			"has_more":     false,
			"total_count":  0,
		})
		return
	}

	// Map response
	var response []TransactionHistoryResponse
	for _, h := range histories {
		response = append(response, TransactionHistoryResponse{
			ID:                   h.ID,
			TicketPurchaseID:     h.TicketPurchaseID,
			Amount:               float64(h.Amount),
			TransactionReference: h.TransactionReference,
			Reference:            h.Reference,
			CustomerEmail:        h.CustomerEmail,
			PaymentStatus:        string(h.PaymentStatus),
			PaymentType:          string(h.PaymentType),
			Currency:             string(h.Currency),
			PaidAt:               h.PaidAt,
			TransactionType:      string(h.TransactionType),
			PaymentMethod:        string(h.PaymentMethod),
			Category:             string(h.Category),
		})
	}

	hasMore := int64(offset+limit) < totalCount
	log.Printf("INFO: SearchTransactionHistoryHandler: User ID: %d, Returning %d transactions for search '%s', total_count: %d, has_more: %t", currentUser.ID, len(response), search, totalCount, hasMore)

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": response,
		"page":         page,
		"has_more":     hasMore,
		"total_count":  totalCount,
	})
}

func GetTransactionByID(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id")
	log.Printf("INFO: GetTransactionByID called for user ID: %d, transaction ID: %s", currentUser.ID, id)

	var tx models.TransactionHistory
	if err := config.DB.First(&tx, "id = ? AND user_id = ?", id, currentUser.ID).Error; err != nil {
		log.Printf("WARN: GetTransactionByID: transaction not found for user ID: %d, transaction ID: %s, error: %v", currentUser.ID, id, err)
		utils.Error(ctx, http.StatusNotFound, "transaction not found")
		return
	}

	log.Printf("INFO: GetTransactionByID: Successfully retrieved transaction ID: %d for user ID: %d", tx.ID, currentUser.ID)
	response := TransactionHistoryResponse{
		ID:                   tx.ID,
		TicketPurchaseID:     tx.TicketPurchaseID,
		Amount:               float64(tx.Amount),
		TransactionReference: tx.TransactionReference,
		Reference:            tx.Reference,
		CustomerEmail:        tx.CustomerEmail,
		PaymentStatus:        string(tx.PaymentStatus),
		PaymentType:          string(tx.PaymentType),
		Currency:             string(tx.Currency),
		PaidAt:               tx.PaidAt,
		TransactionType:      string(tx.TransactionType),
		PaymentMethod:        string(tx.PaymentMethod),
		Category:             string(tx.Category),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":     currentUser.ID,
		"transaction": response,
	})
}
