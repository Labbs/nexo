package persistence

import (
	"github.com/labbs/nexo/domain"
	"gorm.io/gorm"
)

type databasePers struct {
	db *gorm.DB
}

func NewDatabasePers(db *gorm.DB) *databasePers {
	return &databasePers{db: db}
}

func (p *databasePers) Create(database *domain.Database) error {
	return p.db.Create(database).Error
}

func (p *databasePers) GetById(id string) (*domain.Database, error) {
	var database domain.Database
	err := p.db.Debug().
		Preload("Space").
		Preload("User").
		Where("id = ?", id).
		First(&database).Error
	if err != nil {
		return nil, err
	}
	return &database, nil
}

func (p *databasePers) GetBySpaceId(spaceId string) ([]domain.Database, error) {
	var databases []domain.Database
	err := p.db.Debug().
		Preload("User").
		Where("space_id = ?", spaceId).
		Order("created_at DESC").
		Find(&databases).Error
	if err != nil {
		return nil, err
	}
	return databases, nil
}

func (p *databasePers) GetByDocumentId(documentId string) ([]domain.Database, error) {
	var databases []domain.Database
	err := p.db.Debug().
		Preload("User").
		Where("document_id = ?", documentId).
		Order("created_at DESC").
		Find(&databases).Error
	if err != nil {
		return nil, err
	}
	return databases, nil
}

func (p *databasePers) Update(database *domain.Database) error {
	return p.db.Debug().Save(database).Error
}

func (p *databasePers) Delete(id string) error {
	return p.db.Debug().Where("id = ?", id).Delete(&domain.Database{}).Error
}

// DatabaseRow persistence
type databaseRowPers struct {
	db *gorm.DB
}

func NewDatabaseRowPers(db *gorm.DB) *databaseRowPers {
	return &databaseRowPers{db: db}
}

func (p *databaseRowPers) Create(row *domain.DatabaseRow) error {
	return p.db.Debug().Create(row).Error
}

func (p *databaseRowPers) GetById(id string) (*domain.DatabaseRow, error) {
	var row domain.DatabaseRow
	err := p.db.Debug().
		Preload("Database").
		Preload("CreatedUser").
		Preload("UpdatedUser").
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (p *databaseRowPers) GetByDatabaseId(databaseId string, limit, offset int) ([]domain.DatabaseRow, error) {
	var rows []domain.DatabaseRow
	query := p.db.Debug().
		Preload("CreatedUser").
		Preload("UpdatedUser").
		Where("database_id = ?", databaseId).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *databaseRowPers) GetByDatabaseIdWithOptions(databaseId string, options domain.RowQueryOptions) ([]domain.DatabaseRow, error) {
	var rows []domain.DatabaseRow
	query := p.db.Debug().
		Preload("CreatedUser").
		Preload("UpdatedUser").
		Where("database_id = ?", databaseId)

	// Apply filters
	query = p.applyFilters(query, options.Filter)

	// Apply sorting
	if len(options.Sort) > 0 {
		for _, sort := range options.Sort {
			direction := "ASC"
			if sort.Direction == "desc" {
				direction = "DESC"
			}
			// Sort by JSON property value (SQLite compatible)
			query = query.Order("json_extract(properties, '$." + sort.PropertyId + "') " + direction)
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	err := query.Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *databaseRowPers) applyFilters(query *gorm.DB, filter *domain.FilterConfig) *gorm.DB {
	if filter == nil {
		return query
	}

	// Apply AND filters
	for _, rule := range filter.And {
		query = p.applyFilterRule(query, rule, "AND")
	}

	// Apply OR filters
	if len(filter.Or) > 0 {
		orQuery := p.db.Where("1 = 0") // Start with false for OR
		for _, rule := range filter.Or {
			orQuery = p.applyFilterRule(orQuery, rule, "OR")
		}
		query = query.Where(orQuery)
	}

	return query
}

func (p *databaseRowPers) applyFilterRule(query *gorm.DB, rule domain.FilterRule, _ string) *gorm.DB {
	// Use json_extract for SQLite compatibility
	propertyPath := "json_extract(properties, '$." + rule.Property + "')"

	// Get string value safely
	strValue := ""
	if rule.Value != nil {
		if s, ok := rule.Value.(string); ok {
			strValue = s
		}
	}

	switch rule.Condition {
	case "eq":
		if rule.Value == nil || strValue == "" {
			return query // Skip empty eq filters
		}
		return query.Where(propertyPath+" = ?", rule.Value)
	case "neq":
		if rule.Value == nil || strValue == "" {
			return query // Skip empty neq filters
		}
		return query.Where(propertyPath+" != ? OR "+propertyPath+" IS NULL", rule.Value)
	case "gt":
		return query.Where("CAST("+propertyPath+" AS REAL) > ?", rule.Value)
	case "lt":
		return query.Where("CAST("+propertyPath+" AS REAL) < ?", rule.Value)
	case "gte":
		return query.Where("CAST("+propertyPath+" AS REAL) >= ?", rule.Value)
	case "lte":
		return query.Where("CAST("+propertyPath+" AS REAL) <= ?", rule.Value)
	case "contains":
		if strValue == "" {
			return query // Skip empty contains filters
		}
		return query.Where(propertyPath+" LIKE ? COLLATE NOCASE", "%"+strValue+"%")
	case "not_contains":
		if strValue == "" {
			return query // Skip empty not_contains filters
		}
		return query.Where(propertyPath+" NOT LIKE ? COLLATE NOCASE OR "+propertyPath+" IS NULL", "%"+strValue+"%")
	case "starts_with":
		if strValue == "" {
			return query // Skip empty starts_with filters
		}
		return query.Where(propertyPath+" LIKE ? COLLATE NOCASE", strValue+"%")
	case "ends_with":
		if strValue == "" {
			return query // Skip empty ends_with filters
		}
		return query.Where(propertyPath+" LIKE ? COLLATE NOCASE", "%"+strValue)
	case "is_empty":
		return query.Where(propertyPath+" IS NULL OR "+propertyPath+" = ''")
	case "is_not_empty":
		return query.Where(propertyPath+" IS NOT NULL AND "+propertyPath+" != ''")
	default:
		return query
	}
}

func (p *databaseRowPers) GetRowCount(databaseId string) (int64, error) {
	var count int64
	err := p.db.Debug().Model(&domain.DatabaseRow{}).
		Where("database_id = ? AND deleted_at IS NULL", databaseId).
		Count(&count).Error
	return count, err
}

func (p *databaseRowPers) GetRowCountWithFilter(databaseId string, filter *domain.FilterConfig) (int64, error) {
	var count int64
	query := p.db.Debug().Model(&domain.DatabaseRow{}).
		Where("database_id = ? AND deleted_at IS NULL", databaseId)

	query = p.applyFilters(query, filter)

	err := query.Count(&count).Error
	return count, err
}

func (p *databaseRowPers) Update(row *domain.DatabaseRow) error {
	return p.db.Debug().Save(row).Error
}

func (p *databaseRowPers) Delete(id string) error {
	return p.db.Debug().Where("id = ?", id).Delete(&domain.DatabaseRow{}).Error
}

func (p *databaseRowPers) BulkDelete(ids []string) error {
	return p.db.Debug().Where("id IN ?", ids).Delete(&domain.DatabaseRow{}).Error
}
