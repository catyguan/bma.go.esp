package filestoreservice

import (
	"esp/filestore/filestoreprot"
	"fmt"
	"io"
	"logger"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"
)

type request struct {
	path       string
	tempFile   string
	deleteFile bool
}

type session struct {
	sessionToken    filestoreprot.SessionToken
	appId           filestoreprot.AppId
	lastRequestTime time.Time
	requests        []*request
}

func (this *session) String() string {
	return fmt.Sprintf("%s:%s", this.appId, this.sessionToken)
}

// impl
func (this *Service) requestHandler(ev interface{}) (bool, error) {
	switch rv := ev.(type) {
	case func() error:
		return true, rv()
	}
	return true, nil
}

func (this *Service) stopHandler() {
	for _, si := range this.sessions {
		this.doCleanSession(si)
	}
}

func (this *Service) getSessionRoot(s *session) string {
	cfg, ok := this.config.Apps[string(s.appId)]
	r := this.config.TempDir
	if ok && cfg.Temp != "" {
		r = cfg.Temp
	}
	return this.getPathFile(r, string(s.sessionToken))
}

func (this *Service) getPathFile(r string, p string) string {
	return filepath.Join(r, p)
}

func (this *Service) doCleanSession(s *session) {
	logger.Debug("clean session(%s)", s)
	delete(this.sessions, string(s.sessionToken))
	dir := this.getSessionRoot(s)
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		logger.Debug("clean session dir - %s", dir)
		os.RemoveAll(dir)
	}
}

func (this *Service) doBeginCleaner(s *session, d time.Duration) {
	time.AfterFunc(d, func() {
		this.executor.DoNow("cleaner", func() error {
			this.doCleaner(s)
			return nil
		})
	})
}

func (this *Service) doCleaner(s *session) {
	token := string(s.sessionToken)
	if _, ok := this.sessions[token]; !ok {
		return
	}
	t := s.lastRequestTime.Add(time.Duration(this.config.TimeoutSec) * time.Second).Sub(time.Now())
	if t <= 0 {
		logger.Debug(tag, "session(%s) timeout", s)
		this.doCleanSession(s)
	} else {
		this.doBeginCleaner(s, t)
	}
}

func (this *Service) doCreateVcode(s *session) string {
	cfg, _ := this.config.Apps[string(s.appId)]
	key := ""
	if cfg != nil {
		key = cfg.Key
	}
	if key == "" {
		return ""
	}
	return filestoreprot.CreateVCode(string(s.sessionToken), key)
}

func (this *Service) doCheckVcode(s *session, vcode string) bool {
	m := this.doCreateVcode(s)
	if m != "" {
		if strings.ToLower(vcode) != strings.ToLower(m) {
			logger.Warn(tag, "vcode invalid - %s/%s", m, vcode)
			return false
		}
	}
	return true
}

func (this *Service) doBeginSession(appId filestoreprot.AppId) (string, error) {
	_, ok := this.config.Apps[string(appId)]
	if !ok {
		return "", fmt.Errorf("app '%s' not exists", appId)
	}
	var token string
	uid, err := uuid.NewV4()
	if err != nil {
		token = fmt.Sprintf("%d", time.Now().UnixNano())
	} else {
		token = uid.String()
	}
	s := new(session)
	s.appId = appId
	s.sessionToken = filestoreprot.SessionToken(token)
	s.lastRequestTime = time.Now()
	this.sessions[token] = s
	this.doBeginCleaner(s, time.Duration(this.config.TimeoutSec)*time.Second)
	return token, nil
}

func (this *Service) doSendFile(token string, path string, data []byte, vcode string) error {
	logger.Debug(tag, "sendFile(%s,%s,%d,%s)", token, path, len(data), vcode)
	s, _ := this.sessions[token]
	if s == nil {
		return fmt.Errorf("session(%s) invalid", token)
	}
	if !this.doCheckVcode(s, vcode) {
		return fmt.Errorf("vcode(%s) invalid", vcode)
	}
	s.lastRequestTime = time.Now()

	file := this.getPathFile(this.getSessionRoot(s), path)
	if logger.EnableDebug(tag) {
		_, ex := os.Stat(file)
		act := "create"
		if ex == nil {
			act = "append"
		}
		logger.Debug(tag, "%s file '%s'", act, file)
	}
	err := os.MkdirAll(filepath.Dir(file), 0664)
	if err != nil {
		return err
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	_, ex := os.Stat(file)
	if ex == nil {
		flag = os.O_WRONLY | os.O_APPEND
	}
	f, err2 := os.OpenFile(file, flag, 0664)
	if err2 != nil {
		return err
	}
	defer f.Close()
	n, err3 := f.Write(data)
	if err3 == nil && n < len(data) {
		err3 = io.ErrShortWrite
	}
	if err3 != nil {
		return err3
	}
	for _, old := range s.requests {
		if !old.deleteFile && strings.ToLower(old.path) == strings.ToLower(path) {
			return nil
		}
	}
	logger.Debug(tag, "new file request(%s)", path)

	req := new(request)
	req.path = path
	req.tempFile = file
	if s.requests == nil {
		s.requests = make([]*request)
	}
	s.requests = append(s.requests, req)
	return nil
}

func (this *Service) doRemoveFile(token string, path string, vcode string) error {
	logger.Debug(tag, "removeFile(%s,%s,%s)", token, path, vcode)
	s, _ := this.sessions[token]
	if s == nil {
		return fmt.Errorf("session(%s) invalid", token)
	}
	if !this.doCheckVcode(s, vcode) {
		return fmt.Errorf("vcode(%s) invalid", vcode)
	}
	s.lastRequestTime = time.Now()

	for _, old := range s.requests {
		if old.deleteFile && strings.ToLower(old.path) == strings.ToLower(path) {
			return nil
		}
	}

	logger.Debug(tag, "new delete request(%s)", path)
	req := new(request)
	req.deleteFile = true
	req.path = path
	if s.requests == nil {
		s.requests = make([]*request)
	}
	s.requests = append(s.requests, req)
	return nil
}
