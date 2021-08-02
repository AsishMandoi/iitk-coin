package database

import (
	"fmt"
	"time"

	"github.com/AsishMandoi/iitk-coin/global"
	"github.com/gomodule/redigo/redis"
)

func SetJWTid(rollno int, jti string) error {
	conn := pool.Get()
	defer conn.Close()

	// The commands between MULTI and EXEC are made atomic
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HSET", "token:"+fmt.Sprint(rollno), "jwtId", jti)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", "token:"+fmt.Sprint(rollno), (global.TknExpTime*int(time.Minute))/int(time.Second))
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func GetJWTid(rollno int) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	jwtId, err := redis.String(conn.Do("HGET", "token:"+fmt.Sprint(rollno), "jwtId"))
	if err != nil {
		return "", err
	}
	return jwtId, nil
}

func SetPwdResetDetails(resetPwd global.PwdResetObj) error {
	conn := pool.Get()
	defer conn.Close()

	// The commands between MULTI and EXEC are made atomic
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HMSET", "resetPwd:"+fmt.Sprint(resetPwd.Rollno), "pwd", resetPwd.NewPwd, "otp", resetPwd.Otp)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", "resetPwd:"+fmt.Sprint(resetPwd.Rollno), 120)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func GetPwdResetDetails(rollno int) (global.PwdResetObj, string, error) {
	conn := pool.Get()
	defer conn.Close()

	resetPwd := struct {
		Pwd string `redis:"pwd"`
		Otp string `redis:"otp"`
	}{}

	values, err := redis.Values(conn.Do("HGETALL", "resetPwd:"+fmt.Sprint(rollno)))
	if err != nil {
		return global.PwdResetObj{}, "password reset details not found", err
	} else if len(values) == 0 {
		return global.PwdResetObj{}, "password reset details not found", fmt.Errorf("Possible error - password reset request expired")
	}

	if err := redis.ScanStruct(values, &resetPwd); err != nil {
		return global.PwdResetObj{}, "Could not decode password reset details", err
	}

	return global.PwdResetObj{rollno, resetPwd.Pwd, resetPwd.Otp}, "", nil
}

func SetTfrDetails(tfr global.TxnObj) error {
	conn := pool.Get()
	defer conn.Close()

	// The commands between MULTI and EXEC are made atomic
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HMSET", "transferDetails:"+fmt.Sprint(tfr.Sender), "receiver", tfr.Receiver, "amtSent", tfr.AmtSent, "amtRcvd", tfr.AmtRcvd, "descr", tfr.Descr, "otp", tfr.Otp)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", "transferDetails:"+fmt.Sprint(tfr.Sender), 120)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func SetRdmDetails(rdm global.RedeemObj) error {
	conn := pool.Get()
	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("HMSET", "redeemDetails:"+fmt.Sprint(rdm.Redeemer), "itemId", rdm.ItemId, "amount", rdm.Amount, "descr", rdm.Descr, "otp", rdm.Otp)
	if err != nil {
		return err
	}
	err = conn.Send("EXPIRE", "redeemDetails:"+fmt.Sprint(rdm.Redeemer), 120)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func GetTfrDetails(sender int) (global.TxnObj, string, error) {
	tfrDet := struct {
		Receiver int     `redis:"receiver"`
		AmtSent  float64 `redis:"amtSent"`
		AmtRcvd  float64 `redis:"amtRcvd"`
		Descr    string  `redis:"descr"`
		Otp      string  `redis:"otp"`
	}{}

	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", "transferDetails:"+fmt.Sprint(sender)))
	if err != nil {
		return global.TxnObj{}, "transaction details not found", err
	} else if len(values) == 0 {
		return global.TxnObj{}, "transaction details not found", fmt.Errorf("Possible error - transfer request expired")
	}

	if err := redis.ScanStruct(values, &tfrDet); err != nil {
		return global.TxnObj{}, "Could not decode transaction details", err
	}

	return global.TxnObj{sender, tfrDet.Receiver, tfrDet.AmtSent, tfrDet.AmtRcvd, tfrDet.Descr, tfrDet.Otp}, "", nil
}

func GetRdmDetails(redeemer int) (global.RedeemObj, string, error) {
	rdmDet := struct {
		ItemId int     `redis:"itemId"`
		Amount float64 `redis:"amount"`
		Descr  string  `redis:"descr"`
		Otp    string  `redis:"otp"`
	}{}

	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", "redeemDetails:"+fmt.Sprint(redeemer)))
	if err != nil {
		return global.RedeemObj{}, "redeem details not found", err
	} else if len(values) == 0 {
		return global.RedeemObj{}, "redeem details not found", fmt.Errorf("Possible error - redeem request expired")
	}

	if err := redis.ScanStruct(values, &rdmDet); err != nil {
		return global.RedeemObj{}, "Could not decode redeem details", err
	}

	return global.RedeemObj{redeemer, rdmDet.ItemId, rdmDet.Amount, rdmDet.Descr, rdmDet.Otp}, "", nil
}

func DelTfrDetails(sender int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", "transferDetails:"+fmt.Sprint(sender), "receiver", "amtSent", "amtRcvd", "descr", "otp")
	if err != nil {
		return err
	}
	return nil
}

func DelRdmDetails(redeemer int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", "redeemDetails:"+fmt.Sprint(redeemer), "itemId", "amount", "descr", "otp")
	if err != nil {
		return err
	}
	return nil
}

func DelPwdResetDetails(rollno int) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", "resetPwd:"+fmt.Sprint(rollno), "pwd", "otp")
	if err != nil {
		return err
	}
	return nil
}
