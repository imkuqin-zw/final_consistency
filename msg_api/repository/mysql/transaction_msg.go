package mysql

import (
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	"time"
)

func (r *RepoMysql) GetTransMsgById(id uint64) (*models.TransactionMsg, error) {
	m := &models.TransactionMsg{Id: id}
	if err := r.getById(db, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *RepoMysql) GetTransMsgByMsgId(msgId string) (*models.TransactionMsg, error) {
	m := &models.TransactionMsg{}
	err := db.Where("msg_id = ?", msgId).First(m).Error
	return m, err
}

func (r *RepoMysql) UpdateTransMsgStatusByMsgId(msg *models.TransactionMsg, status uint8) (int64, error) {
	data := map[string]interface{}{
		"status":    status,
		"update_at": time.Now(),
	}
	query := "msg_id = ? and status = ?"
	return r.updateTransMsg(data, query, msg.MsgId, msg.Status)
}

func (r *RepoMysql) updateTransMsg(data map[string]interface{}, query string, arg ...interface{}) (int64, error) {
	res := db.Where(query, arg).Update(data)
	return res.RowsAffected, res.Error
}

func (r *RepoMysql) InsertTransMsg(m *models.TransactionMsg) error {
	return r.insert(db, m)
}
