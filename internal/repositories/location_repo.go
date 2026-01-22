type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) CreateLocation(location *models.Location) error {
	return r.db.Create(location).Error
}

func (r *LocationRepository) GetLocationByID(id string) (*models.Location, error) {
	var location models.Location
	if err := r.db.First(&location, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *LocationRepository) UpdateLocation(location *models.Location) error {
	return r.db.Save(location).Error
}

func (r *LocationRepository) DeleteLocation(id string) error {
	return r.db.Delete(&models.Location{}, "id = ?", id).Error
}