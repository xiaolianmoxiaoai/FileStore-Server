package db

import(
	"database/sql"
	mydb "filestore-server/db/mysql"
	"fmt"
)

// 文件上传完成保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`," +
			" `file_size`, `file_addr`, `status`)values(?,?,?,?,1)")
	if err != nil{
		fmt.Println("Failed to prepare statement. err:"+err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil==err{
		if rf<=0{
			fmt.Println("File with hash: %s has been upload success!", filehash)

		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
	UploadAt sql.NullTime
}
//根据文件hash查询文件信息
func GetFileMetaByFilehash(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr, create_at " +
			"from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr, &tfile.UploadAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, err



}