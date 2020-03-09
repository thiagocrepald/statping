package notifications

import (
	"errors"
	"github.com/statping/statping/database"
)

func DB() database.Database {
	return database.DB().Model(&Notification{})
}

func Append(n Notifier) {
	allNotifiers = append(allNotifiers, n)
}

func Find(name string) (Notifier, error) {
	for _, n := range allNotifiers {
		notif := n.Select()
		if notif.Name() == name || notif.Method == name {
			return n, nil
		}
	}
	return nil, errors.New("notifier not found")
}

func All() []Notifier {
	return allNotifiers
}

func (n *Notification) Create() error {
	var notif Notification
	if DB().Where("method = ?", n.Method).Find(&notif).RecordNotFound() {
		return DB().Create(n).Error()
	}
	return nil
}

func (n *Notification) Update() error {
	n.ResetQueue()
	if n.Enabled.Bool {
		n.Close()
		n.Start()
		go Queue(Notifier(n))
	} else {
		n.Close()
	}
	err := DB().Update(n)
	return err.Error()
}

func (n *Notification) Delete() error {
	return nil
}
