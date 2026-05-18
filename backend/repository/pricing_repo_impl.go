package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/azdonald/pharmd/backend/models"
)

type PricingRepoImpl struct {
	db *sql.DB
}

func NewPricingRepositoryImpl(db *sql.DB) PricingRepository {
	return &PricingRepoImpl{db: db}
}

func (r *PricingRepoImpl) ListPrices(ctx context.Context, productID, locationID string, page, limit int) ([]models.ProductPrice, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	where := "pp.organisation_id = ?"
	args := []interface{}{orgID}

	if productID != "" {
		where += " AND pp.product_id = ?"
		args = append(args, productID)
	}
	if locationID != "" {
		where += " AND pp.location_id = ?"
		args = append(args, locationID)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM product_prices pp WHERE %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT pp.id, pp.organisation_id, pp.product_id, p.name, pp.location_id, l.name,
		       pp.selling_price, pp.cost_price, pp.min_price, pp.max_discount,
		       pp.is_active, pp.created_at, pp.updated_at
		 FROM product_prices pp
		 JOIN products p ON p.id = pp.product_id
		 JOIN locations l ON l.id = pp.location_id
		 WHERE %s ORDER BY p.name ASC LIMIT ? OFFSET ?`, where),
		append(args, limit, offset)...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var prices []models.ProductPrice
	for rows.Next() {
		var p models.ProductPrice
		if err := rows.Scan(&p.ID, &p.OrganisationID, &p.ProductID, &p.ProductName, &p.LocationID, &p.LocationName,
			&p.SellingPrice, &p.CostPrice, &p.MinPrice, &p.MaxDiscount,
			&p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		prices = append(prices, p)
	}
	return prices, total, nil
}

func (r *PricingRepoImpl) GetPriceByID(ctx context.Context, id string) (*models.ProductPrice, error) {
	orgID := ctx.Value("organisation_id").(string)
	p := &models.ProductPrice{}
	err := r.db.QueryRowContext(ctx,
		`SELECT pp.id, pp.organisation_id, pp.product_id, pp.location_id,
		        pp.selling_price, pp.cost_price, pp.min_price, pp.max_discount,
		        pp.is_active, pp.created_at, pp.updated_at
		 FROM product_prices pp
		 WHERE pp.id = ? AND pp.organisation_id = ?`,
		id, orgID,
	).Scan(&p.ID, &p.OrganisationID, &p.ProductID, &p.LocationID,
		&p.SellingPrice, &p.CostPrice, &p.MinPrice, &p.MaxDiscount,
		&p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PricingRepoImpl) UpsertPrice(ctx context.Context, price models.ProductPrice) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO product_prices (id, organisation_id, product_id, location_id, selling_price, cost_price, min_price, max_discount, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE selling_price = VALUES(selling_price), cost_price = VALUES(cost_price),
		         min_price = VALUES(min_price), max_discount = VALUES(max_discount),
		         updated_at = VALUES(updated_at)`,
		price.ID, orgID, price.ProductID, price.LocationID,
		price.SellingPrice, price.CostPrice, price.MinPrice, price.MaxDiscount,
		price.IsActive, price.CreatedAt, price.UpdatedAt,
	)
	return err
}

func (r *PricingRepoImpl) DeletePrice(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx, "DELETE FROM product_prices WHERE id = ? AND organisation_id = ?", id, orgID)
	return err
}

func (r *PricingRepoImpl) ListDiscountRules(ctx context.Context, page, limit int) ([]models.DiscountRule, int, error) {
	orgID := ctx.Value("organisation_id").(string)
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM discount_rules WHERE organisation_id = ?", orgID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, organisation_id, name, type, value, min_order_value, max_discount_amount,
		        applies_to, COALESCE(applies_to_id,''), is_active,
		        COALESCE(start_date,''), COALESCE(end_date,''),
		        created_at, updated_at
		 FROM discount_rules WHERE organisation_id = ? ORDER BY name ASC LIMIT ? OFFSET ?`,
		orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []models.DiscountRule
	for rows.Next() {
		var r models.DiscountRule
		if err := rows.Scan(&r.ID, &r.OrganisationID, &r.Name, &r.Type, &r.Value,
			&r.MinOrderValue, &r.MaxDiscountAmount, &r.AppliesTo, &r.AppliesToID,
			&r.IsActive, &r.StartDate, &r.EndDate, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		rules = append(rules, r)
	}
	return rules, total, nil
}

func (r *PricingRepoImpl) GetDiscountRuleByID(ctx context.Context, id string) (*models.DiscountRule, error) {
	orgID := ctx.Value("organisation_id").(string)
	rule := &models.DiscountRule{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, organisation_id, name, type, value, min_order_value, max_discount_amount,
		        applies_to, COALESCE(applies_to_id,''), is_active,
		        COALESCE(start_date,''), COALESCE(end_date,''),
		        created_at, updated_at
		 FROM discount_rules WHERE id = ? AND organisation_id = ?`,
		id, orgID,
	).Scan(&rule.ID, &rule.OrganisationID, &rule.Name, &rule.Type, &rule.Value,
		&rule.MinOrderValue, &rule.MaxDiscountAmount, &rule.AppliesTo, &rule.AppliesToID,
		&rule.IsActive, &rule.StartDate, &rule.EndDate, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func (r *PricingRepoImpl) CreateDiscountRule(ctx context.Context, rule models.DiscountRule) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO discount_rules (id, organisation_id, name, type, value, min_order_value,
		                               max_discount_amount, applies_to, applies_to_id, is_active,
		                               start_date, end_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?,''), ?, NULLIF(?,''), NULLIF(?,''), ?, ?)`,
		rule.ID, orgID, rule.Name, rule.Type, rule.Value, rule.MinOrderValue,
		rule.MaxDiscountAmount, rule.AppliesTo, rule.AppliesToID, rule.IsActive,
		rule.StartDate, rule.EndDate, rule.CreatedAt, rule.UpdatedAt,
	)
	return err
}

func (r *PricingRepoImpl) UpdateDiscountRule(ctx context.Context, id string, rule models.DiscountRule) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx,
		`UPDATE discount_rules SET name=?, type=?, value=?, min_order_value=?, max_discount_amount=?,
		        applies_to=?, applies_to_id=NULLIF(?,''), is_active=?,
		        start_date=NULLIF(?,''), end_date=NULLIF(?,''), updated_at=?
		 WHERE id=? AND organisation_id=?`,
		rule.Name, rule.Type, rule.Value, rule.MinOrderValue, rule.MaxDiscountAmount,
		rule.AppliesTo, rule.AppliesToID, rule.IsActive,
		rule.StartDate, rule.EndDate, rule.UpdatedAt, id, orgID,
	)
	return err
}

func (r *PricingRepoImpl) DeleteDiscountRule(ctx context.Context, id string) error {
	orgID := ctx.Value("organisation_id").(string)
	_, err := r.db.ExecContext(ctx, "DELETE FROM discount_rules WHERE id = ? AND organisation_id = ?", id, orgID)
	return err
}

var _ PricingRepository = (*PricingRepoImpl)(nil)
