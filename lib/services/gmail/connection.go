package gmail

import (
	"binadesa2020-backend/lib/models"
	"binadesa2020-backend/lib/variable"
	"fmt"
	"log"
	"net/smtp"
	"reflect"
)

// Email structure
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Send Email
func (e *Email) Send() error {
	config := variable.GmailConfig
	e.From = config.Email

	msg := "From: " + e.From + "\n" +
		"To: " + e.To + "\n" +
		"Subject: " + e.Subject + "\n\n" +
		e.Body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", e.From, config.Password, "smtp.gmail.com"),
		e.From,
		[]string{e.To},
		[]byte(msg),
	)
	if err != nil {
		log.Printf("Smtp error to %s : %v \n", e.To, err)
		return err
	}
	log.Printf("Smtp SUCCESS send mail to %s\n", e.To)
	return nil
}

// SendReceiveSubmission via email
// parameter is pointer from submission struct
func (e *Email) SendReceiveSubmission(ptrSubmission interface{}) error {
	var (
		typeSubmission string
		uniqueCode     string
	)

	switch v := ptrSubmission.(type) {
	case *models.AktaKelahiran:
		typeSubmission = "Akta Kelahiran"
		uniqueCode = v.UniqueCode
	case *models.KartuKeluarga:
		typeSubmission = "Kartu Keluarga"
		uniqueCode = v.UniqueCode
	case *models.SuratKeterangan:
		typeSubmission = fmt.Sprintf("Surat Keterangan %s", v.Tipe)
		uniqueCode = v.UniqueCode
	default:
		log.Printf("Invalid type in SendReceiveSubmission: %s", reflect.TypeOf(v))
	}

	e.Subject = "Pemberitahuan Kantor Desa Telukjambe"
	e.Body = fmt.Sprintf("Terima kasih telah mengajukan %s, melalui telukjambe.id! Kode unik anda adalah \n\n%s\n\n, yang digunakan untuk mengecek status pengajuan anda. Sehingga, kami sarankan untuk mengecek status anda secara berkala pada laman telukjambe.id \n Mohon tidak membalas pesan ini karena bersifat otomatis. Gunakan kolom pengaduan pada telukjambe.id jika membutuhkan informasi lebih lanjut.", typeSubmission, uniqueCode)

	return e.Send()
}

// SendCompleteSubmission via email
// parameter is pointer from submission struct
func (e *Email) SendCompleteSubmission(ptrSubmission interface{}) error {
	var (
		name           string
		uniqueCode     string
		typeSubmission string
		isPaid         bool = true
	)

	switch v := ptrSubmission.(type) {
	case *models.AktaKelahiran:
		typeSubmission = "Akta Kelahiran"
		name = v.Nama
		uniqueCode = v.UniqueCode
	case *models.KartuKeluarga:
		typeSubmission = "Kartu Keluarga"
		name = v.Nama
		uniqueCode = v.UniqueCode
	case *models.SuratKeterangan:
		typeSubmission = fmt.Sprintf("Surat Keterangan %s", v.Tipe)
		name = v.Nama
		uniqueCode = v.UniqueCode
		isPaid = v.IsPaid
	default:
		log.Printf("Invalid type in SendCompleteSubmission: %s", reflect.TypeOf(ptrSubmission))
	}

	e.Subject = "Pengajuan Berkas Anda Telah Selesai"
	e.Body = fmt.Sprintf("Pemberitahuan kepada %s dengan kode %s, bahwa pengajuan %s anda kini telah SELESAI! Silahkan mengunjungi kantor Desa Telukjambe untuk pengambilannya. ", name, uniqueCode, typeSubmission)

	if isPaid {
		e.Body = e.Body + fmt.Sprintf("Layanan pengajuan anda dikenakan biaya sebesar 20.000 rupiah.")
	} else {
		e.Body = e.Body + fmt.Sprintf("Layanan pengajuan anda dikenakan biaya sebesar 0 rupiah.")
	}

	return e.Send()
}

// SendRejectSubmission via email
// parameter is pointer from submission struct
func (e *Email) SendRejectSubmission(ptrSubmission interface{}) error {
	var (
		name       string
		uniqueCode string
	)

	switch v := ptrSubmission.(type) {
	case *models.AktaKelahiran:
		name = v.Nama
		uniqueCode = v.UniqueCode
	case *models.KartuKeluarga:
		name = v.Nama
		uniqueCode = v.UniqueCode
	case *models.SuratKeterangan:
		name = v.Nama
		uniqueCode = v.UniqueCode
	default:
		log.Printf("Invalid type in SendRejectSubmission: %s", reflect.TypeOf(ptrSubmission))
	}

	e.Subject = "Penolakan Berkas Pengajuan Desa Telukjambe"
	e.Body = fmt.Sprintf("Kepada %s dengan kode %s, mohon maaf pengajuan berkas anda DITOLAK. Mohon mendatangi kantor Desa Telukjambe untuk keterangan lebih lanjut. Terima kasih", name, uniqueCode)

	return e.Send()
}
