package search

import (
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
)

func setupQueue(topic string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, error) {
	uri := config.App.Databases["rabbitmq"]
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return conn, ch, &q, nil
}

func ScrapeFantasyGate(db *sqlx.DB, args *models.QueryBuilder) (interface{}, error) {
	rs := []models.Price{}

	conn, ch, q, err := setupQueue("prices2")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	defer ch.Close()

	collector := colly.NewCollector(
		colly.AllowedDomains("www.fantasygate.gr"),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)

	collector.OnHTML(".sblock4", func(e *colly.HTMLElement) {
		var stock int
		raw_price := e.ChildText(".jshop_price")

		if childHasClass(e, ".btn", "button_buy") {
			stock = 0
		} else {
			stock = 2
		}

		rs = append(rs, models.Price{
			Name:       e.ChildText(".name"),
			StoreId:    2,
			StoreThumb: e.ChildAttr(".jshop_img", "src"),
			Stock:      stock,
			Price:      getPrice(raw_price),
			Url:        e.Request.AbsoluteURL(e.ChildAttr(".name a", "href")),
		})
	})

	collector.Post("https://www.fantasygate.gr/strategygames", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/family-games", map[string]string{
		"limit": "99999",
	})

	collector.Post("https://www.fantasygate.gr/cardgames", map[string]string{
		"limit": "99999",
	})

	collector.Wait()

	err = updateBatch(db, 2)
	if err != nil {
		return nil, err
	}

	for _, item := range rs {
		err = upsertPrice(db, item)
		if err != nil {
			return nil, err
		}

		body, err := item.ToGOB64()
		if err != nil {
			return nil, err
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(body),
			},
		)
	}

	return rs, nil
}
