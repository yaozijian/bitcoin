package main

import (
	"os"
	"os/signal"
	"time"

	"core"

	"router"

	log "github.com/cihub/seelog"
)

const def_console_log_cfg = `<seelog minlevel="info">
		<outputs formatid="detail">
			<console/>
		</outputs>
		<formats>
			<format id="common" format="%Msg%n" />
			<format id="detail" format="[%File:%Line][%Date(2006-01-02 15:04:05.000)] %Msg%n" />
		</formats>
</seelog>`

func idle() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Info("---------- Exit by Signal --------")
	time.Sleep(time.Millisecond * 100)
	os.Exit(100)
}

func main() {

	go idle()

	logger, _ := log.LoggerFromConfigAsString(def_console_log_cfg)

	if logger != nil {
		log.ReplaceLogger(logger)
	}

	r := router.NewRouter()

	ca, _ := core.NewUser("ca", "ca.crt", "ca.crt", "ca.key", r)
	a, _ := core.NewUser("a", "ca.crt", "server.crt", "server.key", r)

	b, _ := core.NewUser("b", "ca.crt", "server.crt", "server.key", r)
	c, _ := core.NewUser("c", "ca.crt", "server.crt", "server.key", r)
	d, _ := core.NewUser("d", "ca.crt", "server.crt", "server.key", r)
	e, _ := core.NewUser("e", "ca.crt", "server.crt", "server.key", r)

	err := ca.Register("a", 100, 0)
	log.Errorf("注册用户a: %v", err)

	err = ca.Register("b", 100, 0)
	log.Errorf("注册用户b: %v", err)

	// 这一步会失败: 只有CA用户能够注册用户
	err = a.Register("c", 100, 0)
	log.Errorf("注册用户c: %v", err)

	err = a.PayTo("b", 50, 0)
	log.Errorf("\n======= 支付1: %v", err)

	err = a.PayTo("b", 30, 0)
	log.Errorf("\n======= 支付2: %v", err)

	err = a.PayTo("b", 30, 0)
	log.Errorf("\n======= 支付3: %v", err)

	err = a.PayTo("c", 60, 0)
	log.Errorf("支付4: %v", err)

	ca.Show()
	a.Show()
	b.Show()
	c.Show()
	d.Show()
	e.Show()

	logger.Close()
}
