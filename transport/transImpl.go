package transport

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/antoha2/loghis/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (t *MQImpl) Init() {
	time.Sleep(time.Second * 12)
	conn, err := amqp.Dial(config.MQURL) // подключение к RabbitMQ
	if err != nil {
		log.Fatalf("%s: %s\n", "ошибка подключения AMQP", err)
	}
	defer conn.Close()

	amqpChannel, err := conn.Channel() // установка канала RabbitMQ
	if err != nil {
		log.Printf("%s: %s\n", "ошибка создания amqpChannel", err)
	}

	defer amqpChannel.Close()

	err = amqpChannel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Printf("%s: %s\n", "Failed to declare an exchange", err)
	}

	queue, err := amqpChannel.QueueDeclare( // объявляет очередь для хранения сообщений и их доставки потребителям.
		"logger", //  имя очереди
		true,     //   Сохранять ли
		false,    // Удаляется ли оно автоматически
		false,    // Это эксклюзив
		false,    // Следует ли блокировать
		nil,      // Дополнительные атрибуты
	)
	if err != nil {
		log.Printf("%s: %s\n", "ошибка создания amqpChannel", err)
	}

	err = amqpChannel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"logs",     // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("%s: %s\n", "Failed to bind a queue", err)
	}

	err = amqpChannel.Qos(3, 0, false) //Qos определяет, сколько сообщений или сколько байтов сервер будет пытаться сохранить в сети для потребителей, прежде чем получит подтверждение доставки.
	if err != nil {
		log.Printf("%s: %s\n", "ошибка конфигурирования Qos", err)
	}

	loggerChannel, err := amqpChannel.Consume( //обработка сообщений
		queue.Name, //  имя очереди
		"",         // Используется для различения нескольких потребителей
		false,      // Следует ли отвечать автоматически
		false,      // Это эксклюзив
		false,      // Если установлено значение true, это означает, что сообщение, отправленное в том же соединении, не может быть доставлено потребителям в этом соединении
		false,      // Заблокирована ли очередь сообщений
		nil,        // Дополнительные атрибуты
	)
	if err != nil {
		log.Printf("%s: %s\n", "ошибка регистрации loggerChannel", err)
	}

	stopChan := make(chan struct{})

	for d := range loggerChannel {

		logMsg := LoggerMsg{}
		if err := json.Unmarshal(d.Body, &logMsg); err != nil {
			log.Fatalf("error decoding JSON: %s", err)
		}

		//запись в файл
		t.service.Write(context.Background(), logMsg.Log)

		if err := d.Ack(false); err != nil { //Ack подтверждает доставку своим тегом доставки, когда он был использован с помощью Channel.Consume или Channel.Get.
			log.Fatalf("ошибка подтверждения: %s\n", err)
		} else {
			log.Println("подтверждение отправлено")
		}
	}

	<-stopChan
}

func (t *MQImpl) Stop() {

	if err := t.server.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
