package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongo struct {
    ID          	primitive.ObjectID 		`bson:"_id,omitempty" json:"id"`
    StudentID   	string             		`bson:"studentId" json:"student_id"`
    AchievementType string         			`bson:"achievementType" json:"achievement_type"` 
    Title           string         			`bson:"title" json:"title"`
    Description     string         			`bson:"description" json:"description"`
    Details 		AchievementDetails     	`bson:"details" json:"details"`
    Attachments 	[]AchievementFile  		`bson:"attachments" json:"attachments"`
    Tags   			[]string  				`bson:"tags" json:"tags"`
    Points 			int       				`bson:"points" json:"points"`
    CreatedAt 		time.Time            	`bson:"createdAt" json:"created_at"`
    UpdatedAt		time.Time            	`bson:"updatedAt" json:"updated_at"`
}

type AchievementDetails struct {
    // competition
    CompetitionName  		*string   			`bson:"competitionName,omitempty" json:"competition_name,omitempty"`
    CompetitionLevel 		*string   			`bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
    Rank             		*int      			`bson:"rank,omitempty" json:"rank,omitempty"`
    MedalType        		*string   			`bson:"medalType,omitempty" json:"medal_type,omitempty"`

    // publication	
    PublicationType  		*string   			`bson:"publicationType,omitempty" json:"publication_type,omitempty"`
    PublicationTitle 		*string   			`bson:"publicationTitle,omitempty" json:"publication_title,omitempty"`
    Authors          		[]string  			`bson:"authors,omitempty" json:"authors,omitempty"`
    Publisher        		*string   			`bson:"publisher,omitempty" json:"publisher,omitempty"`
    ISSN             		*string   			`bson:"issn,omitempty" json:"issn,omitempty"`

    // organization
    OrganizationName 		*string   			`bson:"organizationName,omitempty" json:"organization_name,omitempty"`
    Position         		*string   			`bson:"position,omitempty" json:"position,omitempty"`
    Period           		*Period   			`bson:"period,omitempty" json:"period,omitempty"`

    // certification
    CertificationName   	*string    			`bson:"certificationName,omitempty" json:"certification_name,omitempty"`
    IssuedBy            	*string    			`bson:"issuedBy,omitempty" json:"issued_by,omitempty"`
    CertificationNumber 	*string    			`bson:"certificationNumber,omitempty" json:"certification_number,omitempty"`
    ValidUntil          	*time.Time 			`bson:"validUntil,omitempty" json:"valid_until,omitempty"`

    // general
    EventDate  				*time.Time  		`bson:"eventDate,omitempty" json:"event_date,omitempty"`
    Location   				*string     		`bson:"location,omitempty" json:"location,omitempty"`
    Organizer  				*string     		`bson:"organizer,omitempty" json:"organizer,omitempty"`
    Score      				*int        		`bson:"score,omitempty" json:"score,omitempty"`

    // dynamic fields
    CustomFields 			map[string]any      `bson:"customFields,omitempty" json:"custom_fields,omitempty"`
}

type Period struct {
    Start time.Time `bson:"start" json:"start"`
    End   time.Time `bson:"end" json:"end"`
}

type AchievementFile struct {
    FileName   string    `bson:"fileName" json:"file_name"`
    FileURL    string    `bson:"fileUrl" json:"file_url"`
    FileType   string    `bson:"fileType" json:"file_type"`
    UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}
