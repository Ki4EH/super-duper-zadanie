package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725744678-team-77391/zadanie-6105/service/internal/db/models"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
)

type BidRepository interface {
	CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
	CheckAuthorID(ctx context.Context, bid *models.Bid) (bool, error)
	ListUserBids(ctx context.Context, filter BidFilter) ([]*models.Bid, error)
	ListOrganizationBids(ctx context.Context, filter BidFilter) ([]*models.Bid, error)
	UpdateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
	GetBidByUUID(ctx context.Context, id uuid.UUID) (*models.Bid, error)
	GetBidVersion(ctx context.Context, bidID uuid.UUID, version int) (*models.Bid, error)
	AddBidDecision(ctx context.Context, bidID uuid.UUID, userID uuid.UUID, decision string) error
	GetBidDecisions(ctx context.Context, bidID uuid.UUID) ([]models.BidDecision, error)
	UpdateBidStatus(ctx context.Context, bidID uuid.UUID, status string) error
	CloseTender(ctx context.Context, tenderID uuid.UUID) error
	GetResponsibleCount(ctx context.Context, organizationID uuid.UUID) (int, error)
	AddBidReview(ctx context.Context, review *models.BidReview) error
	GetBidReviewsByAuthor(ctx context.Context, tenderID, authorID uuid.UUID, limit, offset int) ([]models.BidReview, error)
	CheckTenderResponsibility(ctx context.Context, tenderID, userID uuid.UUID) (bool, error)
}

type bidRepository struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

type BidFilter struct {
	TenderId  uuid.UUID
	CreatorID uuid.UUID
	Limit     int
	Offset    int
}

func NewBidRepository(db *gorm.DB) BidRepository {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Ошибка инициализации SQL DB: %v", err))
	}
	return &bidRepository{db: db, sqlDB: sqlDB}
}

func (r *bidRepository) GetBidByUUID(ctx context.Context, id uuid.UUID) (*models.Bid, error) {
	//query, err := readSQLFile("internal/db/repository/sql/get_bid_by_id.sql")
	//if err != nil {
	//	return nil, err
	//}

	query := `SELECT id, name, description, status, version, tender_id, organization_id, author_id, author_type, created_at, updated_at
FROM bids
WHERE id = $1;`

	bid := &models.Bid{}
	row := r.sqlDB.QueryRowContext(ctx, query, id)
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.TenderID, &bid.OrganizationID, &bid.AuthorID, &bid.AuthorType, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса создания предложения: %v", err)
	}
	return bid, nil
}

func (r *bidRepository) SaveBidVersion(ctx context.Context, bid *models.Bid) error {
	//query, err := readSQLFile("internal/db/repository/sql/save_bid_version.sql")
	//if err != nil {
	//	return err
	//}

	query := `insert into bid_versions (bid_id, name, description, status, tender_id, organization_id, author_id, author_type, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), $10)
ON CONFLICT (bid_id, version) DO NOTHING;`

	_, err := r.sqlDB.ExecContext(ctx, query, bid.ID, bid.Name, bid.Description, bid.Status, bid.TenderID, bid.OrganizationID, bid.AuthorID, bid.AuthorType, bid.CreatedAt, bid.Version)
	if err != nil {
		return fmt.Errorf("ошибка сохранения версии Предложения: %v", err)
	}

	return nil
}

func (r *bidRepository) CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {
	//query, err := readSQLFile("internal/db/repository/sql/create_bid.sql")
	//if err != nil {
	//	return nil, err
	//}

	query := `INSERT INTO bids (name, description, status, version, tender_id, organization_id, author_id, author_type, created_at, updated_at)
VALUES ($1, $2, 'Created', 1, $3, $4, $5, $6, NOW(), NOW())
RETURNING id, name, description, status, version, tender_id, organization_id, author_id, author_type, created_at, updated_at;`

	row := r.sqlDB.QueryRowContext(ctx, query, bid.Name, bid.Description, bid.TenderID, bid.OrganizationID, bid.AuthorID, bid.AuthorType)
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.TenderID, &bid.OrganizationID, &bid.AuthorID, &bid.AuthorType, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса создания предложения: %v", err)
	}

	err = r.SaveBidVersion(ctx, bid)
	if err != nil {
		log.Printf("ошибка при создании bid и последующем сохранении: %v", err)
		return nil, err
	}

	return bid, nil
}

func (r *bidRepository) CheckAuthorID(ctx context.Context, bid *models.Bid) (bool, error) {

	query := `SELECT 1 from organization_responsible where user_id = $1;`

	var exists int

	err := r.sqlDB.QueryRowContext(ctx, query, bid.AuthorID).Scan(&exists)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil

	} else if err != nil {
		return false, fmt.Errorf("ошибка выполнения запроса на поиск authorId: %v", err)
	}

	return true, nil

}

