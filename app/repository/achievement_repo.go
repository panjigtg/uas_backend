package repository

import (
    "context"
    "uas/app/models"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type AchievementMongoRepository interface {
	Create(ctx context.Context, data *models.AchievementMongo) (string, error)
	FindByID(ctx context.Context, id string) (*models.AchievementMongo, error)
	SoftDelete(ctx context.Context, id string) error
    Update(ctx context.Context, a *models.AchievementMongo) error
}

type achievementMongoRepository struct {
    col *mongo.Collection
}

func NewAchievementMongoRepository(col *mongo.Collection) AchievementMongoRepository {
    return &achievementMongoRepository{col: col}
}

func (r *achievementMongoRepository) Create(ctx context.Context, a *models.AchievementMongo) (string, error) {
    res, err := r.col.InsertOne(ctx, a)
    if err != nil {
        return "", err
    }

    return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *achievementMongoRepository) FindByID(ctx context.Context, id string) (*models.AchievementMongo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

    opts := options.FindOne().SetProjection(bson.M{
        "description": 0,
        "details":     0,
        "tags":        0,
        "points":      0,
    })
	var out models.AchievementMongo
	err = r.col.FindOne(
        ctx,
        bson.M{
            "_id": objID,
            "isDeleted": bson.M{"$ne": true},
        },
        opts,
    ).Decode(&out)
    
	if err != nil {
		return nil, err
	}
	return &out, nil
}


func (r *achievementMongoRepository) SoftDelete(ctx context.Context, id string) error {
    oid, _ := primitive.ObjectIDFromHex(id)

    update := bson.M{
        "$set": bson.M{
            "isDeleted": true,
            "deletedAt": primitive.NewDateTimeFromTime(time.Now()),
        },
    }

    _, err := r.col.UpdateByID(ctx, oid, update)
    return err
}

func (r *achievementMongoRepository) Update(ctx context.Context, a *models.AchievementMongo) error {

    update := bson.M{
        "$set": bson.M{
            "student_id":       a.StudentID,
            "achievement_type": a.AchievementType,
            "title":            a.Title,
            "description":      a.Description,
            "details":          a.Details,
            "tags":             a.Tags,
            "attachments":      a.Attachments,
            "points":           a.Points,
            "updated_at":       time.Now(),
        },
    }

     _, err := r.col.UpdateByID(ctx, a.ID, update)
    return err
}
