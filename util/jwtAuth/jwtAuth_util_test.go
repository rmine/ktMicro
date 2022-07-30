package jwtAuth

import "testing"

//func TestCreateToken(t *testing.T) {
//	authInfo := &AuthInfo{UserId: 10001}
//	data,err := CreateToken(authInfo)
//	if err != nil {
//		t.Error("err",err)
//	} else {
//		t.Log("tokenstring",data)
//	}
//
//	//data = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoSW5mbyI6eyJVc2VySWQiOjEwMDAxfSwiZXhwIjoxNTkyNjMyOTYzLCJ1c2VyX2lkIjoxMDAwMX0.qvhYiJ90mKGdSmVUdl7WQmCu--b5w9cxbsqvGGn5RNQ"
//
//	token,err := ParseToken(data)
//	if err != nil {
//		t.Error("err2",err)
//	} else {
//		t.Log("token",token)
//	}
//
//	claims,err := VerifyToken(token)
//	if err != nil {
//		t.Error("not ok",err)
//	}
//
//	ret,err := GetValidAuthInfo(claims)
//	if err != nil {
//		t.Error("err3",err)
//	} else {
//		t.Log("ret",ret)
//	}
//}

func TestExpireTokenExp(t *testing.T) {
	data := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoSW5mbyI6eyJVc2VySWQiOjEwMDAxLCJVc2VySW5mbyI6bnVsbH0sImV4cCI6MTU5MjczMDY2NSwidXNlciI6bnVsbCwidXNlcl9pZCI6MTAwMDF9.BRxn_JxR0PkRdpztPnwaN75-O9YS_i8klV9xNJTd1iA"

	token, err := ParseToken(data)
	if err != nil {
		t.Error("err2", err)
	} else {
		t.Log("token", token)
	}

	claims, err := VerifyToken(token)
	if err != nil {
		t.Error("not ok", err)
	} else {
		t.Log("ok ?", claims)
	}

	//ExpireTokenExp(token)

	claims2, err := VerifyToken(token)
	if err != nil {
		t.Error("2== not ok", err)
	} else {
		t.Log("2=== ok ?", claims2)
	}
}
