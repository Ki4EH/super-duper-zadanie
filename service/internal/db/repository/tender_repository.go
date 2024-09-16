package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Ki4EH/super-duper-zadanie/service/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type TenderRepository interface {
	CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error)
	UpdateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error)
	GetTenderByUUID(ctx context.Context, id uuid.UUID) (*models.Tender, error)
	ListTenders(ctx context.Context, filter TenderFilter) ([]*models.Tender, error)
	GetTenderVersion(ctx context.Context, tenderID uuid.UUID, version int) (*models.Tender, error)
}

// tenderRepository реализация интерфейса TenderRepository
type tenderRepository struct {
	db      *gorm.DB
	sqlDB   *sql.DB
	orgRepo OrganizationRepository
}

// TenderFilter структура для фильтрации тендеров
type TenderFilter struct {
	CreatorID    uuid.UUID
	ServiceTypes []string
	Limit        int
	Offset       int
}

func NewTenderRepository(db *gorm.DB) TenderRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Ошибка инициализации SQL DB: %v", err))
	}
	return &tenderRepository{db: db, sqlDB: sqlDB}
}

func (r *tenderRepository) SaveTenderVersion(ctx context.Context, tender *models.Tender) error {

	query := `INSERT INTO tender_versions (tender_id, version, name, description, status, organization_id, creator_id, created_at, service_type)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (tender_id, version) DO NOTHING;`
	_, err := r.sqlDB.ExecContext(ctx, query, tender.ID, tender.Version, tender.Name, tender.Description, tender.Status, tender.OrganizationID, tender.CreatorID, tender.CreatedAt, tender.ServiceType)
	if err != nil {
		return fmt.Errorf("ошибка сохранения версии тендера: %v", err)
	}

	return nil
}

// ListTenders возвращает список тендеров с возможностью фильтрации
func (r *tenderRepository) ListTenders(ctx context.Context, filter TenderFilter) ([]*models.Tender, error) {

	var tenders []*models.Tender
	query := r.db.WithContext(ctx).Model(&models.Tender{})

	if len(filter.ServiceTypes) > 0 {
		query = query.Where("service_type IN ?", filter.ServiceTypes)
	}

	if filter.CreatorID != uuid.Nil {
		query = query.Where("creator_id = ?", filter.CreatorID)
	}

	query = query.Order("name ASC")

	query = query.Limit(filter.Limit).Offset(filter.Offset)
	query = query.Select("id, name, description, status, service_type, version, created_at, creator_id, organization_id")

	if err := query.Find(&tenders).Error; err != nil {
		return nil, err
	}

	return tenders, nil
}

// CreateTender добавляет новый тендер в базу данных
func (r *tenderRepository) CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error) {

	query := `INSERT INTO tenders (name, description, status, version, organization_id, creator_id, created_at, updated_at, service_type)
VALUES ($1, $2, 'Created', 1, $3, $4, NOW(), NOW(), $5)
RETURNING id, name, description, status, version, created_at, service_type;`

	row := r.sqlDB.QueryRowContext(ctx, query, tender.Name, tender.Description, tender.OrganizationID, tender.CreatorID, tender.ServiceType)
	err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.Version, &tender.CreatedAt, &tender.ServiceType)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса создания тендера: %v", err)
	}

	err = r.SaveTenderVersion(ctx, tender)
	if err != nil {
		log.Printf("ошибка при создании tender и последующем сохранении: %v", err)
		return nil, err
	}

	return tender, nil
}

// UpdateTender обновляет существующий тендер в базе данных
func (r *tenderRepository) UpdateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error) {

	query := `UPDATE tenders
SET
    name = COALESCE(NULLIF($1, ''), name),
    description = COALESCE(NULLIF($2, ''), description),
    service_type = COALESCE(NULLIF($3, ''), service_type),
    status = COALESCE(NULLIF($4, ''), status),
    version = version + 1,
    updated_at = NOW()
WHERE id = $5
    RETURNING id, name, description, service_type, status, version, organization_id, creator_id, created_at, updated_at;
`

	row := r.sqlDB.QueryRowContext(ctx, query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.ID)
	err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.Version, &tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt, &tender.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса обновления тендера: %v", err)
	}

	err = r.SaveTenderVersion(ctx, tender)
	if err != nil {
		log.Printf("ошибка при обновлении tender и последующем сохранении: %v", err)
		return nil, err
	}

	return tender, nil
}

func (r *tenderRepository) GetTenderVersion(ctx context.Context, tenderID uuid.UUID, version int) (*models.Tender, error) {

	query := `SELECT tender_id, name, description, status, version, organization_id, creator_id, service_type
FROM tender_versions
WHERE tender_id = $1 AND version = $2;`

	row := r.sqlDB.QueryRowContext(ctx, query, tenderID, version)
	var versionedTender models.Tender
	err := row.Scan(&versionedTender.ID, &versionedTender.Name, &versionedTender.Description, &versionedTender.Status, &versionedTender.Version, &versionedTender.OrganizationID, &versionedTender.CreatorID, &versionedTender.ServiceType)
	if err != nil {
		return nil, fmt.Errorf("версия тендера не найдена: %v", err)
	}

	return &versionedTender, nil
}

func (r *tenderRepository) GetTenderByUUID(ctx context.Context, id uuid.UUID) (*models.Tender, error) {

	query := `SELECT id, name, description, status, version, organization_id, creator_id, created_at, updated_at
FROM tenders
WHERE id = $1;`

	tender := &models.Tender{}
	row := r.sqlDB.QueryRowContext(ctx, query, id)
	err := row.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.Version, &tender.OrganizationID, &tender.CreatorID, &tender.CreatedAt, &tender.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("тендер не найден: %v", err)
	}

	return tender, nil
}