func (r *bidRepository) UpdateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {

	query := `UPDATE bids
SET
    name = COALESCE(NULLIF($1, ''), name),
    description = COALESCE(NULLIF($2, ''), description),
    status = COALESCE(NULLIF($3, ''), status),
    version = version + 1,
    updated_at = NOW()
WHERE id = $4
    RETURNING id, name, description, status, version, tender_id, organization_id, author_id, author_type, created_at, updated_at
`

	row := r.sqlDB.QueryRowContext(ctx, query, bid.Name, bid.Description, bid.Status, bid.ID)
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.Version, &bid.TenderID, &bid.OrganizationID, &bid.AuthorID, &bid.AuthorType, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса обновления bid: %v", err)
	}

	err = r.SaveBidVersion(ctx, bid)
	if err != nil {
		log.Printf("ошибка при обновлении bid и последующем сохранении: %v", err)
		return nil, err
	}

	return bid, nil
}

func (r *bidRepository) ListUserBids(ctx context.Context, filter BidFilter) ([]*models.Bid, error) {
	var bids []*models.Bid
	query := r.db.WithContext(ctx).Model(&models.Bid{})

	if filter.CreatorID != uuid.Nil {
		query = query.Where("author_id = ?", filter.CreatorID)
	}

	query = query.Order("bids.name ASC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Select("bids.id, bids.name, bids.description, bids.status, bids.version, bids.tender_id, bids.organization_id, bids.author_id, bids.author_type, bids.created_at, bids.updated_at")

	if err := query.Find(&bids).Error; err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *bidRepository) ListOrganizationBids(ctx context.Context, filter BidFilter) ([]*models.Bid, error) {
	var bids []*models.Bid
	query := r.db.WithContext(ctx).Model(&models.Bid{})

	query = query.Joins(`
		LEFT JOIN organization_responsible ON bids.organization_id = organization_responsible.organization_id
	`).
		Where("organization_responsible.user_id = ?", filter.CreatorID)

	if filter.TenderId != uuid.Nil {
		query = query.Where("tender_id = ?", filter.TenderId)
	}

	query = query.Order("bids.name ASC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Select("bids.id, bids.name, bids.description, bids.status, bids.version, bids.tender_id, bids.organization_id, bids.author_id, bids.author_type, bids.created_at, bids.updated_at")

	if err := query.Find(&bids).Error; err != nil {
		return nil, err
	}

	return bids, nil
}

func (r *bidRepository) GetBidVersion(ctx context.Context, bidID uuid.UUID, version int) (*models.Bid, error) {
	query := `SELECT bid_id, name, description, status, tender_id, organization_id, author_id, author_type, created_at, updated_at, version
from bid_versions
where bid_id = $1 AND version = $2;`

	row := r.sqlDB.QueryRowContext(ctx, query, bidID, version)
	var bid models.Bid
	err := row.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.OrganizationID, &bid.AuthorID, &bid.AuthorType, &bid.CreatedAt, &bid.UpdatedAt, &bid.Version)
	if err != nil {
		return nil, fmt.Errorf("версия bid не найдена: %v", err)
	}

	return &bid, nil
}

func (r *bidRepository) AddBidDecision(ctx context.Context, bidID uuid.UUID, userID uuid.UUID, decision string) error {
	decisionEntry := models.BidDecision{
		BidID:    bidID,
		UserID:   userID,
		Decision: decision,
	}
	return r.db.WithContext(ctx).Select("bid_id", "user_id", "decision").Create(&decisionEntry).Error
}

func (r *bidRepository) GetBidDecisions(ctx context.Context, bidID uuid.UUID) ([]models.BidDecision, error) {
	var decisions []models.BidDecision
	err := r.db.WithContext(ctx).Where("bid_id = ?", bidID).Find(&decisions).Error
	return decisions, err
}

func (r *bidRepository) UpdateBidStatus(ctx context.Context, bidID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&models.Bid{}).Where("id = ?", bidID).Update("status", status).Error
}

func (r *bidRepository) CloseTender(ctx context.Context, tenderID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Tender{}).Where("id = ?", tenderID).Update("status", "Closed").Error
}

func (r *bidRepository) GetResponsibleCount(ctx context.Context, organizationID uuid.UUID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.OrganizationResponsible{}).Where("organization_id = ?", organizationID).Count(&count).Error
	return int(count), err
}

func (r *bidRepository) AddBidReview(ctx context.Context, review *models.BidReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *bidRepository) GetBidReviewsByAuthor(ctx context.Context, tenderID, authorID uuid.UUID, limit, offset int) ([]models.BidReview, error) {
	var reviews []models.BidReview

	err := r.db.WithContext(ctx).
		Model(&models.BidReview{}).
		Where("bid_id IN (SELECT id FROM bids WHERE tender_id = ?) AND author_id = ?", tenderID, authorID).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, err
}

func (r *bidRepository) CheckTenderResponsibility(ctx context.Context, tenderID, userID uuid.UUID) (bool, error) {
	var cnt []models.OrganizationResponsible
	query := r.db.WithContext(ctx).
		Model(&models.OrganizationResponsible{}).
		Where("organization_id = (SELECT organization_id FROM tenders WHERE id = ?) AND user_id = ?", tenderID, userID)
	if err := query.Find(&cnt).Error; err != nil {
		return false, err
	}

	return len(cnt) > 0, nil
}
