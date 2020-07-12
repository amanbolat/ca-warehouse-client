package crm

import "time"

type Note struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type FileMakerNote struct {
	ID          string    `json:"Id_note"`
	DateCreated time.Time `json:"Date_Created_Timestamp"`
	Content     string    `json:"NoteContent"`
	RelatedId   string    `json:"Id_relationID"`
}

func (fn FileMakerNote) ToNote() *Note {
	return &Note{
		ID:        fn.ID,
		CreatedAt: fn.DateCreated,
		Content:   fn.Content,
	}
}
