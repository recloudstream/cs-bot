.PHONY: start

csbot: main.go body*.md
	go build .

start: csbot
	killall csbot
	nohup ./csbot > /var/log/csbot.log &2>1